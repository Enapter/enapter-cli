package enaptercli

import (
	"context"
	"net/http"

	"github.com/urfave/cli/v3"
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
		ArgsUsage:          "<blueprint id or name>",
		Action: func(ctx context.Context, cm *cli.Command) error {
			return cmd.inspect(ctx, cm.Args().Get(0))
		},
	}
}

func (c *cmdBlueprintsInpsect) Before(
	ctx context.Context, cm *cli.Command,
) (context.Context, error) {
	ctx, err := c.cmdBlueprints.Before(ctx, cm)
	if err != nil {
		return nil, err
	}
	if cm.Args().Get(0) == "" {
		return nil, errBlueprintIDMissed
	}
	return ctx, nil
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
