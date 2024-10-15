package enaptercli

import (
	"context"
	"net/http"

	"github.com/urfave/cli/v2"
)

type cmdBlueprintsInpsect struct {
	cmdBlueprints
}

func buildCmdBlueprintsInspect() *cli.Command {
	cmd := &cmdBlueprintsInpsect{}
	return &cli.Command{
		Name:               "inspect",
		Usage:              "Get blueprint metainfo",
		CustomHelpTemplate: cmd.HelpTemplate(),
		Flags:              cmd.Flags(),
		Before:             cmd.Before,
		Args:               true,
		ArgsUsage:          "<blueprint id or name>",
		Action: func(cliCtx *cli.Context) error {
			return cmd.inspect(cliCtx.Context, cliCtx.Args().Get(0))
		},
	}
}

func (c *cmdBlueprintsInpsect) Before(cliCtx *cli.Context) error {
	if err := c.cmdBlueprints.Before(cliCtx); err != nil {
		return err
	}
	if cliCtx.Args().Get(0) == "" {
		return errBlueprintIDMissed
	}
	return nil
}

func (c *cmdBlueprintsInpsect) inspect(ctx context.Context, blueprintID string) error {
	if isBlueprintID(blueprintID) {
		return c.doHTTPRequest(ctx, doHTTPRequestParams{
			Method: http.MethodGet,
			Path:   "/blueprints/" + blueprintID,
		})
	}

	blueprintName, blueprintTag := parseBlueprintName(blueprintID)
	return c.doHTTPRequest(ctx, doHTTPRequestParams{
		Method: http.MethodGet,
		Path:   "/blueprints/enapter/" + blueprintName + "/" + blueprintTag,
	})
}
