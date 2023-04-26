package config

import (
	"fmt"
	"github.com/j13g/goutil/cli"
	"github.com/j13g/goutil/cron"
	"github.com/j13g/goutil/env"
	"github.com/j13g/goutil/types"
	"github.com/j13g/plexus/postbox"
	"github.com/nats-io/nats.go"
)

var Version = "UNKNOWN"

type Config struct {
	AppName    string `json:"app_name"`
	AppVersion string `json:"app_version"`

	NodeID   string `json:"node_id"`
	NodeArea string `json:"node_area"`

	NatsURL  string     `json:"nats_url"`
	NatsConn *nats.Conn `json:"-"`

	LogLevel   string `json:"log_level"`
	LogConsole bool   `json:"log_console"`

	CLIOutputFormat string `json:"cli_output_format"`

	// what is dependency injection anyway? basically just configurable global variables, right?
	CLI    *cli.CLI        `json:"-"`
	Outbox *postbox.Outbox `json:"-"`
	Inbox  *postbox.Inbox  `json:"-"`
	Cron   *cron.Container `json:"-"`
}

func (c *Config) NodeName() string {
	return fmt.Sprintf("%s.%s.%s", c.AppName, c.NodeArea, c.NodeID)
}

func (c *Config) LoadFromEnv(prefix string) {
	env.SetEnvPrefix(prefix)
	c.AppVersion = Version

	c.NodeID = env.GetStringDefault("NODE_ID", env.NodeID())
	c.NodeArea = env.GetStringDefault("NODE_AREA", "unknown")

	c.LogLevel = env.GetStringDefault("LOG_LEVEL", "debug")
	c.LogConsole = env.GetBool("LOG_CONSOLE").OrElse(false)

	c.CLIOutputFormat = env.GetStringDefault("OUTPUT_FORMAT", "json")

	if x := env.GetString("NATS_URL"); x.IsPresent() {
		c.NatsURL = x.MustGet()
	} else if x := env.GetString("NATS_URL", env.WithPrefix("")); x.IsPresent() {
		c.NatsURL = x.MustGet()
	} else {
		c.NatsURL = "nats://localhost:4222"
	}
}

var Get = types.Singleton[Config]()
