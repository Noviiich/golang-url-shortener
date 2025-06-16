package save

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/Noviiich/golang-url-shortener/internal/http-server/middleware/auth"
	resp "github.com/Noviiich/golang-url-shortener/internal/lib/api/response"
	"github.com/Noviiich/golang-url-shortener/internal/lib/logger/sl"
	"github.com/Noviiich/golang-url-shortener/internal/lib/random"
	"github.com/Noviiich/golang-url-shortener/internal/storage"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type Request struct {
	URL   string `json:"url" validate:"required,url"` // validate нужен для проверки поля url вылидатором
	Alias string `json:"alias,omitempty"`
}

const aliasLength = 6

type Response struct {
	resp.Response
	Alias string `json:"alias"`
}

type URLSaver interface {
	SaveURL(urlToSave string, alias string) (int64, error)
}

//go:generate go run github.com/vektra/mockery/v2@latest --name=URLCacheSaver
type URLCacheSaver interface {
	SaveURL(ctx context.Context, urlToSave string, alias string) error
}

func New(log *slog.Logger, urlSaver URLSaver, cache URLCacheSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.save.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		ctx := r.Context()

		// Проверяем наличие ошибки аутентификации
		if authErr, hasError := auth.ErrorFromContext(ctx); hasError {
			log.Error("authentication failed", sl.Err(authErr))
			render.JSON(w, r, resp.Error("authentication required"))
			return
		}

		// Получаем ID пользователя
		userID, hasUID := auth.UIDFromContext(ctx)
		if !hasUID {
			log.Error("user ID not found in context")
			render.JSON(w, r, resp.Error("unauthorized"))
			return
		}

		// Получаем информацию о том, является ли пользователь админом
		isAdmin, _ := auth.IsAdminFromContext(ctx)
		// if !hasUID {
		// 	log.Error("user ID not found in context")
		// 	render.JSON(w, r, resp.Error("unauthorized"))
		// 	return
		// }
		// Добавляем информацию о пользователе в логи
		log = log.With(
			slog.Int64("user_id", userID),
			slog.Bool("is_admin", isAdmin),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		log.Debug(req.Alias)
		if err != nil {
			// обработка ошибки при пустом теле запроса
			if errors.Is(err, io.EOF) {
				log.Error("request body is empty")
				render.JSON(w, r, resp.Error("empty request"))
				return
			}

			log.Error("failed to decode request body", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to decode request"))
			return
		}

		log.Info("request body decoded", slog.Any("req", req))

		// Создаем объект валидатора
		// и передаем в него структуру, которую нужно провалидировать
		if err := validator.New().Struct(req); err != nil {
			// Приводим ошибку к типу ошибки валидации
			validateErr := err.(validator.ValidationErrors)

			log.Error("invalid request", sl.Err(err))

			render.JSON(w, r, resp.ValidationError(validateErr))
			return
		}

		alias := req.Alias
		if alias == "" {
			alias = random.NewRandomString(aliasLength)
		}

		id, err := urlSaver.SaveURL(req.URL, alias)
		if err != nil {
			if errors.Is(err, storage.ErrURLExists) {
				log.Info("url already exists", slog.String("url", req.URL))
				render.JSON(w, r, resp.Error("url already exists"))
				return
			}
			log.Error("failed to add url", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to add url"))
			return
		}

		err = cache.SaveURL(ctx, req.URL, alias)
		if err != nil {
			log.Info("failed to save cache url ", sl.Err(err))
		}
		log.Info("url added", slog.Int64("id", id))
		responseOK(w, r, alias)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, alias string) {
	render.JSON(w, r, Response{
		Response: resp.OK(),
		Alias:    alias,
	})
}
