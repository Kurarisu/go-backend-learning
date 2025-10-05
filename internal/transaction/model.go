package transaction

import "time"

// Representasi tabel transactions
type Transaction struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Amount    float64   `json:"amount"`
	Type      string    `json:"type"` // deposit / withdraw
	CreatedAt time.Time `json:"created_at"`
}
