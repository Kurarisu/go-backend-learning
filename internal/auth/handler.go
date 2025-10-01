package auth

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"go-backend/internal/response"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

func LoginHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.Error(w, http.StatusBadRequest, "Invalid request")
			return
		}

		// Ambil user dari DB
		var id int
		var hashedPassword string
		err := db.QueryRow("SELECT id, password FROM users WHERE email = $1", req.Email).
			Scan(&id, &hashedPassword)
		if err != nil {
			response.Error(w, http.StatusUnauthorized, "Invalid email or password")
			return
		}

		// Cek password
		if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(req.Password)); err != nil {
			response.Error(w, http.StatusUnauthorized, "Invalid email or password")
			return
		}

		// Generate JWT
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"user_id": id,
			"exp":     time.Now().Add(time.Hour).Unix(),
		})

		tokenString, err := token.SignedString(GetJWTSecret())
		if err != nil {
			response.Error(w, http.StatusInternalServerError, "Failed to generate token")
			return
		}

		// Sukses
		response.Success(w, LoginResponse{Token: tokenString})
	}
}
