package enaptercli

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/urfave/cli/v2"
)

type cmdDevicesAssignBlueprint struct {
	cmdDevices
	blueprintID string
}

func buildCmdDevicesAssignBlueprint() *cli.Command {
	cmd := &cmdDevicesAssignBlueprint{}
	return &cli.Command{
		Name:               "assign_blueprint",
		Usage:              "Assign blueprint to device",
		CustomHelpTemplate: cmd.HelpTemplate(),
		Flags:              cmd.Flags(),
		Before:             cmd.Before,
		Action: func(cliCtx *cli.Context) error {
			return cmd.do(cliCtx.Context)
		},
	}
}

func (c *cmdDevicesAssignBlueprint) Flags() []cli.Flag {
	flags := c.cmdDevices.Flags()
	return append(flags, &cli.StringFlag{
		Name:        "blueprint-id",
		Aliases:     []string{"b"},
		Usage:       "blueprint ID to assign",
		Destination: &c.blueprintID,
		Required:    true,
	})
}

func (c *cmdDevicesAssignBlueprint) do(ctx context.Context) error {
	body, err := json.Marshal(map[string]interface{}{
		"blueprint_id": c.blueprintID,
	})
	if err != nil {
		return fmt.Errorf("build request: %w", err)
	}
	return c.doHTTPRequest(ctx, doHTTPRequestParams{
		Method:      http.MethodPost,
		Path:        "/assign_blueprint",
		Body:        bytes.NewReader(body),
		ContentType: contentTypeJSON,
	})
}
