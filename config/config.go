package config

import (
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

var LogLevels = map[string]int{
	"DEBUG": -4,
	"INFO":  0,
	"WARN":  1,
	"ERROR": 8,
}

type Config struct {
	HTTPServer
	DataBase
	Logger
}

type Logger struct {
	Level string `env:"LOG_LEVEL"`
}

type HTTPServer struct {
	Ip   string `env:"HTTP_IP"`
	Port string `env:"HTTP_PORT"`
}

type DataBase struct {
	Username string `env:"DB_USER"`
	Password string `env:"DB_PASSWORD"`
	Host     string `env:"DB_HOST"`
	Port     string `env:"DB_PORT"`
	Database string `env:"DB_NAME"`
}

func GetConfig() *Config {
	var cfg Config

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		fmt.Printf("environment is not OK: %s\n", err)
		os.Exit(1)
	}

	return &cfg
}
