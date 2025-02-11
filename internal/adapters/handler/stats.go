package handler

import (
	"context"
	"log"
	"net/http"

	"github.com/Noviiich/golang-url-shortener/internal/core/service"
	"github.com/gin-gonic/gin"
)

type StatsFunctionHandler struct {
	statsService *service.StatsService
	linkService  *service.LinkService
}

func NewStatsFunctionHandler(l *service.LinkService, s *service.StatsService) *StatsFunctionHandler {
	return &StatsFunctionHandler{statsService: s, linkService: l}
}

func (h *StatsFunctionHandler) Stats(c *gin.Context) {
	links, err := h.linkService.All(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for i, link := range links {
		stats, err := h.statsService.GetStatsByLinkID(context.Background(), link.Id)
		if err != nil {
			log.Printf("Error getting stats for link '%s': %v", link.Id, err)
			continue
		}
		links[i].Stats = stats
	}
	c.JSON(http.StatusOK, links)
}
