package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"go-backend/internal/user"

	_ "github.com/lib/pq"
)

func main() {
	// Koneksi ke Postgres
	connStr := "postgres://go_user:go_pass@localhost:5432/go_backend?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("❌ Gagal konek ke database:", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("❌ Database tidak bisa diakses:", err)
	}
	fmt.Println("✅ Berhasil konek ke database!")

	// Routing sederhana
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

	// Jalankan server
	fmt.Println("🚀 Server jalan di :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
