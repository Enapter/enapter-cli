package enaptercli

import (
	"context"
	"net/http"

	"github.com/urfave/cli/v2"
)

type cmdBlueprintGet struct {
	cmdBlueprint
	blueprintID string
}

func buildCmdBlueprintGet() *cli.Command {
	cmd := &cmdBlueprintGet{}
	return &cli.Command{
		Name:               "get",
		Usage:              "Retrieve blueprint metadata",
		CustomHelpTemplate: cmd.CommandHelpTemplate(),
		Flags:              cmd.Flags(),
		Before:             cmd.Before,
		Action: func(cliCtx *cli.Context) error {
			return cmd.get(cliCtx.Context)
		},
	}
}

func (c *cmdBlueprintGet) Flags() []cli.Flag {
	flags := c.cmdBlueprint.Flags()
	return append(flags, &cli.StringFlag{
		Name:        "blueprint-id",
		Aliases:     []string{"b"},
		Usage:       "blueprint name or ID to retrieve",
		Destination: &c.blueprintID,
		Required:    true,
	})
}

func (c *cmdBlueprintGet) get(ctx context.Context) error {
	if isBlueprintID(c.blueprintID) {
		return c.doHTTPRequest(ctx, doHTTPRequestParams{
			Method: http.MethodGet,
			Path:   "/blueprints/" + c.blueprintID,
		})
	}

	blueprintName, blueprintTag := parseBlueprintName(c.blueprintID)
	return c.doHTTPRequest(ctx, doHTTPRequestParams{
		Method: http.MethodGet,
		Path:   "/blueprints/enapter/" + blueprintName + "/" + blueprintTag,
	})
}
