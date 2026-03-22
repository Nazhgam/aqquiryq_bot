package bot

import (
	"sync"

	"github.com/Nazhgam/aqquiryq_bot/internal/service"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	api            *tgbotapi.BotAPI
	userService    service.UserService
	contentService service.ContentService

	adminSessions map[int64]*AdminSession
	mu            sync.RWMutex
}

func New(
	api *tgbotapi.BotAPI,
	userService service.UserService,
	contentService service.ContentService,
) *Bot {
	return &Bot{
		api:            api,
		userService:    userService,
		contentService: contentService,
		adminSessions:  make(map[int64]*AdminSession),
	}
}
