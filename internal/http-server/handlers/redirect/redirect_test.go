package redirect_test

import (
	"net/http/httptest"
	"testing"

	"github.com/Noviiich/golang-url-shortener/internal/http-server/handlers/redirect"
	"github.com/Noviiich/golang-url-shortener/internal/http-server/handlers/redirect/mocks"
	"github.com/Noviiich/golang-url-shortener/internal/lib/api"
	"github.com/Noviiich/golang-url-shortener/internal/lib/logger/handlers/slogdiscard"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetHandler(t *testing.T) {
	cases := []struct {
		alias     string
		respError string
		mockError error
		url       string
		name      string
	}{
		{
			name:  "Success",
			alias: "test",
			url:   "https://google.com",
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {

			urlGetterMock := mocks.NewURLGetter(t)

			if tc.respError == "" || tc.mockError != nil {
				urlGetterMock.On("GetURL", tc.alias).
					Return(tc.url, tc.mockError).Once()
			}

			r := chi.NewRouter()
			r.Get("/{alias}", redirect.New(slogdiscard.NewDiscardLogger(), urlGetterMock))

			ts := httptest.NewServer(r)
			defer ts.Close()

			redirectedToURL, err := api.GetRedirect(ts.URL + "/" + tc.alias)
			require.NoError(t, err)

			// Check the final URL after redirection.
			assert.Equal(t, tc.url, redirectedToURL)

		})
	}
}
