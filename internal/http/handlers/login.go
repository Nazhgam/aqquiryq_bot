package handlers

import (
	"net/http"
	"net/url"
)

func (h *handler) Login(w http.ResponseWriter, r *http.Request) {
	redirectTo := r.URL.Query().Get("redirect_to")
	authURL := h.cfg.BaseURL + "/auth/telegram"
	if redirectTo != "" {
		authURL += "?redirect_to=" + url.QueryEscape(redirectTo)
	}

	data := struct {
		BotUsername string
		AuthURL     string
	}{
		BotUsername: h.cfg.Telegram.BotUsername,
		AuthURL:     authURL,
	}

	h.templates.ExecuteTemplate(w, "login.html", data)
}
