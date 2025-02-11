package config

import "time"

type Config struct {
	Server struct {
		Host         string        `yaml:"host"`
		Port         string        `yaml:"port"`
		ResponseTime time.Duration `yaml:"response_time"`
		RPS          int           `yaml:"rps"`
	} `yaml:"server"`
	DB struct {
		URL string `yaml:"url"`
	} `yaml:"db"`
}

// func New(path string) *Config {

// }
