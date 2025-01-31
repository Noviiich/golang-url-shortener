package main

import (
	"log"

	"github.com/Noviiich/golang-url-shortener/internal/config"
	"github.com/Noviiich/golang-url-shortener/internal/handler"
	"github.com/Noviiich/golang-url-shortener/internal/repository"
	"github.com/Noviiich/golang-url-shortener/internal/service"
	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	cfg := config.LoadConfig()

	repo, err := repository.NewURLRepository(cfg)
	if err != nil {
		log.Fatalf("ошибка создания репозитория: %v", err)
	}

	service := service.NewURLService(repo)
	handler := handler.NewURLHandler(service)

	router := gin.Default()
	router.POST("/", handler.CreateShortLink)
	router.GET("/:shortID", handler.RedirectShortURL)

	log.Println("Server started on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
