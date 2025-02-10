package repository

import (
	"context"
	"fmt"

	"github.com/Noviiich/golang-url-shortener/internal/config"
	"github.com/Noviiich/golang-url-shortener/internal/core/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type StatsRepository struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func NewStatsRepository(cfg *config.Config) (*StatsRepository, error) {
	if cfg.MongoURI == "" {
		return nil, fmt.Errorf("MongoDB URI is empty")
	}
	if cfg.Database == "" {
		return nil, fmt.Errorf("MongoDB database is empty")
	}
	if cfg.Collection == "" {
		return nil, fmt.Errorf("MongoDB collection is empty")
	}

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().
		ApplyURI(cfg.MongoURI).
		SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(context.Background(), opts)
	if err != nil {
		return nil, fmt.Errorf("ошибка при подключении к mongoDB: %w", err)
	}

	if err := client.Database("admin").RunCommand(context.Background(), bson.D{{Key: "ping", Value: 1}}).Err(); err != nil {
		return nil, fmt.Errorf("ошибка проверки подключения: %w", err)
	}

	coll := client.Database(cfg.Database).Collection(cfg.Collection)

	return &StatsRepository{
		client:     client,
		collection: coll,
	}, nil
}

func (r *StatsRepository) All(ctx context.Context) ([]domain.Link, error) {
	cursor, err := r.collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, fmt.Errorf("документ не найден: %w", err)
	}
	defer cursor.Close(ctx)

	var links []domain.Link
	if err = cursor.All(ctx, &links); err != nil {
		return nil, fmt.Errorf("ошибка декодирования документа: %w", err)
	}

	return links, nil
}

func (r *StatsRepository) Get(ctx context.Context, shortID string) (*domain.Link, error) {
	var link domain.Link
	err := r.collection.FindOne(ctx, bson.M{"short_id": shortID}).Decode(&link)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("документ с short_id '%s' не найден", shortID)
		}
		return nil, fmt.Errorf("ошибка поиска документа: %w", err)
	}

	return &link, nil
}

func (r *StatsRepository) Create(ctx context.Context, link *domain.Link) error {
	_, err := r.collection.InsertOne(ctx, link)
	if err != nil {
		return fmt.Errorf("не удалось добавить элемент в MongoDB: %w", err)
	}
	return nil
}

func (r *StatsRepository) Delete(ctx context.Context, shortID string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"short_id": shortID})
	if err != nil {
		return fmt.Errorf("ошибка удаления элемента из MongoDB: %w", err)
	}

	return nil
}
