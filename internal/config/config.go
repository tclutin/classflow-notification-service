package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"log"
)

const (
	dev  string = "local"
	prod string = "prod"
)

type Config struct {
	Environment string `env:"ENVIRONMENT"`
	HTTPServer
	Postgres
	Telegram
}

type Telegram struct {
	Token string `env:"TELEGRAM_TOKEN"`
}

type HTTPServer struct {
	Address string `env:"HTTP_HOST"`
	Port    string `env:"HTTP_PORT"`
}

type Postgres struct {
	Host     string `env:"POSTGRES_HOST"`
	Port     string `env:"POSTGRES_PORT"`
	DbName   string `env:"POSTGRES_DB"`
	User     string `env:"POSTGRES_USER"`
	Password string `env:"POSTGRES_PASSWORD"`
}

func MustLoad() *Config {
	var config Config

	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("error loading .env file: %v\n", err)
	}

	if err := cleanenv.ReadEnv(&config); err != nil {
		log.Fatalf("error loading configuration: %v\n", err)
	}

	return &config
}

func (c *Config) IsProd() bool {
	return c.Environment == prod
}

func (c *Config) IsLocal() bool {
	return c.Environment == dev
}
