package middleware

import (
	"context"
	"net/http"
	"os"

	"uniproj/db"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" || len(tokenString) < 7 || tokenString[:7] != "Bearer " {
			http.Error(w, `{"error":"Unauthorized"}`, http.StatusUnauthorized)
			return
		}
		tokenString = tokenString[7:]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})
		if err != nil || !token.Valid {
			http.Error(w, `{"error":"Invalid token"}`, http.StatusUnauthorized)
			return
		}

		var exists bool
		err = db.Conn.QueryRow(context.Background(),
			"SELECT EXISTS(SELECT 1 FROM sessions WHERE token = $1 AND expires_at > NOW())",
			tokenString).Scan(&exists)
		if err != nil || !exists {
			http.Error(w, `{"error":"Session expired or invalid"}`, http.StatusUnauthorized)
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		ctx := context.WithValue(r.Context(), "user_id", claims["user_id"])
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
