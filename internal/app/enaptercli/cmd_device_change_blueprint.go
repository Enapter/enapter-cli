package enaptercli

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/urfave/cli/v2"
)

type cmdDevicesChangeBlueprint struct {
	cmdDevices
	deviceID      string
	blueprintID   string
	blueprintPath string
}

func buildCmdDevicesChangeBlueprint() *cli.Command {
	cmd := &cmdDevicesChangeBlueprint{}
	return &cli.Command{
		Name:               "change-blueprint",
		Usage:              "Change blueprint to device",
		CustomHelpTemplate: cmd.CommandHelpTemplate(),
		Flags:              cmd.Flags(),
		Before:             cmd.Before,
		Action: func(cliCtx *cli.Context) error {
			return cmd.do(cliCtx.Context)
		},
	}
}

func (c *cmdDevicesChangeBlueprint) Flags() []cli.Flag {
	flags := c.cmdDevices.Flags()
	return append(flags, &cli.StringFlag{
		Name:        "device-id",
		Aliases:     []string{"d"},
		Usage:       "device ID",
		Destination: &c.deviceID,
		Required:    true,
	}, &cli.StringFlag{
		Name:        "blueprint-id",
		Aliases:     []string{"b"},
		Usage:       "blueprint ID to use as new device blueprint",
		Destination: &c.blueprintID,
	}, &cli.StringFlag{
		Name:        "blueprint-path",
		Usage:       "blueprint path (zip file or directory) to use as new device blueprint",
		Destination: &c.blueprintPath,
	})
}

func (c *cmdDevicesChangeBlueprint) Before(cliCtx *cli.Context) error {
	if err := c.cmdDevices.Before(cliCtx); err != nil {
		return err
	}
	if c.blueprintID != "" && c.blueprintPath != "" {
		return errOnlyOneBlueprinFlag
	}
	if c.blueprintID == "" && c.blueprintPath == "" {
		return errMissedBlueprintFlag
	}
	return c.validateExpandFlag(cliCtx)
}

func (c *cmdDevicesChangeBlueprint) do(ctx context.Context) error {
	if c.blueprintPath != "" {
		blueprintID, err := uploadBlueprintAndReturnBlueprintID(ctx, c.blueprintPath, c.cmdBase.doHTTPRequest)
		if err != nil {
			return fmt.Errorf("upload blueprint: %w", err)
		}
		c.blueprintID = blueprintID
	}

	body, err := json.Marshal(map[string]any{
		"blueprint_id": c.blueprintID,
	})
	if err != nil {
		return fmt.Errorf("build request: %w", err)
	}
	return c.doHTTPRequest(ctx, doHTTPRequestParams{
		Method:      http.MethodPost,
		Path:        "/" + c.deviceID + "/assign_blueprint",
		Body:        bytes.NewReader(body),
		ContentType: contentTypeJSON,
	})
}
