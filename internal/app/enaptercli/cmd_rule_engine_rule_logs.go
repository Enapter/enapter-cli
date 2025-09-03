package enaptercli

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/urfave/cli/v2"
)

type cmdRuleEngineRuleLogs struct {
	cmdRuleEngineRule
	ruleID string
	follow bool
}

func buildCmdRuleEngineRuleLogs() *cli.Command {
	cmd := &cmdRuleEngineRuleLogs{}
	return &cli.Command{
		Name:               "logs",
		Usage:              "Show rule logs",
		CustomHelpTemplate: cmd.CommandHelpTemplate(),
		Flags:              cmd.Flags(),
		Before:             cmd.Before,
		Action: func(cliCtx *cli.Context) error {
			return cmd.do(cliCtx)
		},
	}
}

func (c *cmdRuleEngineRuleLogs) Flags() []cli.Flag {
	return append(c.cmdRuleEngineRule.Flags(),
		&cli.StringFlag{
			Name:        "rule-id",
			Usage:       "rule ID",
			Destination: &c.ruleID,
			Required:    true,
		},
		&cli.BoolFlag{
			Name:        "follow",
			Aliases:     []string{"f"},
			Usage:       "follow the log output",
			Destination: &c.follow,
		},
	)
}

func (c *cmdRuleEngineRuleLogs) do(cliCtx *cli.Context) error {
	if !c.follow {
		return cli.Exit("Currently, only follow mode (--follow) is supported.", 1)
	}

	path := fmt.Sprintf("/site/rule_engine/rules/%s/logs", c.ruleID)

	return c.runWebSocket(cliCtx.Context, runWebSocketParams{
		Path: path,
		RespProcessor: func(r io.Reader) error {
			var msg struct {
				Timestamp int64  `json:"timestamp"`
				Severity  string `json:"severity"`
				Message   string `json:"message"`
			}
			if err := json.NewDecoder(r).Decode(&msg); err != nil {
				return fmt.Errorf("parse payload: %w", err)
			}
			ts := time.Unix(msg.Timestamp, 0).Format(time.RFC3339)
			fmt.Fprintf(c.writer, "%s [%s] %s\n", ts, msg.Severity, msg.Message)
			return nil
		},
	})
}
