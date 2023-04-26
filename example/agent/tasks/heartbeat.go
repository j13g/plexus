package tasks

import (
	"context"
	"time"

	"github.com/j13g/goutil/cron"
	"github.com/j13g/goutil/log"
	"github.com/j13g/plexus/example/payloads"
	"github.com/j13g/plexus/config"
	"github.com/j13g/plexus/postbox"
)

func GetHeartbeatTask() *cron.Task {
	cfg := config.Get()
	l := log.Get()

	return cron.NewTask("heartbeat", 10*time.Second, func(ctx context.Context) {
		payload := payloads.GetHeartbeat()
		l.Trace().Interface("heartbeat", payload).Msg("sending heartbeat")

		env, err := postbox.NewEnvelope(payloads.GetHeartbeat())
		if err != nil {
			panic(err) // TODO
		}

		err = cfg.Outbox.Send("heartbeat", env)
		if err != nil {
			panic(err) // TODO
		}
	})
}
