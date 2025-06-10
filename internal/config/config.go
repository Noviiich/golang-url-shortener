package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load("C:/Users/nowik/VSCode/Go/golang-url-shortener/.env")
	if err != nil {
		log.Print("Error loading .env file: ", err)
	}
}

type Config struct {
	Env             string `yaml:"env" env-required:"true"`
	StoragePath     string `yaml:"storage_path" env-required:"true"`
	HTTPServer      `yaml:"http_server"`
	MongoURI        string
	Database        string
	LinksCollection string
	StatsCollection string
	Redis           RedisConfig
}

type HTTPServer struct {
	Address     string `yaml:"address" env-default:"localhost:8082"`
	Timeout     string `yaml:"timeout" env-default:"4s"`
	IdleTimeout string `yaml:"idle_timeout" env-default:"30s"`
	User        string `yaml:"user" env-required:"true"`
	Password    string `yaml:"password" env-required:"true" env:"HTTP_SERVER_PASSWORD"`
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	// Check if the config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("error reading config file: %s", err)
	}

	return &cfg
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
