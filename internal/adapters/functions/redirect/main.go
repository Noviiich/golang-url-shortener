package main

import (
	"log"

	"github.com/Noviiich/golang-url-shortener/internal/adapters/cache"
	"github.com/Noviiich/golang-url-shortener/internal/adapters/handler"
	"github.com/Noviiich/golang-url-shortener/internal/adapters/repository"
	"github.com/Noviiich/golang-url-shortener/internal/config"
	"github.com/Noviiich/golang-url-shortener/internal/core/service"
	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	cfg := config.LoadConfig()
	redisAddress, redisPassword, redisDB := cfg.GetRedisParams()

	repo := repository.NewLinkRepository(cfg)
	cache := cache.NewRedisCache(redisAddress, redisPassword, redisDB)
	service := service.NewLinkService(repo, cache)
	handler := handler.NewRedirectFunctionHandler(service)

	router := gin.Default()
	router.GET("/:shortID", handler.RedirectShortURL)
	if err := router.Run(":8082"); err != nil {
		log.Fatal(err)
	}
}
