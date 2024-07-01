package objects

import "time"

const (
	New        = "NEW"
	Processing = "PROCESSING"
	Invalid    = "INVALID"   //End value
	Processed  = "PROCESSED" //End value
)

type Order struct {
	UserID   uint64    `json:"user_id"`
	Status   string    `json:"status"`
	UploadAt time.Time `json:"upload_at" json:"upload_at"`
	Number   *uint64   `json:"number" json:"number,omitempty"`
	Accrual  *float64  `json:"accrual,omitempty"`
}
