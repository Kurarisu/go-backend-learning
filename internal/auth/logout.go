// internal/auth/logout.go
package auth

import (
	"net/http"

	"go-backend/internal/response"
)

func LogoutHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Stateless JWT → logout cukup di client.
		// Server hanya mengembalikan response sukses.
		response.Success(w, map[string]string{
			"message": "Logout berhasil, silakan hapus token di client.",
		})
	}
}
