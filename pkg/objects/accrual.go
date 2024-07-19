package objects

const (
	AccrualStatusRegistered = "REGISTERED"
	AccrualStatusProcessing = "PROCESSING"
	AccrualStatusInvalid    = "INVALID"
	AccrualStatusProcessed  = "PROCESSED"
)

var AccrualStatusToOrderStatus = map[string]string{
	AccrualStatusRegistered: OrderStatusProcessing,
	AccrualStatusProcessing: OrderStatusProcessing,
	AccrualStatusInvalid:    OrderStatusInvalid,
	AccrualStatusProcessed:  OrderStatusProcessed,
}

type Accrual struct {
	Order   string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float32 `json:"accrual,omitempty"`
}
