package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Env      string   `yaml:"env"`
	HTTP     HTTP     `yaml:"http"`
	Telegram Telegram `yaml:"telegram"`
	Database Database `yaml:"database"`
	BaseURL  string   `yaml:"base_url"`
}

type HTTP struct {
	Port string `yaml:"port"`
}

type Telegram struct {
	BotToken    string `yaml:"bot_token"`
	BotUsername string `yaml:"bot_username"`
	ChannelID   int64  `yaml:"channel_id"`
}

type Database struct {
	DSN string `yaml:"dsn"`
}

func Load(path string) (*Config, error) {
	var cfg Config

	// Попытка загрузить из YAML, если файл существует
	if _, err := os.Stat(path); err == nil {
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("read config file: %w", err)
		}
		if err := yaml.Unmarshal(data, &cfg); err != nil {
			return nil, fmt.Errorf("parse config file: %w", err)
		}
	}

	// Переопределение из переменных окружения
	if env := os.Getenv("ENV"); env != "" {
		cfg.Env = env
	}
	if port := os.Getenv("PORT"); port != "" {
		if port[0] != ':' {
			cfg.HTTP.Port = ":" + port
		} else {
			cfg.HTTP.Port = port
		}
	}
	if token := os.Getenv("BOT_TOKEN"); token != "" {
		cfg.Telegram.BotToken = token
	}
	if dsn := os.Getenv("DATABASE_URL"); dsn != "" {
		cfg.Database.DSN = dsn
	}
	if baseURL := os.Getenv("BASE_URL"); baseURL != "" {
		cfg.BaseURL = baseURL
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func (c *Config) validate() error {
	if c.Env == "" {
		return fmt.Errorf("config error: env is empty")
	}

	if c.HTTP.Port == "" {
		return fmt.Errorf("config error: http.port is empty")
	}

	if c.Telegram.BotToken == "" {
		return fmt.Errorf("config error: telegram.bot_token is empty")
	}

	if c.Database.DSN == "" {
		return fmt.Errorf("config error: database.dsn is empty")
	}

	if c.BaseURL == "" {
		return fmt.Errorf("config error: base_url is empty")
	}

	return nil
}
