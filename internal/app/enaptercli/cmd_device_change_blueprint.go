package enaptercli

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/urfave/cli/v2"
)

type cmdDeviceChangeBlueprint struct {
	cmdDevice
	deviceID      string
	blueprintID   string
	blueprintPath string
}

func buildCmdDeviceChangeBlueprint() *cli.Command {
	cmd := &cmdDeviceChangeBlueprint{}
	return &cli.Command{
		Name:               "change-blueprint",
		Usage:              "Change device blueprint",
		CustomHelpTemplate: cmd.CommandHelpTemplate(),
		Flags:              cmd.Flags(),
		Before:             cmd.Before,
		Action: func(cliCtx *cli.Context) error {
			return cmd.do(cliCtx.Context)
		},
	}
}

func (c *cmdDeviceChangeBlueprint) Flags() []cli.Flag {
	flags := c.cmdDevice.Flags()
	return append(flags, &cli.StringFlag{
		Name:        "device-id",
		Aliases:     []string{"d"},
		Usage:       "Device ID",
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

func (c *cmdDeviceChangeBlueprint) Before(cliCtx *cli.Context) error {
	if err := c.cmdDevice.Before(cliCtx); err != nil {
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

func (c *cmdDeviceChangeBlueprint) do(ctx context.Context) error {
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
