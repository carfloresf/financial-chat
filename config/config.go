package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		HTTP     `yaml:"http"`
		DB       `yaml:"db"`
		Auth     `yaml:"auth"`
		RabbitMQ `yaml:"rabbitmq"`
	}

	HTTP struct {
		Addr string `env-required:"true" yaml:"address" env:"HTTP_ADDR"`
		Port string `env-required:"true" yaml:"port" env:"HTTP_PORT"`
	}

	DB struct {
		DBFile string `env-required:"true" yaml:"db_file" env:"DB_FILE"`
	}

	Auth struct {
		Pepper string `env-required:"true" yaml:"pepper" env:"PEPPER"`
		Cookie string `env-required:"true" yaml:"cookie" env:"COOKIE"`
	}

	RabbitMQ struct {
		URL string `env-required:"true" yaml:"url" env:"RABBITMQ_URL"`
	}
)

// NewConfig returns app config.
func NewConfig(path string) (*Config, error) {
	cfg := &Config{}

	err := cleanenv.ReadConfig(path, cfg)
	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
