package objects

type AccrualBalance struct {
	Balance  float64 `json:"current"`
	Withdraw float64 `json:"withdraw"`
}
