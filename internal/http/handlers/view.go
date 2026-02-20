package handlers

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/gorilla/mux"
)

func (h *handler) View(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	// Get user from cookie
	cookie, err := r.Cookie("user_id")
	if err != nil {
		// Deep link redirection: Send them to login, then bring them back here
		returnURL := "/view/" + id
		http.Redirect(w, r, "/login?redirect_to="+url.QueryEscape(returnURL), http.StatusSeeOther)
		return
	}

	userIDStr := cookie.Value
	var userID int64
	fmt.Sscanf(userIDStr, "%d", &userID)

	allowed, err := h.tgBot.IsUserMember(userID, h.cfg.Telegram.ChannelID)
	if err != nil {
		log.Printf("Error checking membership: %v", err)
		http.Error(w, "Error checking permissions", http.StatusInternalServerError)
		return
	}

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

	// Prevent caching so permissions are checked every time
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	if !allowed {
		http.Error(w, "You are not a member of the required channel to view this presentation.", http.StatusForbidden)
		return
	}

	h.templates.ExecuteTemplate(w, "viewer.html", map[string]string{
		"Title":    presentation.Title,
		"CanvaURL": presentation.CanvaURL,
	})
}
