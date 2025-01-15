package enaptercli

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gorilla/websocket"
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
		CustomHelpTemplate: cmd.HelpTemplate(),
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
			Usage:       "follow log output",
			Destination: &c.follow,
		},
	)
}

func (c *cmdRuleEngineRuleLogs) do(cliCtx *cli.Context) error {
	if !c.follow {
		return cli.Exit("Currently, only follow mode (--follow) is supported.", 1)
	}

	path := fmt.Sprintf("/site/rule_engine/rules/%s/logs/ws", c.ruleID)
	for retry := false; ; retry = true {
		if retry {
			fmt.Fprintln(c.writer, "Reconnecting...")
			time.Sleep(time.Second)
		}

		conn, err := c.dialWebSocket(cliCtx.Context, path)
		if err != nil {
			select {
			case <-cliCtx.Done():
				return nil
			default:
				fmt.Fprintln(c.writer, "Failed to retrieve logs:", err)
				continue
			}
		}
		fmt.Fprintln(c.writer, "Connection established")

		closeCh := make(chan struct{})
		go func() {
			select {
			case <-cliCtx.Done():
			case <-closeCh:
			}
			conn.Close()
		}()

		if err := c.runProxy(conn); err != nil {
			select {
			case <-cliCtx.Done():
				return nil
			default:
				fmt.Fprintln(c.writer, "Failed to retrieve logs:", err)
				close(closeCh)
			}
		}
	}
}

func (c *cmdRuleEngineRuleLogs) runProxy(conn *websocket.Conn) error {
	for {
		_, r, err := conn.NextReader()
		if err != nil {
			return fmt.Errorf("read: %w", err)
		}

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
	}
}
