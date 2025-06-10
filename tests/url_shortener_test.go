package tests

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/Noviiich/golang-url-shortener/internal/http-server/handlers/url/save"
	"github.com/Noviiich/golang-url-shortener/internal/lib/api"
	"github.com/Noviiich/golang-url-shortener/internal/lib/random"
	"github.com/brianvoe/gofakeit"
	"github.com/gavv/httpexpect/v2"
	"github.com/stretchr/testify/require"
)

const (
	host = "localhost:8082"
)

func TestURLShortener_HappyPath(t *testing.T) {
	u := url.URL{
		Scheme: "http",
		Host:   host,
	}

	// Создаем клиент httpexpect
	e := httpexpect.Default(t, u.String())

	e.POST("/url"). // Отправляем POST-запрос, путь - '/url'
			WithJSON(save.Request{ // Формируем тело запроса
			URL:   gofakeit.URL(),             // Генерируем случайный URL
			Alias: random.NewRandomString(10), // Генерируем случайную строку
		}).
		WithBasicAuth("my_user", "my_pass"). // Добавляем к запросу креды авторизации
		Expect().                            // Далее перечисляем наши ожидания от ответа
		Status(200).                         // Код должен быть 200
		JSON().Object().                     // Получаем JSON-объект тела ответа
		ContainsKey("alias")                 // Проверяем, что в нём есть ключ 'alias'
}

func TestURLShortener_SaveRedirect(t *testing.T) {
	testCases := []struct {
		name  string
		url   string
		alias string
		error string
	}{
		{
			name:  "Valid URL",
			url:   gofakeit.URL(),
			alias: gofakeit.Word() + gofakeit.Word(),
		},
		{
			name:  "Invalid URL",
			url:   "invalid_url",
			alias: gofakeit.Word(),
			error: "field URL is not a valid URL",
		},
		{
			name:  "Empty Alias",
			url:   gofakeit.URL(),
			alias: "",
		},
		// TODO: add more test cases
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			u := url.URL{
				Scheme: "http",
				Host:   host,
			}

			e := httpexpect.Default(t, u.String())

			// Save

			resp := e.POST("/url").
				WithJSON(save.Request{
					URL:   tc.url,
					Alias: tc.alias,
				}).
				WithBasicAuth("my_user", "my_pass").
				Expect().Status(http.StatusOK).
				JSON().Object()

			if tc.error != "" {
				resp.NotContainsKey("alias")

				resp.Value("error").String().IsEqual(tc.error)

				return
			}

			alias := tc.alias

			if tc.alias != "" {
				resp.Value("alias").String().IsEqual(tc.alias)
			} else {
				resp.Value("alias").String().NotEmpty()

				alias = resp.Value("alias").String().Raw()
			}

			// Redirect

			testRedirect(t, alias, tc.url)

			// remove
		})
	}
}

func testRedirect(t *testing.T, alias string, urlToRedirect string) {
	u := url.URL{
		Scheme: "http",
		Host:   host,
		Path:   alias,
	}

	redirectedToURL, err := api.GetRedirect(u.String())
	require.NoError(t, err)

	require.Equal(t, urlToRedirect, redirectedToURL)
}
