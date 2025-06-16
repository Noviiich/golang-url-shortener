package redis

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Noviiich/golang-url-shortener/internal/storage"
	"github.com/redis/go-redis/v9"
)

type Storage struct {
	client *redis.Client
}

func New(addr string, passw string, db int) (*Storage, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: passw,
		DB:       db,
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		fmt.Printf("failed to connect to redis server: %s\n", err.Error())
		return nil, err
	}

	return &Storage{client: client}, nil
}

func (r *Storage) SaveURL(ctx context.Context, urlToSave string, alias string) error {
	const op = "storage.redis.SetURL"

	if err := r.client.Set(ctx, alias, urlToSave, time.Hour).Err(); err != nil {
		return fmt.Errorf("%s: failed to set short URL for identifier: %w", op, err)
	}

	return nil
}

func (r *Storage) GetURL(ctx context.Context, alias string) (string, error) {
	const op = "storage.redis.GetURL"

	resURL, err := r.client.Get(ctx, alias).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", fmt.Errorf("%s: %w", op, storage.ErrURLNotFound)
		}
		return "", fmt.Errorf("%s: execute statement: %w", op, err)
	}

	return resURL, nil
}

func (r *Storage) DeleteURL(ctx context.Context, alias string) error {
	const op = "storage.redis.DeleteURL"

	if err := r.client.Del(ctx, alias).Err(); err != nil {
		if errors.Is(err, redis.Nil) {
			return fmt.Errorf("%s: %w", op, storage.ErrURLNotFound)
		}
		return fmt.Errorf("%s: execute statement: %w", op, err)
	}
	return nil
}
