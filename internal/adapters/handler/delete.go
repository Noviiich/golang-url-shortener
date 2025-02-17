package handler

import (
	"context"
	"net/http"

	"github.com/Noviiich/golang-url-shortener/internal/core/service"
	"github.com/gin-gonic/gin"
)

type DeleteFunctionHandler struct {
	linkService  *service.LinkService
	statsService *service.StatsService
}

func NewDeleteFunctionHandler(l *service.LinkService, s *service.StatsService) *DeleteFunctionHandler {
	return &DeleteFunctionHandler{linkService: l, statsService: s}
}

func (s *DeleteFunctionHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := s.linkService.Delete(context.Background(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	if err := s.statsService.Delete(context.Background(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusNoContent, gin.H{})
}
