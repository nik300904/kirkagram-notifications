package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"time"
)

type Config struct {
	Env         string    `yaml:"env" env-required:"true"`
	StoragePath string    `yaml:"storage_path" env-required:"true"`
	HttpServe   HttpServe `yaml:"http_serve" env-required:"true"`
}

type HttpServe struct {
	Address     string        `yaml:"address" env-default:"8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"5s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"5s"`
}

func New() *Config {
	path, exists := os.LookupEnv("CONFIG_PATH")

	if !exists {
		panic(fmt.Sprintf("Переменная окружения не установлена CONFIG_PATH"))
	}

	var cfg Config

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic(err)
	}

	return &Config{
		Env:         cfg.Env,
		StoragePath: cfg.StoragePath,
		HttpServe: HttpServe{
			Address:     cfg.HttpServe.Address,
			Timeout:     cfg.HttpServe.Timeout,
			IdleTimeout: cfg.HttpServe.IdleTimeout,
		},
	}
}
