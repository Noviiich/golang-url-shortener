package api

import (
	"errors"
	"fmt"
	"net/http"
)

var (
	ErrInvalidStatusCode = errors.New("invalid status code")
)

// GetRedirect returns the final URL after redirection.
func GetRedirect(url string) (string, error) {
	const op = "api.GetRedirect"

	// Делаем клиента, который проверяет редиректы
	//
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse // Наш клиент останавливается после первого редиректа
		},
	}

	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer func() { _ = resp.Body.Close() }()

	//Проверяем, что 302 статус возвращается
	if resp.StatusCode != http.StatusFound {
		return "", fmt.Errorf("%s: %w: %d", op, ErrInvalidStatusCode, resp.StatusCode)
	}

	// // Возвращаем тот url, на который происходит redirect
	return resp.Header.Get("Location"), err
}
