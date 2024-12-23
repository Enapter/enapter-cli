package enaptercli

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/urfave/cli/v3"
)

type cmdRuleEngineRuleDisable struct {
	cmdRuleEngineRule
}

func buildCmdRuleEngineRuleDisable() *cli.Command {
	cmd := &cmdRuleEngineRuleDisable{}
	return &cli.Command{
		Name:               "disable",
		Usage:              "Disable one or more rules",
		ArgsUsage:          "RULE [RULE...]",
		CustomHelpTemplate: cmd.HelpTemplate(),
		Flags:              cmd.Flags(),
		Before:             cmd.Before,
		Action: func(ctx context.Context, cm *cli.Command) error {
			return cmd.do(ctx, cm.Args().Slice())
		},
	}
}

func (c *cmdRuleEngineRuleDisable) Before(
	ctx context.Context, cm *cli.Command,
) (context.Context, error) {
	ctx, err := c.cmdRuleEngineRule.Before(ctx, cm)
	if err != nil {
		return nil, err
	}

	if cm.Args().Get(0) == "" {
		return nil, errRequiresAtLeastOneArgument
	}

	return ctx, nil
}

func (c *cmdRuleEngineRuleDisable) do(ctx context.Context, rules []string) error {
	body, err := json.Marshal(map[string]any{
		"rule_ids": rules,
	})
	if err != nil {
		return fmt.Errorf("build request: %w", err)
	}
	return c.doHTTPRequest(ctx, doHTTPRequestParams{
		Method: http.MethodPost,
		Path:   "/batch_disable",
		Body:   bytes.NewReader(body),
	})
}
