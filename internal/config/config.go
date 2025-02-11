package config

import (
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	SecretKey string `env:"SECRET" env-required:"true"`
	Server    struct {
		Host         string        `yaml:"host" env-default:"localhost"`
		Port         string        `yaml:"port" env-default:"8080"`
		ResponseTime time.Duration `yaml:"response_time" env-default:"50ms"`
		RPS          int           `yaml:"rps" env-default:"1000"`
	} `yaml:"server"`
	DB struct {
		URL string `yaml:"url" env-required:"true"`
	} `yaml:"db"`
}

func New(path string) (*Config, error) {
	var cfg *Config

	if err := cleanenv.ReadConfig(path, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
