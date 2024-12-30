package enaptercli

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/urfave/cli/v2"
)

type cmdProvisioningLua struct {
	cmdProvisioning
	deviceName  string
	runtimeID   string
	blueprintID string
}

func buildCmdProvisioningLua() *cli.Command {
	cmd := &cmdProvisioningLua{}
	return &cli.Command{
		Name:               "lua-device",
		Usage:              "Create a new Lua device",
		CustomHelpTemplate: cmd.HelpTemplate(),
		Flags:              cmd.Flags(),
		Before:             cmd.Before,
		Action: func(cliCtx *cli.Context) error {
			return cmd.do(cliCtx.Context)
		},
	}
}

func (c *cmdProvisioningLua) Flags() []cli.Flag {
	flags := c.cmdProvisioning.Flags()
	return append(flags, &cli.StringFlag{
		Name:        "runtime-id",
		Aliases:     []string{"r"},
		Usage:       "runtime UCM device ID where to run a new Lua device",
		Destination: &c.runtimeID,
		Required:    true,
	}, &cli.StringFlag{
		Name:        "device-name",
		Aliases:     []string{"n"},
		Usage:       "name of a new Lua device",
		Destination: &c.deviceName,
		Required:    true,
	}, &cli.StringFlag{
		Name:        "blueprint-id",
		Aliases:     []string{"b"},
		Usage:       "blueprint ID of a new Lua device",
		Destination: &c.blueprintID,
		Required:    true,
	})
}

func (c *cmdProvisioningLua) do(ctx context.Context) error {
	body, err := json.Marshal(map[string]interface{}{
		"runtime_id":   c.runtimeID,
		"name":         c.deviceName,
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
