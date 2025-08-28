package enaptercli

import (
	"context"
	"net/http"

	"github.com/urfave/cli/v2"
)

type cmdRuleEngineRuleDelete struct {
	cmdRuleEngineRule
	ruleID string
}

func buildCmdRuleEngineRuleDelete() *cli.Command {
	cmd := &cmdRuleEngineRuleDelete{}
	return &cli.Command{
		Name:               "delete",
		Usage:              "Delete a rule",
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
		&cli.StringFlag{
			Name:        "rule-id",
			Usage:       "Rule ID or slug",
			Required:    true,
			Destination: &c.ruleID,
		},
	)
}

func (c *cmdRuleEngineRuleDelete) do(ctx context.Context) error {
	return c.doHTTPRequest(ctx, doHTTPRequestParams{
		Method:      http.MethodDelete,
		Path:        "/" + c.ruleID,
		ContentType: contentTypeJSON,
	})
}
