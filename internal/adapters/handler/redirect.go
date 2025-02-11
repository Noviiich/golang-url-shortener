package handler

import (
	"context"
	"net/http"

	"github.com/Noviiich/golang-url-shortener/internal/core/service"
	"github.com/gin-gonic/gin"
)

type RedirectFunctionHandler struct {
	linkService *service.LinkService
}

func NewRedirectFunctionHandler(l *service.LinkService) *RedirectFunctionHandler {
	return &RedirectFunctionHandler{linkService: l}
}

func (h *RedirectFunctionHandler) RedirectShortURL(c *gin.Context) {
	shortID := c.Param("shortID")
	link, err := h.linkService.GetOriginalURL(context.Background(), shortID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.Redirect(http.StatusFound, *link)
}
