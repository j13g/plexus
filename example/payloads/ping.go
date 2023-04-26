package payloads

import (
	"time"
)

type PingPayload struct {
	TS time.Time `json:"ts"`
}

func (p PingPayload) GetName() string {
	return "Ping"
}

func (p PingPayload) GetVersion() string {
	return "1.0.0"
}
