package enaptercli

import (
	"context"
	"net/http"

	"github.com/urfave/cli/v2"
)

type cmdBlueprintsGet struct {
	cmdBlueprints
	blueprintID string
}

func buildCmdBlueprintsGet() *cli.Command {
	cmd := &cmdBlueprintsGet{}
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

func (c *cmdBlueprintsGet) Flags() []cli.Flag {
	flags := c.cmdBlueprints.Flags()
	return append(flags, &cli.StringFlag{
		Name:        "blueprint-id",
		Aliases:     []string{"b"},
		Usage:       "blueprint name or ID to retrieve",
		Destination: &c.blueprintID,
		Required:    true,
	})
}

func (c *cmdBlueprintsGet) get(ctx context.Context) error {
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
