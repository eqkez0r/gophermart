package objects

type User struct {
	UserID         uint64 `json:"user_id,omitempty"`
	Login          string `json:"login"`
	Password       string `json:"password"`
	AccrualBalance `json:"accrual_balance"`
}
