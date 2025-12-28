package enaptercli

import (
	"bytes"
	"cmp"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/urfave/cli/v2"
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
		Action: func(cliCtx *cli.Context) error {
			return cmd.do(cliCtx.Context)
		},
	}
}

func (c *cmdDeviceCreateStandalone) Flags() []cli.Flag {
	flags := c.cmdDeviceCreate.Flags()
	return append(flags, &cli.StringFlag{
		Name:        "site-id",
		Aliases:     []string{"s"},
		Usage:       "site ID where the device will be created",
		Destination: &c.siteID,
	}, &cli.StringFlag{
		Name:        "device-name",
		Aliases:     []string{"n"},
		Usage:       "name for the new device",
		Destination: &c.deviceName,
		Required:    true,
	}, &cli.StringFlag{
		Name:        "device-slug",
		Usage:       "slug for the new standalone device",
		Destination: &c.deviceSlug,
	})
}

func (c *cmdDeviceCreateStandalone) do(ctx context.Context) error {
	if c.siteID != "" && c.cmdBase.siteID != "" && c.cmdBase.siteID != c.siteID {
		return errSiteIDMismatch
	}

	siteID := cmp.Or(c.siteID, c.cmdBase.siteID)
	if siteID == "" {
		return errSiteIDMissing
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
