package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load("C:/Users/nowik/VSCode/Go/golang-url-shortener/.env")
	if err != nil {
		log.Print("Error loading .env file: ", err)
	}
}

type Config struct {
	MongoURI        string
	Database        string
	LinksCollection string
	StatsCollection string
	Redis           RedisConfig
}

type RedisConfig struct {
	Address  string
	Password string
	DB       int
}

func LoadConfig() *Config {
	return &Config{
		MongoURI:        os.Getenv("MONGO_URI"),
		Database:        os.Getenv("MONGO_DATABASE"),
		LinksCollection: os.Getenv("MONGO_COLLECTION"),
		StatsCollection: os.Getenv("MONGO_STATS"),
		Redis: RedisConfig{
			Address:  "localhost:6379",
			Password: "",
			DB:       0,
		},
	}
}

func (c *Config) GetRedisParams() (string, string, int) {
	address, ok := os.LookupEnv("REDIS_ADDRESS")
	if !ok {
		fmt.Println("Need REDIS_ADDRESS environment variable")
		return c.Redis.Address, c.Redis.Password, c.Redis.DB
	}

	password, ok := os.LookupEnv("REDIS_PASSWORD")
	if !ok {
		fmt.Println("Need REDIS_PASSWORD evironment variable")
		return address, c.Redis.Password, c.Redis.DB
	}

	dbStr, ok := os.LookupEnv("REDIS_DB")
	if !ok {
		fmt.Println("Need REDIS_DB evironment variable")
		return address, password, c.Redis.DB
	}

	db, err := strconv.Atoi(dbStr)
	if err != nil {
		fmt.Printf("REDIS_DB evironment variable is not a valid integer: %v\n", err)
		return address, password, c.Redis.DB
	}

	return address, password, db
}
