package enaptercli

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/urfave/cli/v2"
)

type cmdDeviceCreateLua struct {
	cmdDeviceCreate
	deviceName    string
	deviceSlug    string
	runtimeID     string
	blueprintID   string
	blueprintPath string
}

func buildCmdDeviceCreateLua() *cli.Command {
	cmd := &cmdDeviceCreateLua{}
	return &cli.Command{
		Name:               "lua-device",
		Usage:              "Create a new Lua device",
		CustomHelpTemplate: cmd.CommandHelpTemplate(),
		Flags:              cmd.Flags(),
		Before:             cmd.Before,
		Action: func(cliCtx *cli.Context) error {
			return cmd.do(cliCtx.Context)
		},
	}
}

func (c *cmdDeviceCreateLua) Flags() []cli.Flag {
	flags := c.cmdDeviceCreate.Flags()
	return append(flags, &cli.StringFlag{
		Name:        "runtime-id",
		Aliases:     []string{"r"},
		Usage:       "UCM device ID where the new Lua device will run",
		Destination: &c.runtimeID,
		Required:    true,
	}, &cli.StringFlag{
		Name:        "device-name",
		Aliases:     []string{"n"},
		Usage:       "name for the new Lua device",
		Destination: &c.deviceName,
		Required:    true,
	}, &cli.StringFlag{
		Name:        "device-slug",
		Usage:       "slug for the new Lua device",
		Destination: &c.deviceSlug,
	}, &cli.StringFlag{
		Name:        "blueprint-id",
		Aliases:     []string{"b"},
		Usage:       "blueprint ID to use for the new Lua device",
		Destination: &c.blueprintID,
	}, &cli.StringFlag{
		Name:        "blueprint-path",
		Usage:       "blueprint path (zip file or directory) to use for the new Lua device",
		Destination: &c.blueprintPath,
	})
}

func (c *cmdDeviceCreateLua) Before(cliCtx *cli.Context) error {
	if err := c.cmdDeviceCreate.Before(cliCtx); err != nil {
		return err
	}
	if c.blueprintID != "" && c.blueprintPath != "" {
		return errOnlyOneBlueprinFlag
	}
	if c.blueprintID == "" && c.blueprintPath == "" {
		return errMissedBlueprintFlag
	}
	return nil
}

func (c *cmdDeviceCreateLua) do(ctx context.Context) error {
	if c.blueprintPath != "" {
		blueprintID, err := uploadBlueprintAndReturnBlueprintID(ctx, c.blueprintPath, c.cmdBase.doHTTPRequest)
		if err != nil {
			return fmt.Errorf("upload blueprint: %w", err)
		}
		c.blueprintID = blueprintID
	}

	body, err := json.Marshal(map[string]interface{}{
		"runtime_id":   c.runtimeID,
		"name":         c.deviceName,
		"slug":         c.deviceSlug,
		"blueprint_id": c.blueprintID,
	})
	if err != nil {
		return fmt.Errorf("build request: %w", err)
	}
	return c.doHTTPRequest(ctx, doHTTPRequestParams{
		Method:      http.MethodPost,
		Path:        "/provisioning/lua_device",
		Body:        bytes.NewReader(body),
		ContentType: contentTypeJSON,
	})
}
