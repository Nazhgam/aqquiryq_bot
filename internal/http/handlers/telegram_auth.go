package handlers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

func (h *handler) TelegramAuth(w http.ResponseWriter, r *http.Request) {
	user, err := checkTelegramAuth(r.URL.Query(), h.cfg.Telegram.BotToken)
	if err != nil {
		http.Error(w, "Authentication failed: "+err.Error(), http.StatusUnauthorized)
		return
	}

	// Set a cookie or session here for the user
	http.SetCookie(w, &http.Cookie{
		Name:     "user_id",
		Value:    user["id"],
		Path:     "/",
		HttpOnly: true,
		Secure:   true,                  // Required for SameSite=None
		SameSite: http.SameSiteNoneMode, // Allows cross-site cookie for Telegram redirect
		MaxAge:   3600 * 24 * 30,        // 30 days
	})

	// Redirect to home or dashboard (or the requested page)
	redirectTo := r.URL.Query().Get("redirect_to")
	if redirectTo != "" {
		http.Redirect(w, r, redirectTo, http.StatusSeeOther)
	} else {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func checkTelegramAuth(query url.Values, token string) (map[string]string, error) {
	// Check if hash is present
	hash := query.Get("hash")
	if hash == "" {
		return nil, errors.New("hash is missing")
	}

	// Create data-check-string
	var args []string
	for k, v := range query {
		if k == "hash" || k == "redirect_to" {
			continue
		}
		args = append(args, fmt.Sprintf("%s=%s", k, v[0]))
	}
	sort.Strings(args)
	dataCheckString := strings.Join(args, "\n")

	// Compute secret key
	sha256Hash := sha256.New()
	sha256Hash.Write([]byte(token))
	secretKey := sha256Hash.Sum(nil)

	// Compute HMAC-SHA256 signature
	hmacHash := hmac.New(sha256.New, secretKey)
	hmacHash.Write([]byte(dataCheckString))
	signature := hex.EncodeToString(hmacHash.Sum(nil))

	// DEBUG LOGGING
	if signature != hash {
		log.Printf("Auth FAILED. Expected: %s, Got: %s", signature, hash)
		log.Printf("Data Check String was:\n%s", dataCheckString)
	}

	// Compare signatures
	if signature != hash {
		return nil, errors.New("signature mismatch")
	}

	// Check auth_date
	authDateStr := query.Get("auth_date")
	var authDate int64
	fmt.Sscanf(authDateStr, "%d", &authDate)
	if time.Now().Unix()-authDate > 86400 {
		return nil, errors.New("auth data is outdated")
	}

	user := make(map[string]string)
	user["id"] = query.Get("id")
	user["first_name"] = query.Get("first_name")
	user["username"] = query.Get("username")
	return user, nil
}
