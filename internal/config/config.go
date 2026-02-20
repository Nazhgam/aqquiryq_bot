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
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse config file: %w", err)
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
