package payloads

type EchoPayload struct {
	Msg string `json:"msg"`
}

func (e EchoPayload) GetName() string {
	return "Echo"
}

func (e EchoPayload) GetVersion() string {
	return "1.0.0"
}
