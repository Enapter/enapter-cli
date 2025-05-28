package enaptercli

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/urfave/cli/v2"
)

type cmdDevicesExecuteCommand struct {
	cmdDevices
	deviceID  string
	cmdName   string
	cmdArgs   string
	ephemeral bool
}

func buildCmdDevicesExecuteCommand() *cli.Command {
	cmd := &cmdDevicesExecuteCommand{}
	return &cli.Command{
		Name:               "execute-command",
		Usage:              "Execute a device command",
		CustomHelpTemplate: cmd.HelpTemplate(),
		Flags:              cmd.Flags(),
		Before:             cmd.Before,
		Action: func(cliCtx *cli.Context) error {
			return cmd.do(cliCtx.Context)
		},
	}
}

func (c *cmdDevicesExecuteCommand) Flags() []cli.Flag {
	flags := c.cmdDevices.Flags()
	return append(flags,
		&cli.StringFlag{
			Name:        "device-id",
			Usage:       "device ID",
			Destination: &c.deviceID,
			Required:    true,
		},
		&cli.StringFlag{
			Name:        "cmd-name",
			Usage:       "command name",
			Destination: &c.cmdName,
			Required:    true,
		},
		&cli.StringFlag{
			Name:        "cmd-args",
			Usage:       "command args (should be a JSON string)",
			Destination: &c.cmdArgs,
		},
		&cli.BoolFlag{
			Name:        "ephemeral",
			Usage:       "run command in ephemeral mode",
			Destination: &c.ephemeral,
			Hidden:      true,
		},
	)
}

func (c *cmdDevicesExecuteCommand) do(ctx context.Context) error {
	reqBody := struct {
		Name      string          `json:"name"`
		Args      json.RawMessage `json:"arguments,omitempty"`
		Ephemeral bool            `json:"ephemeral,omitempty"`
	}{
		Name:      c.cmdName,
		Args:      json.RawMessage(c.cmdArgs),
		Ephemeral: c.ephemeral,
	}
	data, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("build request: %w", err)
	}

	return c.doHTTPRequest(ctx, doHTTPRequestParams{
		Method: http.MethodPost,
		Path:   "/" + c.deviceID + "/execute_command",
		Body:   bytes.NewReader(data),
	})
}
