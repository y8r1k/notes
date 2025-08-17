package config

import (
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
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
	// Получаем путь до конфиг-файла из env-переменной CONFIG_PATH
	// configPath := os.Getenv("CONFIG_PATH")
	configPath := `C:\Users\NoRoot\Прога\projects\Practice\Golang\Practice by Chatgpt\8. Application architecture\REST API NOTES\config\local.yaml`

	if configPath == "" {
		fmt.Fprintf(os.Stderr, "CONFIG_PATH environment variable is not set")
		os.Exit(1)
	}

	// Проверяем существование конфиг-файла
	if _, err := os.Stat(configPath); err != nil {
		fmt.Fprintf(os.Stderr, "error opening config file: %s", err)
		os.Exit(1)
	}
	var cfg Config

	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading config file: %s", err)
		os.Exit(1)
	}

	return &cfg
}
