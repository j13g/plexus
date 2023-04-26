package natsutil

import (
	"fmt"
	"github.com/j13g/plexus/config"

	"github.com/nats-io/nats.go"
)

func Connect() (*nats.Conn, error) {
	cfg := config.Get()

	conn, err := nats.Connect(
		cfg.NatsURL,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to nats: %w", err)
	}

	return conn, nil
}
