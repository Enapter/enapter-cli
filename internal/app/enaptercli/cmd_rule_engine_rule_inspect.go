package enaptercli

import (
	"context"
	"net/http"

	"github.com/urfave/cli/v2"
)

type cmdRuleEngineRuleInspect struct {
	cmdRuleEngineRule
	ruleID string
}

func buildCmdRuleEngineRuleInspect() *cli.Command {
	cmd := &cmdRuleEngineRuleInspect{}
	return &cli.Command{
		Name:               "inspect",
		Usage:              "Inspect a rule",
		CustomHelpTemplate: cmd.HelpTemplate(),
		Flags:              cmd.Flags(),
		Before:             cmd.Before,
		Action: func(cliCtx *cli.Context) error {
			return cmd.do(cliCtx.Context)
		},
	}
}

func (c *cmdRuleEngineRuleInspect) Flags() []cli.Flag {
	return append(c.cmdRuleEngineRule.Flags(),
		&cli.StringFlag{
			Name:        "rule-id",
			Usage:       "Rule ID or slug",
			Destination: &c.ruleID,
			Required:    true,
		},
	)
}

func (c *cmdRuleEngineRuleInspect) do(ctx context.Context) error {
	return c.doHTTPRequest(ctx, doHTTPRequestParams{
		Method: http.MethodGet,
		Path:   "/" + c.ruleID,
	})
}
