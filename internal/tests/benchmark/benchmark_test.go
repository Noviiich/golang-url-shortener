package benchmark

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/Noviiich/golang-url-shortener/internal/core/domain"
	"github.com/Noviiich/golang-url-shortener/internal/core/service"
	"github.com/Noviiich/golang-url-shortener/internal/tests/mock"
)

func GetService() *service.LinkService {
	cache := mock.NewMockRedisCache()
	mockLink := mock.NewMockRepository()

	linkService := service.NewLinkService(mockLink, cache)

	return linkService
}

func BenchmarkLinkServiceGetAll(b *testing.B) {
	service := GetService()
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := service.All(ctx)
		if err != nil {
			b.Fatalf("Benchmark GetAll failed: %v", err)
		}
	}
}

func BenchmarkLinkServiceGetOriginalURL(b *testing.B) {
	service := GetService()
	ctx := context.Background()
	shortID := "testid2"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := service.GetOriginalURL(ctx, shortID)
		if err != nil {
			b.Fatalf("Benchmark GetOriginalURL failed: %v", err)
		}
	}
}

func BenchmarkLinkServiceCreate(b *testing.B) {
	service := GetService()
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		newID := fmt.Sprint(i)
		link := &domain.Link{
			Id:          newID,
			OriginalURL: "https://example.com/" + newID,
			CreateAt:    time.Now(),
		}

		err := service.Create(ctx, link)
		if err != nil {
			b.Fatalf("Benchmark Create failed: %v", err)
		}
	}
}
