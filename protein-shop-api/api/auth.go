package api

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"uniproj/db"
	"uniproj/models"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

type registerRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"Invalid request"}`, http.StatusBadRequest)
		return
	}

	hashedPassword, err := HashPassword(req.Password)
	if err != nil {
		http.Error(w, `{"error":"Failed to hash password"}`, http.StatusInternalServerError)
		return
	}

	userID := uuid.New()
	_, err = db.Conn.Exec(context.Background(),
		"INSERT INTO users (id, email, password_hash) VALUES ($1, $2, $3)",
		userID, req.Email, hashedPassword)
	if err != nil {
		http.Error(w, `{"error":"Email already exists"}`, http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "User registered successfully",
		"user_id": userID.String(),
	})
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"Invalid request"}`, http.StatusBadRequest)
		return
	}

	var user models.User
	err := db.Conn.QueryRow(context.Background(),
		"SELECT id, email, password_hash FROM users WHERE email = $1", req.Email).
		Scan(&user.ID, &user.Email, &user.PasswordHash)
	if err == pgx.ErrNoRows {
		http.Error(w, `{"error":"Invalid credentials"}`, http.StatusUnauthorized)
		return
	}
	if err != nil {
		http.Error(w, `{"error":"Server error"}`, http.StatusInternalServerError)
		return
	}

	valid, err := VerifyPassword(user.PasswordHash, req.Password)
	if err != nil || !valid {
		http.Error(w, `{"error":"Invalid credentials"}`, http.StatusUnauthorized)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID.String(),
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		http.Error(w, `{"error":"Failed to generate token"}`, http.StatusInternalServerError)
		return
	}

	_, err = db.Conn.Exec(context.Background(),
		"INSERT INTO sessions (id, user_id, token, expires_at) VALUES ($1, $2, $3, $4)",
		uuid.New(), user.ID, tokenString, time.Now().Add(24*time.Hour))
	if err != nil {
		http.Error(w, `{"error":"Failed to store session"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"token":      tokenString,
		"expires_at": time.Now().Add(24 * time.Hour).Format(time.RFC3339),
	})
}
