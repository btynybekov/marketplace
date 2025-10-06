package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Server struct {
		Port int `yaml:"port"`
	} `yaml:"server"`
	Database struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Name     string `yaml:"name"`
		SSLMode  string `yaml:"sslmode"`
	} `yaml:"database"`
}

var Cfg Config

func LoadConfig(path string) {
	f, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("Failed to read config.yaml: %v", err)
	}

	if err := yaml.Unmarshal(f, &Cfg); err != nil {
		log.Fatalf("Failed to parse config.yaml: %v", err)
	}

	_ = godotenv.Load()

	if url := os.Getenv("DATABASE_URL"); url != "" {
		Cfg.Database.Port = 0
		Cfg.Database.User = ""
		Cfg.Database.Password = ""
		Cfg.Database.Name = ""
		Cfg.Database.SSLMode = ""
	}
}

func GetDSN() string {
	if url := os.Getenv("DATABASE_URL"); url != "" {
		return url
	}

	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		Cfg.Database.User,
		Cfg.Database.Password,
		Cfg.Database.Host,
		Cfg.Database.Port,
		Cfg.Database.Name,
		Cfg.Database.SSLMode,
	)
}
