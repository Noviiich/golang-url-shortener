package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Print("Error loading .env file: ", err)
	}
}

type Config struct {
	MongoURI      string
	Database      string
	Collection    string
	redisAddress  string
	redisPassword string
	redisDB       int
}

func LoadConfig() *Config {
	return &Config{
		MongoURI:      os.Getenv("MONGO_URI"),
		Database:      os.Getenv("MONGO_DATABASE"),
		Collection:    os.Getenv("MONGO_COLLECTION"),
		redisAddress:  "localhost:6379",
		redisPassword: "",
		redisDB:       0,
	}
}

func (c *Config) GetRedisParams() (string, string, int) {
	address, ok := os.LookupEnv("REDIS_ADDRESS")
	if !ok {
		fmt.Println("Need REDIS_ADDRESS environment variable")
		return c.redisAddress, c.redisPassword, c.redisDB
	}

	password, ok := os.LookupEnv("REDIS_PASSWORD")
	if !ok {
		fmt.Println("Need REDIS_PASSWORD evironment variable")
		return address, c.redisPassword, c.redisDB
	}

	dbStr, ok := os.LookupEnv("REDIS_DB")
	if !ok {
		fmt.Println("Need REDIS_DB evironment variable")
		return address, password, c.redisDB
	}

	db, err := strconv.Atoi(dbStr)
	if err != nil {
		fmt.Printf("REDIS_DB evironment variable is not a valid integer: %v\n", err)
		return address, password, c.redisDB
	}

	return address, password, db
}
