package auth

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/Noviiich/golang-url-shortener/internal/lib/jwt"
	"github.com/Noviiich/golang-url-shortener/internal/lib/logger/sl"
)

// Определяем тип для ключей контекста (лучшая практика)
type contextKey string

// Определяем константы для ключей
const (
	errorKey   contextKey = "auth_error"
	uidKey     contextKey = "user_id"
	isAdminKey contextKey = "is_admin"
)

var (
	ErrInvalidToken       = errors.New("invalid token")
	ErrFailedIsAdminCheck = errors.New("failed to check if user is admin")
)

type PermissionProvider interface {
	IsAdmin(ctx context.Context, userID int64) (bool, error)
}

// New creates new auth middleware.
func New(
	log *slog.Logger,
	appSecret string,
	permProvider PermissionProvider,
) func(next http.Handler) http.Handler {
	const op = "middleware.auth.New"

	log = log.With(slog.String("op", op))

	// Возвращаем функцию-обработчик
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Получаем JWT-токен из запроса
			tokenStr := extractBearerToken(r)
			if tokenStr == "" {
				// It's ok, if user is not authorized
				next.ServeHTTP(w, r)
				return
			}

			// Парсим и валидируем токен, использeуя appSecret
			claims, err := jwt.Parse(tokenStr, appSecret)
			if err != nil {
				log.Warn("failed to parse token", sl.Err(err))

				// But if token is invalid, we shouldn't handle request
				ctx := context.WithValue(r.Context(), errorKey, ErrInvalidToken)
				next.ServeHTTP(w, r.WithContext(ctx))

				return
			}

			log.Info("user authorized", slog.Any("claims", claims))

			//Отправляем запрос для проверки, является ли пользователь админов
			isAdmin, err := permProvider.IsAdmin(r.Context(), claims.UID)
			if err != nil {
				log.Error("failed to check if user is admin", sl.Err(err))

				ctx := context.WithValue(r.Context(), errorKey, ErrFailedIsAdminCheck)
				next.ServeHTTP(w, r.WithContext(ctx))

				return
			}

			// Полученны данные сохраняем в контекст,
			// откуда его смогут получить следующие хэндлеры.
			ctx := context.WithValue(r.Context(), uidKey, claims.UID)
			ctx = context.WithValue(ctx, isAdminKey, isAdmin)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func extractBearerToken(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	splitToken := strings.Split(authHeader, "Bearer ")
	if len(splitToken) != 2 {
		return ""
	}

	return splitToken[1]
}

func UIDFromContext(ctx context.Context) (int64, bool) {
	uid, ok := ctx.Value(uidKey).(int64)
	return uid, ok
}

func ErrorFromContext(ctx context.Context) (error, bool) {
	err, ok := ctx.Value(errorKey).(error)
	return err, ok
}

func IsAdminFromContext(ctx context.Context) (bool, bool) {
	isAdmin, ok := ctx.Value(isAdminKey).(bool)
	return isAdmin, ok
}
