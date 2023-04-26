package outbox

import (
	"github.com/j13g/goutil/cli"
	"github.com/j13g/plexus/config"
	"github.com/j13g/plexus/postbox"
	"github.com/spf13/cobra"
	"time"
)

func OutboxCmd(cfg *config.Config) {
	OutboxSendCmd(cfg)
	OutboxRequestCmd(cfg)
}

func OutboxRequestCmd(cfg *config.Config) {
	outboxRequest := &cobra.Command{
		Aliases: []string{"r", "req"},
		RunE: func(cmd *cobra.Command, args []string) error {

			name, err := cmd.Flags().GetString("name")
			if err != nil {
				return err
			}

			version, err := cmd.Flags().GetString("version")
			if err != nil {
				return err
			}

			payload, err := cmd.Flags().GetString("payload")
			if err != nil {
				return err
			}

			subject, err := cmd.Flags().GetString("subject")
			if err != nil {
				return err
			}

			e := &postbox.Envelope{
				Name:    name,
				Version: version,
				TS:      time.Now().UTC(),
				Payload: []byte(payload),
			}
			outbox := postbox.NewOutbox(cfg.NatsConn)
			resp, err := outbox.Request(subject, e)
			if err != nil {
				return err
			}

			return cli.Print(resp)
		},
	}
	outboxRequest.Flags().StringP("subject", "s", "", "subject")
	outboxRequest.Flags().StringP("name", "n", "", "payload name")
	outboxRequest.Flags().StringP("version", "v", "", "payload version")
	outboxRequest.Flags().StringP("payload", "p", "", "envelope payload")

	for _, x := range []string{"subject", "name", "version", "payload"} {
		outboxRequest.MarkFlagRequired(x)
	}

	cfg.CLI.Add("outbox.request", outboxRequest)
}

func OutboxSendCmd(cfg *config.Config) {
	outboxSend := &cobra.Command{
		RunE: func(cmd *cobra.Command, args []string) error {

			name, err := cmd.Flags().GetString("name")
			if err != nil {
				return err
			}

			version, err := cmd.Flags().GetString("version")
			if err != nil {
				return err
			}

			payload, err := cmd.Flags().GetString("payload")
			if err != nil {
				return err
			}

			subject, err := cmd.Flags().GetString("subject")
			if err != nil {
				return err
			}

			e := &postbox.Envelope{
				Name:    name,
				Version: version,
				TS:      time.Now().UTC(),
				Payload: []byte(payload),
			}
			outbox := postbox.NewOutbox(cfg.NatsConn)
			return outbox.Send(subject, e)
		},
	}
	outboxSend.Flags().StringP("subject", "s", "", "subject")
	outboxSend.Flags().StringP("name", "n", "", "payload name")
	outboxSend.Flags().StringP("version", "v", "", "payload version")
	outboxSend.Flags().StringP("payload", "p", "", "envelope payload")

	for _, x := range []string{"subject", "name", "version", "payload"} {
		outboxSend.MarkFlagRequired(x)
	}

	cfg.CLI.Add("outbox.send", outboxSend)
}
