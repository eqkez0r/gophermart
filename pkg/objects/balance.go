package objects

type AccrualBalance struct {
	Balance  float64 `json:"sum"`
	Withdraw float64 `json:"withdraw"`
}
