package enaptercli

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v2"

	"github.com/enapter/enapter-cli/internal/cloudapi"
)

type cmdRulesLogs struct {
	cmdRules
}

func buildCmdRulesLogs() *cli.Command {
	cmd := &cmdRulesLogs{}

	return &cli.Command{
		Name:               "logs",
		Usage:              "Stream logs from a rule",
		CustomHelpTemplate: cmd.HelpTemplate(),
		Flags:              cmd.Flags(),
		Before:             cmd.Before,
		Action: func(cliCtx *cli.Context) error {
			return cmd.run(cliCtx.Context, cliCtx.App.Version)
		},
	}
}

func (c *cmdRulesLogs) run(ctx context.Context, version string) error {
	writer := func(topic, msg string) {
		fmt.Fprintf(c.writer, "[%s] %s\n", topic, msg)
	}

	streamer, err := cloudapi.NewRuleLogsWriter(c.websocketsURL, c.token,
		version, c.ruleID, writer)
	if err != nil {
		return fmt.Errorf("create streamer: %w", err)
	}

	return streamer.Run(ctx)
}
