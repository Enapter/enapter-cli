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
}

func buildCmdRuleEngineRuleDisable() *cli.Command {
	cmd := &cmdRuleEngineRuleDisable{}
	return &cli.Command{
		Name:               "disable",
		Usage:              "Disable one or more rules",
		Args:               true,
		ArgsUsage:          "RULE [RULE...]",
		CustomHelpTemplate: cmd.HelpTemplate(),
		Flags:              cmd.Flags(),
		Before:             cmd.Before,
		Action: func(cliCtx *cli.Context) error {
			return cmd.do(cliCtx.Context, cliCtx.Args().Slice())
		},
	}
}

func (c *cmdRuleEngineRuleDisable) Before(cliCtx *cli.Context) error {
	if err := c.cmdRuleEngineRule.Before(cliCtx); err != nil {
		return err
	}

	if cliCtx.Args().Get(0) == "" {
		return errRequiresAtLeastOneArgument
	}

	return nil
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
