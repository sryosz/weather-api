package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"time"
)

type Config struct {
	EmailFrom    string        `yaml:"from"`
	EmailTo      string        `yaml:"to"`
	Password     string        `yaml:"password"`
	Subject      string        `yaml:"subject"`
	PollInterval time.Duration `yaml:"poll_interval"`
	Endpoint     string        `yaml:"endpoint"`
	Latitude     float64       `yaml:"latitude"`
	Longitude    float64       `yaml:"longitude"`
	Host         string        `yaml:"host"`
	Port         string        `yaml:"port"`
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		panic("config path is empty")
	}

	if _, err := os.Stat(configPath); err != nil {
		panic("config file does not exist: " + configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("cannot read config: " + err.Error())
	}

	return &cfg
}
