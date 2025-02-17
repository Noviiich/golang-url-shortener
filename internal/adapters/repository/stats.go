package repository

import (
	"context"
	"fmt"
	"log"

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

func NewStatsRepository(cfg *config.Config) *StatsRepository {
	if cfg.MongoURI == "" {
		log.Fatal("MongoDB URI is empty")
		return nil
	}
	if cfg.Database == "" {
		log.Fatal("MongoDB database is empty")
		return nil
	}
	if cfg.StatsCollection == "" {
		log.Fatal("MongoDB collection is empty")
		return nil
	}

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().
		ApplyURI(cfg.MongoURI).
		SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(context.Background(), opts)
	if err != nil {
		log.Fatalf("ошибка при подключении к mongoDB: %v", err)
		return nil
	}

	if err := client.Database("admin").RunCommand(context.Background(), bson.D{{Key: "ping", Value: 1}}).Err(); err != nil {
		log.Fatalf("ошибка проверки подключения: %v", err)
		return nil
	}

	stats := client.Database(cfg.Database).Collection(cfg.StatsCollection)

	return &StatsRepository{
		client:     client,
		collection: stats,
	}
}

func (r *StatsRepository) All(ctx context.Context) ([]domain.Stats, error) {
	cursor, err := r.collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, fmt.Errorf("документ не найден: %w", err)
	}
	defer cursor.Close(ctx)

	var links []domain.Stats
	if err = cursor.All(ctx, &links); err != nil {
		return nil, fmt.Errorf("ошибка декодирования документа: %w", err)
	}

	return links, nil
}

func (r *StatsRepository) Get(ctx context.Context, id string) (*domain.Stats, error) {
	var stats domain.Stats
	err := r.collection.FindOne(ctx, bson.M{"link_id": id}).Decode(&stats)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("документ с id '%s' не найден", id)
		}
		return nil, fmt.Errorf("ошибка поиска документа: %w", err)
	}

	return &stats, nil
}

func (r *StatsRepository) Create(ctx context.Context, stats *domain.Stats) error {
	_, err := r.collection.InsertOne(ctx, stats)
	if err != nil {
		return fmt.Errorf("не удалось добавить элемент в MongoDB: %w", err)
	}
	return nil
}

func (r *StatsRepository) Delete(ctx context.Context, id string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"id": id})
	if err != nil {
		return fmt.Errorf("ошибка удаления элемента из MongoDB: %w", err)
	}

	return nil
}

func (r *StatsRepository) GetStatsByLinkID(ctx context.Context, id string) ([]domain.Stats, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"id": id})
	if err != nil {
		return nil, fmt.Errorf("документ не найден: %w", err)
	}

	var links []domain.Stats
	if err = cursor.All(ctx, &links); err != nil {
		return nil, fmt.Errorf("ошибка декодирования документа: %w", err)
	}
	return links, nil
}
