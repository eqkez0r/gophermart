package objects

import "time"

type Withdraw struct {
	WithdrawID  uint64    `json:"-"`
	UserID      uint64    `json:"-"`
	Order       string    `json:"order"`
	Sum         float32   `json:"sum"`
	ProcessedAt time.Time `json:"processed_at"`
}
