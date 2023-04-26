package handlers

import (
	"context"

	"github.com/j13g/plexus/example/payloads"
	"github.com/j13g/plexus/postbox"
)

func EchoHandler(ctx context.Context, envelope *postbox.Envelope) *postbox.Envelope {
	payload, err := postbox.GetPayload[payloads.EchoPayload](envelope)
	if err != nil {
		panic(err) // TODO
	}
	msg := payload.Msg

	env, err := postbox.NewEnvelope(payloads.EchoPayload{
		Msg: msg,
	})
	if err != nil {
		panic(err) // TODO
	}
	return env
}
