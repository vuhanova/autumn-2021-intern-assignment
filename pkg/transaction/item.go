package transaction

import "time"

type User struct {
	UserID   int     `json:"id"`
	Balance  float64 `json:"balance"`
	ToUserID int     `json:"id_to,omitempty"`
	Field    string  `json:"field,omitempty"`
}

type Transaction struct {
	ToID    *int      `json:"to_id"`
	FromID  *int      `json:"from_id"`
	Money   float64   `json:"money"`
	Created time.Time `json:"created"`
}
