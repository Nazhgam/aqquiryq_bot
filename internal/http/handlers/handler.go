package handlers

import (
	"html/template"
	"net/http"

	"github.com/Nazhgam/aqquiryq_bot/internal/bot"
	"github.com/Nazhgam/aqquiryq_bot/internal/config"
	"github.com/Nazhgam/aqquiryq_bot/internal/repo"
)

type Handlers interface {
	Dashboard(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
	Logout(w http.ResponseWriter, r *http.Request)
	TelegramAuth(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	cfg         *config.Config
	userRepo    repo.UserRepository
	contentRepo repo.ContentRepository
	templates   *template.Template
	tgBot       *bot.Bot
}

func New(cfg *config.Config,
	templates *template.Template,
	userRepo repo.UserRepository,
	contentRepo repo.ContentRepository,
	tgBot *bot.Bot) Handlers {
	return &handler{
		cfg:         cfg,
		userRepo:    userRepo,
		contentRepo: contentRepo,
		templates:   templates,
		tgBot:       tgBot,
	}
}
