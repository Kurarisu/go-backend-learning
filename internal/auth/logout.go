// internal/auth/logout.go
package auth

import (
	"encoding/json"
	"net/http"
)

func LogoutHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Stateless JWT → logout cukup di client.
		// Di sini kita cuma balas sukses.
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Logout berhasil, silakan hapus token di client.",
		})
	}
}
