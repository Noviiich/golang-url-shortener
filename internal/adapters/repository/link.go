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

type LinkRepository struct {
	Client     *mongo.Client
	Collection *mongo.Collection
}

func NewLinkRepository(cfg *config.Config) *LinkRepository {
	if cfg.MongoURI == "" {
		log.Fatal("MongoDB URI is empty")
		return nil
	}
	if cfg.Database == "" {
		log.Fatal("MongoDB database is empty")
		return nil
	}
	if cfg.LinksCollection == "" {
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

	coll := client.Database(cfg.Database).Collection(cfg.LinksCollection)

	return &LinkRepository{
		Client:     client,
		Collection: coll,
	}
}

func (r *LinkRepository) All(ctx context.Context) ([]domain.Link, error) {
	cursor, err := r.Collection.Find(ctx, bson.D{})
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

func (r *LinkRepository) Get(ctx context.Context, shortID string) (*domain.Link, error) {
	var link domain.Link
	err := r.Collection.FindOne(ctx, bson.M{"id": shortID}).Decode(&link)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("документ с id '%s' не найден", shortID)
		}
		return nil, fmt.Errorf("ошибка поиска документа: %w", err)
	}

	return &link, nil
}

func (r *LinkRepository) Create(ctx context.Context, link *domain.Link) error {
	_, err := r.Collection.InsertOne(ctx, link)
	if err != nil {
		return fmt.Errorf("не удалось добавить элемент в MongoDB: %w", err)
	}
	return nil
}

func (r *LinkRepository) Delete(ctx context.Context, shortID string) error {
	_, err := r.Collection.DeleteOne(ctx, bson.M{"id": shortID})
	if err != nil {
		return fmt.Errorf("ошибка удаления элемента из MongoDB: %w", err)
	}

	return nil
}
