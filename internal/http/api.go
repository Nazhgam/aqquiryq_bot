package http

import (
	"context"
	"html/template"
	"log"
	"net/http"

	"github.com/Nazhgam/aqquiryq_bot/internal/bot"
	"github.com/Nazhgam/aqquiryq_bot/internal/config"
	"github.com/Nazhgam/aqquiryq_bot/internal/http/handlers"
	"github.com/Nazhgam/aqquiryq_bot/internal/repo"
	"github.com/gorilla/mux"
)

func Start(ctx context.Context, cfg *config.Config, tgBot *bot.Bot, userRepo repo.UserRepository, contentRepo repo.ContentRepository) {
	templates, err := template.ParseGlob("templates/*.html")
	if err != nil {
		log.Fatalf("Error parsing templates: %v", err)
	}

	handlers := handlers.New(cfg, templates, userRepo, contentRepo, tgBot)

	r := mux.NewRouter()

	// Serve static files
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	r.HandleFunc("/", handlers.Dashboard).Methods("GET")
	r.HandleFunc("/login", handlers.Login).Methods("GET")
	r.HandleFunc("/logout", handlers.Logout).Methods("GET")
	r.HandleFunc("/auth/telegram", handlers.TelegramAuth).Methods("GET")
	//r.HandleFunc("/view/{id}", viewHandler).Methods("GET")
	//r.HandleFunc("/post/{id}", postHandler).Methods("POST") // Use POST for actions

	log.Printf("Server starting on %s", cfg.HTTP.Port)
	err = http.ListenAndServe(cfg.HTTP.Port, r)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
