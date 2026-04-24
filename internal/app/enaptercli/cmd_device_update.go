package enaptercli

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/urfave/cli/v3"

	"github.com/enapter/enapter-cli/v3/internal/app/cliflags"
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
		Action: func(ctx context.Context, _ *cli.Command) error {
			return cmd.do(ctx)
		},
	}
}

func (c *cmdDeviceUpdate) Flags() []cli.Flag {
	flags := c.cmdDevice.Flags()
	return append(flags,
		&cli.StringFlag{
			Name:        "device-id",
			Aliases:     []string{"d"},
			Usage:       "Device ID",
			Destination: &c.deviceID,
			Required:    true,
		},
		&cli.StringFlag{
			Name:        "name",
			Usage:       "Device name",
			Destination: &c.name,
			Action:      cliflags.TrimSpaceAction(&c.name),
		},
		&cli.StringFlag{
			Name:        "slug",
			Usage:       "Device slug",
			Destination: &c.slug,
			Action:      cliflags.TrimSpaceAction(&c.slug),
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
