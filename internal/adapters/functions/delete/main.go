package main

import (
	"log"

	"github.com/Noviiich/golang-url-shortener/internal/adapters/cache"
	"github.com/Noviiich/golang-url-shortener/internal/adapters/handler"
	"github.com/Noviiich/golang-url-shortener/internal/adapters/repository"
	"github.com/Noviiich/golang-url-shortener/internal/config"
	"github.com/Noviiich/golang-url-shortener/internal/core/service"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()
	redisAddress, redisPassword, redisDB := cfg.GetRedisParams()
	cache := cache.NewRedisCache(redisAddress, redisPassword, redisDB)

	linkRepo := repository.NewLinkRepository(cfg)
	statsRepo := repository.NewStatsRepository(cfg)

	linkService := service.NewLinkService(linkRepo, cache)
	statsService := service.NewStatsService(statsRepo, cache)
	handler := handler.NewDeleteFunctionHandler(linkService, statsService)

	router := gin.Default()
	router.DELETE("/:id", handler.Delete)
	if err := router.Run(":8084"); err != nil {
		log.Fatal(err)
	}
}
