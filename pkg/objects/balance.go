package objects

type AccrualBalance struct {
	Balance  float32 `json:"current"`
	Withdraw float32 `json:"withdraw"`
}
