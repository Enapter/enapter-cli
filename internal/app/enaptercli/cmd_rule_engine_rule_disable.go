package enaptercli

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/urfave/cli/v2"
)

type cmdRuleEngineRuleDisable struct {
	cmdRuleEngineRule
	ruleIDs []string
}

func buildCmdRuleEngineRuleDisable() *cli.Command {
	cmd := &cmdRuleEngineRuleDisable{}
	return &cli.Command{
		Name:               "disable",
		Usage:              "Disable one or more rules",
		CustomHelpTemplate: cmd.HelpTemplate(),
		Flags:              cmd.Flags(),
		Before:             cmd.Before,
		Action: func(cliCtx *cli.Context) error {
			return cmd.do(cliCtx.Context)
		},
	}
}

func (c *cmdRuleEngineRuleDisable) Flags() []cli.Flag {
	return append(c.cmdRuleEngineRule.Flags(),
		&cli.MultiStringFlag{
			Target: &cli.StringSliceFlag{
				Name:     "rule-id",
				Usage:    "Rule IDs or slugs",
				Required: true,
			},
			Destination: &c.ruleIDs,
		},
	)
}

func (c *cmdRuleEngineRuleDisable) do(ctx context.Context) error {
	body, err := json.Marshal(map[string]any{
		"rule_ids": c.ruleIDs,
	})
	if err != nil {
		return fmt.Errorf("build request: %w", err)
	}
	return c.doHTTPRequest(ctx, doHTTPRequestParams{
		Method:      http.MethodPost,
		Path:        "/batch_disable",
		Body:        bytes.NewReader(body),
		ContentType: contentTypeJSON,
	})
}
