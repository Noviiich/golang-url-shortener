package config

import (
	"os"

	_ "github.com/joho/godotenv/autoload"
)

type Config struct {
	MongoURI   string
	Database   string
	Collection string
}

func LoadConfig() *Config {
	return &Config{
		MongoURI:   os.Getenv("MONGO_URI"),
		Database:   os.Getenv("MONGO_DATABASE"),
		Collection: os.Getenv("MONGO_COLLECTION"),
	}
}
