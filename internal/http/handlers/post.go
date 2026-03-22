package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (h *handler) Post(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	// Security: In a real app, check if the requester is an Admin
	// Here we assume only admins access the dashboard.

	contentID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid presentation ID", http.StatusBadRequest)
		return
	}

	presentation, err := h.contentRepo.GetByID(r.Context(), int64(contentID))
	if err != nil {
		http.Error(w, "Presentation not found", http.StatusNotFound)
		return
	}

	// Construct the View URL
	viewURL := h.cfg.BaseURL + "/view/" + id

	// Send to Telegram
	err = h.tgBot.PostPresentationToChannel(h.cfg.Telegram.ChannelID, presentation.Title, viewURL)
	if err != nil {
		log.Printf("Error posting to channel: %v", err)
		http.Error(w, "Failed to post to Telegram: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Redirect back to dashboard with success message (or just back)
	http.Redirect(w, r, "/?posted=true", http.StatusSeeOther)
}
