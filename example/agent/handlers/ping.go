package handlers

import (
	"context"
	"time"

	"github.com/j13g/plexus/example/payloads"
	"github.com/j13g/plexus/postbox"
)

func PingHandler(ctx context.Context, envelope *postbox.Envelope) *postbox.Envelope {
	env, err := postbox.NewEnvelope(payloads.PingPayload{
		TS: time.Now().UTC(),
	})
	if err != nil {
		panic(err) // TODO
	}
	return env
}
