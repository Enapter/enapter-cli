package enaptercli

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/urfave/cli/v3"
)

type cmdRuleEngineRuleEnable struct {
	cmdRuleEngineRule
	ruleIDs []string
}

func buildCmdRuleEngineRuleEnable() *cli.Command {
	cmd := &cmdRuleEngineRuleEnable{}
	return &cli.Command{
		Name:               "enable",
		Usage:              "Enable one or more rules",
		CustomHelpTemplate: cmd.CommandHelpTemplate(),
		Flags:              cmd.Flags(),
		Before:             cmd.Before,
		Action: func(ctx context.Context, _ *cli.Command) error {
			return cmd.do(ctx)
		},
	}
}

func (c *cmdRuleEngineRuleEnable) Flags() []cli.Flag {
	return append(c.cmdRuleEngineRule.Flags(),
		&cli.StringSliceFlag{
			Name:        "rule-id",
			Usage:       "Rule IDs or slugs",
			Required:    true,
			Destination: &c.ruleIDs,
		},
	)
}

func (c *cmdRuleEngineRuleEnable) do(ctx context.Context) error {
	body, err := json.Marshal(map[string]any{
		"rule_ids": c.ruleIDs,
	})
	if err != nil {
		return fmt.Errorf("build request: %w", err)
	}
	return c.doHTTPRequest(ctx, doHTTPRequestParams{
		Method:      http.MethodPost,
		Path:        "/batch_enable",
		Body:        bytes.NewReader(body),
		ContentType: contentTypeJSON,
	})
}
