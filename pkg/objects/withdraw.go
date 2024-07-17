package objects

import "time"

type Withdraw struct {
	UserID      uint64    `json:"userID"`
	Order       string    `json:"order"`
	Sum         float64   `json:"sum"`
	ProcessedAt time.Time `json:"processed_at"`
}
