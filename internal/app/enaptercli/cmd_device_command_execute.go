package enaptercli

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/urfave/cli/v2"
)

type cmdDeviceCommandExecute struct {
	cmdDevices
	deviceID  string
	cmdName   string
	cmdArgs   string
	ephemeral bool
}

func buildCmdDeviceCommandExecute() *cli.Command {
	cmd := &cmdDeviceCommandExecute{}
	return &cli.Command{
		Name:               "execute",
		Usage:              "Execute a device command",
		CustomHelpTemplate: cmd.CommandHelpTemplate(),
		Flags:              cmd.Flags(),
		Before:             cmd.Before,
		Action: func(cliCtx *cli.Context) error {
			return cmd.do(cliCtx.Context)
		},
	}
}

func (c *cmdDeviceCommandExecute) Flags() []cli.Flag {
	flags := c.cmdDevices.Flags()
	return append(flags,
		&cli.StringFlag{
			Name:        "device-id",
			Aliases:     []string{"d"},
			Usage:       "device ID",
			Destination: &c.deviceID,
			Required:    true,
		},
		&cli.StringFlag{
			Name:        "name",
			Usage:       "command name",
			Destination: &c.cmdName,
			Required:    true,
		},
		&cli.StringFlag{
			Name:        "arguments",
			Usage:       "command arguments (should be a JSON string)",
			Destination: &c.cmdArgs,
		},
		&cli.BoolFlag{
			Name:        "ephemeral",
			Usage:       "run command in ephemeral mode",
			Destination: &c.ephemeral,
		},
	)
}

func (c *cmdDeviceCommandExecute) do(ctx context.Context) error {
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
