package handler

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"

	"github.com/Noviiich/golang-url-shortener/internal/core/model"
	"github.com/gin-gonic/gin"
)

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
	if !IsValidLink(requestBody.Long) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неправильный URL формат"})
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

func generateShortUrl(longUrl string) string {
	hasher := md5.New()
	hasher.Write([]byte(longUrl))
	shortUrl := hex.EncodeToString(hasher.Sum(nil))
	return shortUrl[:8]
}
