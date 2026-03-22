package config

import (
	"fmt"
	"os"
	"strconv"
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

func Load() (*Config, error) {
	botToken := os.Getenv("BOT_TOKEN")
	dsn := os.Getenv("DATABASE_URL")
	channelIDStr := os.Getenv("CHANNEL_ID")
	botUsername := os.Getenv("BOT_USERNAME")

	channelID, err := strconv.ParseInt(channelIDStr, 10, 64)
	if err != nil {
		return nil, err
	}

	cfg := &Config{
		Env: "prod",
		Telegram: Telegram{
			BotToken:    botToken,
			BotUsername: botUsername,
			ChannelID:   channelID,
		},
		Database: Database{DSN: dsn},
	}

	cfg.Telegram.BotToken = botToken
	cfg.Database.DSN = dsn

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) validate() error {
	if c.Telegram.BotToken == "" {
		return fmt.Errorf("config error: telegram.bot_token is empty")
	}

	if c.Database.DSN == "" {
		return fmt.Errorf("config error: database.dsn is empty")
	}

	return nil
}
