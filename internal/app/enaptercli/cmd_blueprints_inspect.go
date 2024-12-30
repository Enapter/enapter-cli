package enaptercli

import (
	"context"
	"net/http"

	"github.com/urfave/cli/v2"
)

type cmdBlueprintsInpsect struct {
	cmdBlueprints
	blueprintID string
}

func buildCmdBlueprintsInspect() *cli.Command {
	cmd := &cmdBlueprintsInpsect{}
	return &cli.Command{
		Name:               "inspect",
		Usage:              "Get blueprint metainfo",
		CustomHelpTemplate: cmd.HelpTemplate(),
		Flags:              cmd.Flags(),
		Before:             cmd.Before,
		Action: func(cliCtx *cli.Context) error {
			return cmd.inspect(cliCtx.Context)
		},
	}
}

func (c *cmdBlueprintsInpsect) Flags() []cli.Flag {
	flags := c.cmdBlueprints.Flags()
	return append(flags, &cli.StringFlag{
		Name:        "blueprint-id",
		Aliases:     []string{"b"},
		Usage:       "blueprint name or ID to inspect",
		Destination: &c.blueprintID,
		Required:    true,
	})
}

func (c *cmdBlueprintsInpsect) inspect(ctx context.Context) error {
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
