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

type cmdDeviceCreateStandalone struct {
	cmdDeviceCreate
	siteID     string
	deviceName string
	deviceSlug string
}

func buildCmdDeviceCreateStandalone() *cli.Command {
	cmd := &cmdDeviceCreateStandalone{}
	return &cli.Command{
		Name:               "standalone",
		Usage:              "Create a new standalone device",
		CustomHelpTemplate: cmd.CommandHelpTemplate(),
		Flags:              cmd.Flags(),
		Before:             cmd.Before,
		Action: func(ctx context.Context, _ *cli.Command) error {
			return cmd.do(ctx)
		},
	}
}

func (c *cmdDeviceCreateStandalone) Flags() []cli.Flag {
	flags := c.cmdDeviceCreate.Flags()
	return append(flags, &cli.StringFlag{
		Name:        "site-id",
		Aliases:     []string{"s"},
		Usage:       "Site ID where the device will be created",
		Destination: &c.siteID,
	}, &cli.StringFlag{
		Name:        "device-name",
		Aliases:     []string{"n"},
		Usage:       "Name for the new device",
		Destination: &c.deviceName,
		Required:    true,
		Action:      cliflags.TrimSpaceAction(&c.deviceName),
	}, &cli.StringFlag{
		Name:        "device-slug",
		Usage:       "Slug for the new standalone device",
		Destination: &c.deviceSlug,
		Action:      cliflags.TrimSpaceAction(&c.deviceSlug),
	})
}

func (c *cmdDeviceCreateStandalone) do(ctx context.Context) error {
	siteID, err := c.chooseSiteID(c.siteID)
	if err != nil {
		return err
	}

	body, err := json.Marshal(map[string]any{
		"site_id": siteID,
		"name":    c.deviceName,
		"slug":    c.deviceSlug,
	})
	if err != nil {
		return fmt.Errorf("build request: %w", err)
	}
	return c.doHTTPRequest(ctx, doHTTPRequestParams{
		Method:      http.MethodPost,
		Path:        "/provisioning/standalone",
		Body:        bytes.NewReader(body),
		ContentType: contentTypeJSON,
	})
}
