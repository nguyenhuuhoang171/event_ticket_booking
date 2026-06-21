package config

import (
	"encoding/json"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Server         ServerConfig
	Db             DbConfig
	Authentication AuthenticationConfig
	Redis          RedisConfig
}

type ServerConfig struct {
	Port string
}

type DbConfig struct {
	Username string
	Password string
	Host     string
	Port     string
	DbName   string
}

type AuthenticationConfig struct {
	AccessTokenExpirationMinutes int // minutes
}

type RedisConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	Db       int
}

func NewConfig() Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Fail to load file .env: %v", err)
	}

	var configJson []byte
	config := Config{}
	switch env := os.Getenv("ENVIRONMENT"); env {
	case "dev":
		configJson, err = os.ReadFile("config/config-dev.json")
	default:
		log.Fatalf("Environment invalid")
	}
	if err != nil {
		log.Fatalf("Fail to get config: %v", err)
	}
	err = json.Unmarshal(configJson, &config)
	if err != nil {
		log.Fatalf("Fail to unmarshal config: %v", err)
	}

	return config
}
