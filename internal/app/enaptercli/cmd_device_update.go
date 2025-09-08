package enaptercli

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/urfave/cli/v2"
)

type cmdDeviceUpdate struct {
	cmdDevice
	deviceID string
	name     string
	slug     string
}

func buildCmdDeviceUpdate() *cli.Command {
	cmd := &cmdDeviceUpdate{}
	return &cli.Command{
		Name:               "update",
		Usage:              "Update a device",
		CustomHelpTemplate: cmd.CommandHelpTemplate(),
		Flags:              cmd.Flags(),
		Before:             cmd.Before,
		Action: func(cliCtx *cli.Context) error {
			return cmd.do(cliCtx.Context)
		},
	}
}

func (c *cmdDeviceUpdate) Flags() []cli.Flag {
	flags := c.cmdDevice.Flags()
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
			Usage:       "device name",
			Destination: &c.name,
		},
		&cli.StringFlag{
			Name:        "slug",
			Usage:       "device slug",
			Destination: &c.slug,
		},
	)
}

func (c *cmdDeviceUpdate) do(ctx context.Context) error {
	payload := map[string]string{
		"name": c.name,
		"slug": c.slug,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("build request: %w", err)
	}

	return c.doHTTPRequest(ctx, doHTTPRequestParams{
		Method:      http.MethodPatch,
		Path:        "/" + c.deviceID,
		Body:        bytes.NewReader(body),
		ContentType: contentTypeJSON,
	})
}
