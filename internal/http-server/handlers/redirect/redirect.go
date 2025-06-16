package redirect

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	resp "github.com/Noviiich/golang-url-shortener/internal/lib/api/response"
	"github.com/Noviiich/golang-url-shortener/internal/lib/logger/sl"
	"github.com/Noviiich/golang-url-shortener/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type URLGetter interface {
	GetURL(alias string) (string, error)
}

//go:generate go run github.com/vektra/mockery/v2@latest --name=URLCacheGetter
type URLCacheGetter interface {
	GetURL(ctx context.Context, alias string) (string, error)
}

type Request struct {
	Alias string `json:"alias"`
}

type Response struct {
	Response resp.Response `json:"response"`
	URL      string        `json:"url"`
}

func New(log *slog.Logger, urlGetter URLGetter, cache URLCacheGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.redirect.New"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")
		if alias == "" {
			log.Info("alias is empty")
			render.JSON(w, r, resp.Error("invalid request"))
			return
		}

		cacheURL, err := cache.GetURL(r.Context(), alias)
		if err == nil {
			log.Info("cache url found", slog.String("url", cacheURL))
			http.Redirect(w, r, cacheURL, http.StatusFound)
			return
		}
		log.Info("failed to get url from cache", sl.Err(err))

		resURL, err := urlGetter.GetURL(alias)
		if err != nil {
			if errors.Is(err, storage.ErrURLNotFound) {
				log.Info("url not found", "alias", alias)
				render.JSON(w, r, resp.Error("not found"))
				return
			}

			log.Error("failed to get url", sl.Err(err))
			render.JSON(w, r, resp.Error("internal error"))
			return
		}

		log.Info("got url", slog.String("url", resURL))

		// редирект на оригинальный URL
		http.Redirect(w, r, resURL, http.StatusFound)
	}
}
