package service

import (
	"context"
	"fmt"

	"github.com/Noviiich/golang-url-shortener/internal/core/model"
)

type URLService struct {
	repo model.DB
}

func NewURLService(repo model.DB) *URLService {
	return &URLService{repo: repo}
}

func (s *URLService) All(ctx context.Context) ([]model.Link, error) {
	links, err := s.repo.All(ctx)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении всех коротких urls: %w", err)
	}
	return links, nil
}

func (s *URLService) Get(ctx context.Context, short string) (*model.Link, error) {
	link, err := s.repo.Get(ctx, short)
	if err != nil {
		return link, fmt.Errorf("ошибка при получении короткого url for indentifier '%s': %w", short, err)
	}
	return link, nil
}

func (s *URLService) Create(ctx context.Context, link *model.Link) error {
	err := s.repo.Create(ctx, link)
	if err != nil {
		return fmt.Errorf("ошибка создания короткого url: %w", err)
	}
	return nil
}

func (s *URLService) Delete(ctx context.Context, short string) error {
	err := s.repo.Delete(ctx, short)
	if err != nil {
		return fmt.Errorf("ошибка удаления короткого url for indentifier '%s': %w", short, err)
	}
	return nil
}
