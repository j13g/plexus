package mainutil

import (
	"github.com/j13g/plexus/config"
	"github.com/j13g/plexus/natsutil"
	"github.com/j13g/plexus/postbox"
	"strings"

	"github.com/j13g/goutil/cli"
	"github.com/j13g/goutil/cli/outputter"
	"github.com/j13g/goutil/log"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

type AppSetup struct {
	Name      string
	Version   string
	EnvPrefix string
	SetupFunc func()
}

var rootCmd = &cobra.Command{
	PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
		cfg := config.Get()
		err := cfg.NatsConn.Drain()
		if err != nil {
			return err
		}

		return nil
	},
}

func CLIMain(setup AppSetup) error {
	cfg := config.Get()
	cfg.LoadFromEnv(strings.ToUpper(setup.EnvPrefix))
	cfg.AppName = setup.Name
	cfg.AppVersion = setup.Version

	// Logging
	err := logSetup(cfg, setup)
	if err != nil {
		return err
	}
	l := log.Get()

	// NATs, Outbox, Inbox
	conn, err := natsutil.Connect()
	if err != nil {
		return err
	}
	cfg.NatsConn = conn

	config.ProvideValue[*postbox.Outbox](postbox.NewOutbox(cfg.NatsConn))
	config.ProvideValue[*postbox.Inbox](postbox.NewInbox(cfg.NatsConn, cfg.AppName, cfg.AppVersion, cfg.NodeName()))

	// CLI
	cliInstance := cli.NewCLI(setup.Name).SetRoot(rootCmd)
	config.ProvideValue[*cli.CLI](cliInstance)
	err = cli.SetOutputterByName(cfg.CLIOutputFormat)
	if err != nil {
		l := log.Get()
		l.Error().Err(err).Str("format", cfg.CLIOutputFormat).Msg("failed to set CLI output format")
		cli.SetOutputter(outputter.JSONOutputter{}) // using JSON output format as a fallback
	}

	// App-specific setup
	l.Trace().Msg("running setup")
	setup.SetupFunc()

	// run CLI
	l.Trace().Msg("running cli")
	err = config.Invoke[*cli.CLI]().MustGet().Run()
	if err != nil {
		l.Error().Err(err).Msg("failed to run")
		return err
	}
	return nil
}

func logSetup(cfg *config.Config, setup AppSetup) error {
	lvl, err := zerolog.ParseLevel(cfg.LogLevel)
	if err != nil {
		return err
	}
	log.WithLevel(lvl)

	options := []log.Option{
		log.WithAppName(setup.Name),
		log.WithStdout(),
		log.WithLevel(lvl),
	}

	if cfg.LogConsole {
		options = append(options, log.WithConsoleOutput())
	}

	log.SetupLogging(options...)
	return nil
}
