package postbox

import (
	"encoding/json"
	"time"
)

type Payload interface {
	GetName() string
	GetVersion() string
}

type Envelope struct {
	Name    string          `json:"n"`
	Version string          `json:"v"`
	TS      time.Time       `json:"t"`
	Error   string          `json:"e,omitempty"`
	Meta    map[string]any  `json:"m,omitempty"`
	Payload json.RawMessage `json:"p,omitempty"`
}

func NewErrorEnvelope(err error) *Envelope {
	return &Envelope{
		Name:    "Error",
		Version: "1.0.0",
		TS:      time.Now().UTC(),
		Error:   err.Error(),
	}
}

func NewEnvelope[T Payload](payload T) (*Envelope, error) {
	p, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	return &Envelope{
		Name:    payload.GetName(),
		Version: payload.GetVersion(),
		Payload: p,
		Meta:    make(map[string]any),
	}, nil
}

func GetPayload[T Payload](e *Envelope) (T, error) {
	var p T
	err := json.Unmarshal(e.Payload, &p)
	return p, err
}

func writeEnvelope(envelope *Envelope) ([]byte, error) {
	return json.Marshal(envelope)
}

func readEnvelope(data []byte) (*Envelope, error) {
	x := &Envelope{
		Meta: make(map[string]any),
	}
	err := json.Unmarshal(data, x)
	return x, err
}
