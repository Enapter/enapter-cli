package enaptercli

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/urfave/cli/v2"
)

type cmdRuleEngineRuleDelete struct {
	cmdRuleEngineRule
	ruleIDs []string
}

func buildCmdRuleEngineRuleDelete() *cli.Command {
	cmd := &cmdRuleEngineRuleDelete{}
	return &cli.Command{
		Name:               "delete",
		Usage:              "Delete one or more rules",
		CustomHelpTemplate: cmd.CommandHelpTemplate(),
		Flags:              cmd.Flags(),
		Before:             cmd.Before,
		Action: func(cliCtx *cli.Context) error {
			return cmd.do(cliCtx.Context)
		},
	}
}

func (c *cmdRuleEngineRuleDelete) Flags() []cli.Flag {
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

func (c *cmdRuleEngineRuleDelete) do(ctx context.Context) error {
	body, err := json.Marshal(map[string]any{
		"rule_ids": c.ruleIDs,
	})
	if err != nil {
		return fmt.Errorf("build request: %w", err)
	}
	return c.doHTTPRequest(ctx, doHTTPRequestParams{
		Method:      http.MethodPost,
		Path:        "/batch_delete",
		Body:        bytes.NewReader(body),
		ContentType: contentTypeJSON,
	})
}
