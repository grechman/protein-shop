package api

import (
	"context"
	"encoding/json"
	"net/http"

	"uniproj/db"
	"uniproj/models"
)

func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		http.Error(w, `{"error":"Unauthorized"}`, http.StatusUnauthorized)
		return
	}

	var user models.User
	err := db.Conn.QueryRow(context.Background(),
		"SELECT id, email, created_at FROM users WHERE id = $1", userID).
		Scan(&user.ID, &user.Email, &user.CreatedAt)
	if err != nil {
		http.Error(w, `{"error":"Failed to fetch profile"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
