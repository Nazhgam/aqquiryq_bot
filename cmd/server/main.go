package main

import (
	"context"
	"fmt"
	"log"

	"github.com/Nazhgam/aqquiryq_bot/internal/bot"
	"github.com/Nazhgam/aqquiryq_bot/internal/config"
	"github.com/Nazhgam/aqquiryq_bot/internal/repo"
	"github.com/Nazhgam/aqquiryq_bot/internal/service"
	"github.com/Nazhgam/aqquiryq_bot/internal/storage"
	"github.com/Nazhgam/aqquiryq_bot/migrations"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	cfg, err := config.Load("config/config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	if err := Run(ctx, cfg); err != nil {
		log.Fatal(err)
	}
}

func Run(ctx context.Context, cfg *config.Config) error {
	// 1️⃣ DB
	storage, err := storage.New(ctx, cfg.Database.DSN)
	if err != nil {
		return fmt.Errorf("db init failed: %w", err)
	}

	// 2️⃣ Migrations
	if err := migrations.Migration(storage.Pool); err != nil {
		return fmt.Errorf("migration failed: %w", err)
	}

	// 3️⃣ Repositories
	userRepo := repo.NewUserRepository(storage.Pool)
	contentRepo := repo.NewContentRepository(storage.Pool)

	// 4️⃣ Services
	userService := service.NewUserService(userRepo)
	contentService := service.NewContentService(contentRepo)

	// 5️⃣ Telegram API
	api, err := tgbotapi.NewBotAPI(cfg.Telegram.BotToken)
	if err != nil {
		return fmt.Errorf("telegram init failed: %w", err)
	}

	log.Printf("Authorized on account %s", api.Self.UserName)

	// 6️⃣ bot
	telegramBot := bot.New(api, userService, contentService)

	go telegramBot.Start(ctx)

	log.Println("Application started")

	select {} // блокируем main
}
