package postbox

import (
	"time"

	"github.com/nats-io/nats.go"
)

func NewOutbox(conn *nats.Conn) *Outbox {
	return &Outbox{conn: conn}
}

type Outbox struct {
	conn *nats.Conn
}

func (o Outbox) Send(subject string, env *Envelope) error {
	env.TS = time.Now().UTC()
	data, err := writeEnvelope(env)
	if err != nil {
		panic(err) // TODO
	}
	return o.conn.Publish(subject, data)
}

func (o Outbox) Request(subject string, env *Envelope) (*Envelope, error) {
	env.TS = time.Now().UTC()
	data, err := writeEnvelope(env)
	if err != nil {
		return nil, err
	}
	resp, err := o.conn.Request(subject, data, 10*time.Second) // TODO timeout
	if err != nil {
		return nil, err
	}
	return readEnvelope(resp.Data)
}
