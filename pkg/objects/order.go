package objects

import "time"

const (
	OrderStatusNew        = "NEW"
	OrderStatusProcessing = "PROCESSING"
	OrderStatusInvalid    = "INVALID"   //End value
	OrderStatusProcessed  = "PROCESSED" //End value
)

type Order struct {
	UserID   uint64    `json:"user_id"`
	Status   string    `json:"status"`
	UploadAt time.Time `json:"upload_at"`
	Number   *uint64   `json:"number,omitempty"`
	Accrual  *float64  `json:"accrual,omitempty"`
}
