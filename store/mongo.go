package store

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Noviiich/golang-url-shortener/types"
	_ "github.com/joho/godotenv/autoload"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDBStore struct {
	Client         *mongo.Client
	Collection     *mongo.Collection
	URI            string
	databaseName   string
	collectionName string
}

func NewMongoDBStore(ctx context.Context, database, collection string) *MongoDBStore {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	MONGO_URI := os.Getenv("MONGO_URI")
	opts := options.Client().ApplyURI(MONGO_URI).SetServerAPIOptions(serverAPI)
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		log.Fatalf("ошибка при подключении к mongoDB: %v", err)
	}

	if err := client.Database("admin").RunCommand(ctx, bson.D{{Key: "ping", Value: 1}}).Err(); err != nil {
		log.Fatalf("ошибка проверки подключения: %v", err)
	}
	coll := client.Database(database).Collection(collection)

	return &MongoDBStore{
		Client:         client,
		Collection:     coll,
		collectionName: collection,
		URI:            MONGO_URI,
		databaseName:   database,
	}
}

func (d *MongoDBStore) All(ctx context.Context) ([]types.Link, error) {
	cursor, err := d.Collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, fmt.Errorf("документ не найден: %w", err)
	}
	defer cursor.Close(ctx)

	var links []types.Link
	if err = cursor.All(ctx, &links); err != nil {
		return nil, fmt.Errorf("ошибка декодирования документа: %w", err)
	}

	return links, nil
}

func (d *MongoDBStore) Get(ctx context.Context, short string) (*types.Link, error) {
	var link types.Link
	err := d.Collection.FindOne(ctx, bson.M{"short": short}).Decode(&link)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("документ с short '%s' не найден", short)
		}
		return nil, fmt.Errorf("ошибка поиска документа: %w", err)
	}

	return &link, nil
}

func (d *MongoDBStore) Create(ctx context.Context, link types.Link) error {
	_, err := d.Collection.InsertOne(ctx, link)
	if err != nil {
		return fmt.Errorf("не удалось добавить элемент в MongoDB: %w", err)
	}

	return nil
}

func (d *MongoDBStore) Delete(ctx context.Context, id string) error {
	_, err := d.Collection.DeleteOne(ctx, bson.M{"id": id})
	if err != nil {
		return fmt.Errorf("ошибка удаления элемента из MongoDB: %w", err)
	}

	return nil
}
