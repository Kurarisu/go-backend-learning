package transaction

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"go-backend/internal/response"
)

func DepositHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input struct {
			UserID int     `json:"user_id"`
			Amount float64 `json:"amount"`
		}

		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			response.Error(w, http.StatusBadRequest, "Invalid request")
			return
		}

		// Insert transaksi
		_, err := db.Exec(
			"INSERT INTO transactions (user_id, amount, type) VALUES ($1, $2, 'deposit')",
			input.UserID, input.Amount,
		)
		if err != nil {
			response.Error(w, http.StatusInternalServerError, "Failed to insert transaction")
			return
		}

		// Update saldo user
		_, err = db.Exec(
			"UPDATE users SET balance = balance + $1 WHERE id = $2",
			input.Amount, input.UserID,
		)
		if err != nil {
			response.Error(w, http.StatusInternalServerError, "Failed to update balance")
			return
		}

		response.Success(w, map[string]string{"message": "deposit success"})
	}
}

// WithdrawHandler untuk tarik uang
func WithdrawHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input struct {
			UserID int     `json:"user_id"`
			Amount float64 `json:"amount"`
		}

		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			response.Error(w, http.StatusBadRequest, "Invalid request")
			return
		}

		// Cek saldo dulu
		var balance float64
		err := db.QueryRow("SELECT balance FROM users WHERE id = $1", input.UserID).Scan(&balance)
		if err != nil {
			response.Error(w, http.StatusNotFound, "User not found")
			return
		}
		if balance < input.Amount {
			response.Error(w, http.StatusBadRequest, "Insufficient balance")
			return
		}

		// Insert transaksi
		_, err = db.Exec(
			"INSERT INTO transactions (user_id, amount, type) VALUES ($1, $2, 'withdraw')",
			input.UserID, input.Amount,
		)
		if err != nil {
			response.Error(w, http.StatusInternalServerError, "Failed to insert transaction")
			return
		}

		// Update saldo user
		_, err = db.Exec(
			"UPDATE users SET balance = balance - $1 WHERE id = $2",
			input.Amount, input.UserID,
		)
		if err != nil {
			response.Error(w, http.StatusInternalServerError, "Failed to update balance")
			return
		}

		response.Success(w, map[string]string{"message": "withdraw success"})
	}
}

// ListTransactionsHandler untuk melihat riwayat transaksi user
func ListTransactionsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.URL.Query().Get("user_id")
		if userID == "" {
			response.Error(w, http.StatusBadRequest, "user_id is required")
			return
		}

		rows, err := db.Query(
			"SELECT id, user_id, amount, type, created_at FROM transactions WHERE user_id = $1 ORDER BY created_at DESC",
			userID,
		)
		if err != nil {
			response.Error(w, http.StatusInternalServerError, "Failed to query transactions")
			return
		}
		defer rows.Close()

		var transactions []Transaction
		for rows.Next() {
			var t Transaction
			if err := rows.Scan(&t.ID, &t.UserID, &t.Amount, &t.Type, &t.CreatedAt); err != nil {
				response.Error(w, http.StatusInternalServerError, "Failed to scan transaction")
				return
			}
			transactions = append(transactions, t)
		}

		response.Success(w, transactions)
	}
}
