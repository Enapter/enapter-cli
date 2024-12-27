package enaptercli

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/urfave/cli/v2"
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
		CustomHelpTemplate: cmd.HelpTemplate(),
		Flags:              cmd.Flags(),
		Before:             cmd.Before,
		Action: func(cliCtx *cli.Context) error {
			return cmd.do(cliCtx.Context)
		},
	}
}

func (c *cmdRuleEngineRuleEnable) Flags() []cli.Flag {
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
