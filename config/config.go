package config

import (
	"log"

	"github.com/caarlos0/env"

	"github.com/joho/godotenv"
)

type DBConfig struct {
	Host     string `env:"DB_HOST"`
	Name     string `env:"DB_NAME"`
	Port     string `env:"DB_PORT"`
	User     string `env:"DB_USER"`
	Password string `env:"DB_PASSWORD"`
}
type APIConfig struct {
	Agify_API       string `env:"AGIFY_API"`
	Genderize_API   string `env:"GENDERIZE_API"`
	Nationalize_API string `env:"NATIONALIZE_API"`
}
type AppConfig struct {
	LogLevel string `env:"LOG_LEVEL"`
	Env      string `env:"APP_ENV" envDefault:"development"`
}

type Config struct {
	DB  DBConfig
	API APIConfig
	App AppConfig
}

var Cfg *Config

func LoadConfig() {
	Cfg = &Config{}
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Ошибка загрузки .env файла")
	}

	if err := env.Parse(&Cfg.DB); err != nil {
		log.Fatalf("Ошибка загрузки DB конфигурации:%v", err)
	}
	if err := env.Parse(&Cfg.API); err != nil {
		log.Fatalf("Ошибка загрузки API конфигурации: %v", err)
	}
	if err := env.Parse(&Cfg.App); err != nil {
		log.Fatalf("Ошибка загрузки APP конфигурвции: %v", err)
	}

}
