package handler

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *URLHandler) RedirectShortURL(c *gin.Context) {
	shortID := c.Param("shortID")
	link, err := h.service.Get(context.Background(), shortID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.Redirect(http.StatusFound, link.OriginalURL)
}
