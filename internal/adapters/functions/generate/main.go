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

	repo, err := repository.NewURLRepository(cfg)
	if err != nil {
		log.Fatalf("ошибка создания репозитория: %v", err)
	}
	cache := cache.NewRedisCache(redisAddress, redisPassword, redisDB)
	service := service.NewURLService(repo, cache)
	handler := handler.NewURLHandler(service)

	router := gin.Default()
	router.POST("/genarate", handler.CreateShortLink)
	if err := router.Run(":8081"); err != nil {
		log.Fatal(err)
	}
}
