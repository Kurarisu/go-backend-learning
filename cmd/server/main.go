package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"

	"go-backend/internal/auth"
	"go-backend/internal/user"
)

func main() {
	// Koneksi ke PostgreSQL
	connStr := "postgres://go_user:go_pass@localhost:5432/go_backend?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Gagal konek DB:", err)
	}
	defer db.Close()

	// Cek koneksi
	if err := db.Ping(); err != nil {
		log.Fatal("DB tidak bisa di-ping:", err)
	}
	fmt.Println("✅ Berhasil konek ke database!")

	// Routes
	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			user.GetUsersHandler(db)(w, r)
			return
		}
		if r.Method == http.MethodPost {
			user.CreateUserHandler(db)(w, r)
			return
		}
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	})

	http.HandleFunc("/me", auth.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
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
	}))

	http.HandleFunc("/login", auth.LoginHandler(db))

	http.HandleFunc("/logout", auth.LogoutHandler())

	// Start server
	fmt.Println("🚀 Server jalan di :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
