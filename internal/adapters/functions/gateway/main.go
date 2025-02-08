package main

import (
	"log"
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
)

func reverseProxy(target string) gin.HandlerFunc {
	return func(c *gin.Context) {
		targetURL, err := url.Parse(target)
		if err != nil {
			log.Fatalf("ошибка парсинга url: %v", err)
		}

		proxy := httputil.NewSingleHostReverseProxy(targetURL)
		proxy.ServeHTTP(c.Writer, c.Request)
	}
}

func main() {
	router := gin.Default()
	router.POST("/genarate", reverseProxy("https://golang-url-shortener-esii.onrender.com:8081"))
	router.GET("/:shortID", reverseProxy("https://golang-url-shortener-esii.onrender.com:8082"))

	log.Println("API Gateway запущен на порту :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
