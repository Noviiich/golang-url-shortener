package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	ssogrpc "github.com/Noviiich/golang-url-shortener/internal/clients/sso/grcp"
	"github.com/Noviiich/golang-url-shortener/internal/config"
	"github.com/Noviiich/golang-url-shortener/internal/http-server/handlers/redirect"
	"github.com/Noviiich/golang-url-shortener/internal/http-server/handlers/url/save"
	"github.com/Noviiich/golang-url-shortener/internal/http-server/middleware/auth"
	mwLogger "github.com/Noviiich/golang-url-shortener/internal/http-server/middleware/logger"
	"github.com/Noviiich/golang-url-shortener/internal/lib/logger/handlers/slogpretty"
	"github.com/Noviiich/golang-url-shortener/internal/storage/redis"
	"github.com/Noviiich/golang-url-shortener/internal/storage/sqlite"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var (
	redisAddr     = "localhost:6379"
	redisPassword = ""
	redisDB       = 0
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)
	log.Info("url-shortener started", slog.String("env", cfg.Env))

	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Error("failed to initialize storage")
		os.Exit(1)
	}

	log.Info("staring db")

	cache, err := redis.New(
		redisAddr, redisPassword, redisDB,
	)
	if err != nil {
		log.Error("failed to initialize redis")
		os.Exit(1)
	}

	log.Info("staring redis")

	ssoClient, err := ssogrpc.New(
		context.Background(),
		log, cfg.Clients.SSO.Address,
		cfg.Clients.SSO.Timeout,
		cfg.Clients.SSO.RetriesCount,
	)
	if err != nil {
		log.Error("failed to initialize sso client")
		os.Exit(1)
	}

	//middleware
	router := chi.NewRouter()
	router.Use(middleware.RequestID) // Добавляет request_id в каждый запрос, для трейсинга
	router.Use(middleware.Logger)    // Логирование всех запросов
	router.Use(mwLogger.New(log))    // внутренний логгер
	router.Use(middleware.Recoverer) // Если где-то внутри сервера (обработчика запроса) произойдет паника, приложение не должно упасть
	router.Use(middleware.URLFormat) // Парсер URLов поступающих запросов

	// Хэндлер redirect остается снаружи, в основном роутере
	router.Get("/{alias}", redirect.New(log, storage, cache))

	// Группа маршрутов с SSO аутентификацией
	router.Route("/url", func(r chi.Router) {
		// Заменяем BasicAuth на SSO middleware
		r.Use(auth.New(log, cfg.AppSecret, ssoClient))

		// TODO: добавить обработчики
		r.Post("/", save.New(log, storage, cache))
		// r.Delete("/{alias}", delete.New(log, storage))
		// r.Get("/{alias}/stats", stats.New(log, storage))
	})

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server")
	}

	log.Error("server stopped")

}

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		log = setupPrettySlog()
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	default: // если env конфиг не валидный, устанавливает настройки prod
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
