package unit

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGenerateUnit(t *testing.T) {
	apiHandler := setupTestGenerate()

	tests := []struct {
		longURL            string
		expectedStatusCode int
		expectedBody       string
	}{
		{
			longURL:            "https://example.com/abcdefg",
			expectedStatusCode: http.StatusCreated,
			expectedBody:       "",
		},
		{
			longURL:            "",
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       "URL не должен быть пустым",
		},
		{
			longURL:            "invalid",
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       "URL должен быть больше 15 символов",
		},
	}

	for _, test := range tests {
		t.Run(test.longURL, func(t *testing.T) {
			body := `{"long": "` + test.longURL + `"}`
			req, err := http.NewRequest("POST", "/generate", strings.NewReader(body))
			assert.NoError(t, err)

			res := httptest.NewRecorder()
			apiHandler.ServeHTTP(res, req)

			assert.Equal(t, test.expectedStatusCode, res.Code)
			if test.expectedStatusCode != http.StatusCreated {
				assert.Contains(t, res.Body.String(), test.expectedBody)
			}
		})

	}
}

func setupTestGenerate() *gin.Engine {
	handler := SetupTest()

	router := gin.Default()
	router.POST("/generate", handler.CreateShortLink)

	return router
}

// func TestCreateShortLink_Success(t *testing.T) {
// 	body := `{"long": "https://example.com/abcdefg"}`
// 	response, err := setupTestGenerate(body)

// 	assert.NoError(t, err)
// 	assert.Equal(t, http.StatusCreated, response.Code)

// 	// var link model.Link
// 	// err := json.Unmarshal([]byte(w.Body.String()), &link)                     потом добавить, чтобы возвращался ответ
// 	// assert.NoError(t, err)
// 	// assert.Equal(t, "https://example.com/abcdefg", link.OriginalURL)
// }

// func TestCreateShortLink_EmptyString(t *testing.T) {
// 	body := `{"long":""}`
// 	response, err := setupTestGenerate(body)

// 	assert.NoError(t, err)
// 	assert.Equal(t, http.StatusBadRequest, response.Code)
// }

// func TestCreateShoreLink_InvalidURL(t *testing.T) {
// 	body := `{"long": "invalid"}`
// 	response, err := setupTestGenerate(body)

// 	assert.NoError(t, err)
// 	assert.Equal(t, http.StatusBadRequest, response.Code)
// }
