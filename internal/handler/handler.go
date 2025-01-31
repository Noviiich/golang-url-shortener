package handler

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"

	"github.com/Noviiich/golang-url-shortener/internal/model"
	"github.com/Noviiich/golang-url-shortener/internal/service"
	"github.com/gin-gonic/gin"
)

type URLHandler struct {
	service *service.URLService
}

func NewURLHandler(s *service.URLService) *URLHandler {
	return &URLHandler{service: s}
}

type RequestBody struct {
	Long string `json:"long"`
}

func (h *URLHandler) CreateShortLink(c *gin.Context) {
	var requestBody RequestBody
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		fmt.Println(requestBody)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if requestBody.Long == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "URL не должен быть пустым"})
		return
	}
	if len(requestBody.Long) < 15 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "URL должен быть больше 15 символов"})
		return
	}

	link := model.Link{
		ShortId:     generateShortUrl(requestBody.Long),
		OriginalURL: requestBody.Long,
		CreateAt:    time.Now(),
		Clicks:      0,
	}

	err := h.service.Create(context.Background(), &link)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{})
}

func (h *URLHandler) RedirectShortURL(c *gin.Context) {
	shortID := c.Param("shortID")
	link, err := h.service.Get(context.Background(), shortID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.Redirect(http.StatusFound, link.OriginalURL)
}

// func (h *UrlShortenerHandler) GetOriginalUrl(w http.ResponseWriter, r *http.Request) {
// 	shortLink := r.URL.Path[1:]
// 	longLink, err := h.link.Get(context.Background(), shortLink)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusNotFound)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusOK)
// 	response := map[string]string{
// 		"long": longLink.OriginalURL,
// 	}
// 	json.NewEncoder(w).Encode(response)

// }

func generateShortUrl(longUrl string) string {
	hasher := md5.New()
	hasher.Write([]byte(longUrl))
	shortUrl := hex.EncodeToString(hasher.Sum(nil))
	return shortUrl[:8]
}
