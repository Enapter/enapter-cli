package enaptercli

import (
	"context"
	"net/http"

	"github.com/urfave/cli/v3"
)

type cmdRuleEngineRuleList struct {
	cmdRuleEngineRule
	limit int
}

func buildCmdRuleEngineRuleList() *cli.Command {
	cmd := &cmdRuleEngineRuleList{}
	return &cli.Command{
		Name:               "list",
		Usage:              "List rules",
		CustomHelpTemplate: cmd.CommandHelpTemplate(),
		Flags:              cmd.Flags(),
		Before:             cmd.Before,
		Action: func(ctx context.Context, _ *cli.Command) error {
			return cmd.do(ctx)
		},
	}
}

func (c *cmdRuleEngineRuleList) Flags() []cli.Flag {
	flags := c.cmdRuleEngineRule.Flags()
	return append(flags, &cli.IntFlag{
		Name:        "limit",
		Usage:       "maximum number of rules to retrieve",
		Destination: &c.limit,
		DefaultText: "retrieves all",
	})
}

func (c *cmdRuleEngineRuleList) do(ctx context.Context) error {
	doPaginateRequestParams := paginateHTTPRequestParams{
		ObjectName: "rules",
		Limit:      c.limit,
		DoFn:       c.doHTTPRequest,
		BaseParams: doHTTPRequestParams{
			Method: http.MethodGet,
			Path:   "",
		},
	}

	return c.doPaginateRequest(ctx, doPaginateRequestParams)
}
