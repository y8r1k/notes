package config

import (
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env        string `yaml:"env" env-default:"local" env-required:"true"`
	DBConfig   `yaml:"db_postgres"`
	HTTPServer `yaml:"http_server"`
}

type HTTPServer struct {
	Addres string `yaml:"addres" env-default:"localhost:8080"`
}

type DBConfig struct {
	User     string `yaml:"user" env-required:"true"`
	Password string `yaml:"password" env-required:"true"`
	DBName   string `yaml:"dbname" env-required:"true"`
	SSLMode  string `yaml:"sslmode" env-required:"true"`
}

func MustLoad() *Config {
	// Getting CONFIG_PATH
	configPath := os.Getenv("CONFIG_PATH")

	if configPath == "" {
		log.Fatal("CONFIG_PATH environment variable is not set")
	}

	// Check existing CONFIG_PATH
	if _, err := os.Stat(configPath); err != nil {
		log.Fatalf("failed to open config file: %s", err)
	}
	var cfg Config

	// Reading configuration
	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		log.Fatalf("failed to read config file: %s", err)
	}

	return &cfg
}
