package handler

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"

	"github.com/Noviiich/golang-url-shortener/internal/core/domain"
	"github.com/Noviiich/golang-url-shortener/internal/core/service"
	"github.com/gin-gonic/gin"
)

type RequestBody struct {
	Long string `json:"long"`
}

type GenerateFunctionHandler struct {
	linkService  *service.LinkService
	statsService *service.StatsService
}

func NewGenerateFunctionHandler(l *service.LinkService, s *service.StatsService) *GenerateFunctionHandler {
	return &GenerateFunctionHandler{linkService: l, statsService: s}
}

func (h *GenerateFunctionHandler) CreateShortLink(c *gin.Context) {
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

	link := domain.Link{
		Id:          generateShortUrl(requestBody.Long),
		OriginalURL: requestBody.Long,
		CreateAt:    time.Now(),
	}

	err := h.linkService.Create(context.Background(), &link)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = h.statsService.Create(context.Background(), &domain.Stats{
		Platform:  domain.PlatformTwitter,
		LinkID:    link.Id,
		CreatedAt: link.CreateAt,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"key": link.Id})
}

func generateShortUrl(longUrl string) string {
	hasher := md5.New()
	hasher.Write([]byte(longUrl))
	shortUrl := hex.EncodeToString(hasher.Sum(nil))
	return shortUrl[:8]
}
