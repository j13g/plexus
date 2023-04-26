package main

import (
	"github.com/j13g/goutil/cron"
	"github.com/j13g/goutil/log"
	"github.com/j13g/goutil/sig"
	"github.com/j13g/plexus/example/agent/handlers"
	"github.com/j13g/plexus/example/agent/tasks"
	"github.com/j13g/plexus/config"
	"github.com/j13g/plexus/mainutil"
	"github.com/j13g/plexus/postbox"
	"github.com/spf13/cobra"
)

const appName = "plexus_agent"

func main() {
	mainutil.CLIMain(mainutil.AppSetup{
		Name:      appName,
		Version:   "UNKNOWN",
		SetupFunc: setup,
	})
}

func setup() {
	cfg := config.Get()
	cfg.CLI.Add("start", &cobra.Command{
		RunE: func(cmd *cobra.Command, args []string) error {
			l := log.Get()

			l.Info().Msg("starting scheduled tasks")
			StartScheduledTasks()
			defer func() {
				StopScheduledTasks()
				l.Info().Msg("scheduled tasks stopped")
			}()

			l.Info().Msg("starting message handlers")
			StartMessageHandlers()
			defer func() {
				StopMessageHandlers()
				l.Info().Msg("message handlers stopped")
			}()

			sig.WaitUntilSignalled()
			return nil
		},
	})
}

func StartScheduledTasks() {
	cfg := config.Get()
	cfg.Cron = cron.NewContainer()
	cfg.Cron.Add(tasks.GetHeartbeatTask())
	cfg.Cron.Start()
}

func StopScheduledTasks() {
	cfg := config.Get()
	cfg.Cron.Stop()
}

func StartMessageHandlers() {
	cfg := config.Get()
	router := cfg.Inbox.Router()
	router.Register("Ping", "1.x", handlers.PingHandler)
	router.Register("Echo", "1.x", handlers.EchoHandler)

	cfg.Inbox.Start(postbox.NewSubjectSpec().AddMultiF("%s.%s.%s", cfg.AppName, cfg.NodeArea, cfg.NodeID))
}

func StopMessageHandlers() {
	cfg := config.Get()
	cfg.Inbox.Stop()
}
