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

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterResponse struct {
	Token string `json:"token"`
}

func RegisterHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RegisterRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.Error(w, http.StatusBadRequest, "Invalid request")
			return
		}

		// Hash password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			response.Error(w, http.StatusInternalServerError, "Failed to hash password")
			return
		}

		// Insert user
		var userID int
		err = db.QueryRow(
			"INSERT INTO users (name, email, password, balance) VALUES ($1, $2, $3, 0) RETURNING id",
			req.Name, req.Email, string(hashedPassword),
		).Scan(&userID)
		if err != nil {
			response.Error(w, http.StatusInternalServerError, "Failed to create user")
			return
		}

		// Generate JWT
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"user_id": userID,
			"exp":     time.Now().Add(time.Hour).Unix(),
		})

		tokenString, err := token.SignedString(GetJWTSecret())
		if err != nil {
			response.Error(w, http.StatusInternalServerError, "Failed to generate token")
			return
		}

		response.Success(w, RegisterResponse{Token: tokenString})
	}
}
