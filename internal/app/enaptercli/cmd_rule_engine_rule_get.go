package enaptercli

import (
	"context"
	"net/http"

	"github.com/urfave/cli/v3"
)

type cmdRuleEngineRuleGet struct {
	cmdRuleEngineRule
	ruleID string
}

func buildCmdRuleEngineRuleGet() *cli.Command {
	cmd := &cmdRuleEngineRuleGet{}
	return &cli.Command{
		Name:               "get",
		Usage:              "Retrieve a rule",
		CustomHelpTemplate: cmd.CommandHelpTemplate(),
		Flags:              cmd.Flags(),
		Before:             cmd.Before,
		Action: func(ctx context.Context, _ *cli.Command) error {
			return cmd.do(ctx)
		},
	}
}

func (c *cmdRuleEngineRuleGet) Flags() []cli.Flag {
	return append(c.cmdRuleEngineRule.Flags(),
		&cli.StringFlag{
			Name:        "rule-id",
			Usage:       "Rule ID or slug",
			Destination: &c.ruleID,
			Required:    true,
		},
	)
}

func (c *cmdRuleEngineRuleGet) do(ctx context.Context) error {
	return c.doHTTPRequest(ctx, doHTTPRequestParams{
		Method: http.MethodGet,
		Path:   "/" + c.ruleID,
	})
}
