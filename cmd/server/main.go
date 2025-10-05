package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"go-backend/internal/auth"
	"go-backend/internal/middleware"
	"go-backend/internal/transaction"
	"go-backend/internal/user"
)

// CORS
func enableCORSFlexible(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	// Load .env
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found or failed to load, using system env vars")
	}

	// Contoh ambil variabel
	dbUser := os.Getenv("DB_USER")
	log.Println("DB_USER =", dbUser)

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET tidak ditemukan di .env")
	}
	auth.SetJWTSecret(jwtSecret)

	// DB connection
	connStr := "postgres://go_user:go_pass@localhost:5432/go_backend?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Gagal konek DB:", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("DB tidak bisa di-ping:", err)
	}
	fmt.Println("✅ Berhasil konek ke database!")

	// pakai mux biar gampang wrap middleware
	mux := http.NewServeMux()

	// Users
	mux.Handle("/users", enableCORSFlexible(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			user.GetUsersHandler(db)(w, r)
			return
		}
		if r.Method == http.MethodPost {
			user.CreateUserHandler(db)(w, r)
			return
		}
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	})))

	// Profile (/me)
	mux.Handle("/me", enableCORSFlexible(auth.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value(auth.UserIDKey).(int)

		var name, email string
		err := db.QueryRow("SELECT name, email FROM users WHERE id = $1", userID).Scan(&name, &email)
		if err != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"id":    userID,
			"name":  name,
			"email": email,
		})
	}))))

	// Auth
	mux.Handle("/login", enableCORSFlexible(auth.LoginHandler(db)))
	mux.Handle("/logout", enableCORSFlexible(auth.LogoutHandler()))
	mux.Handle("/register", enableCORSFlexible(auth.RegisterHandler(db)))

	// Transactions
	mux.Handle("/deposit", enableCORSFlexible(auth.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			transaction.DepositHandler(db)(w, r)
			return
		}
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}))))

	mux.Handle("/withdraw", enableCORSFlexible(auth.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			transaction.WithdrawHandler(db)(w, r)
			return
		}
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}))))

	mux.Handle("/transactions", enableCORSFlexible(auth.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			transaction.ListTransactionsHandler(db)(w, r)
			return
		}
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}))))

	// Bungkus semua dengan LoggingMiddleware
	handler := middleware.LoggingMiddleware(mux)

	fmt.Println("🚀 Server jalan di :8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
