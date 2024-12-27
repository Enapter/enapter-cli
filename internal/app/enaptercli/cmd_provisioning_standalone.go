package enaptercli

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/urfave/cli/v2"
)

type cmdProvisioningStandalone struct {
	cmdProvisioning
	siteID     string
	deviceName string
}

func buildCmdProvisioningStandalone() *cli.Command {
	cmd := &cmdProvisioningStandalone{}
	return &cli.Command{
		Name:               "standalone",
		Usage:              "Create a new standalone device",
		CustomHelpTemplate: cmd.HelpTemplate(),
		Flags:              cmd.Flags(),
		Before:             cmd.Before,
		Action: func(cliCtx *cli.Context) error {
			return cmd.do(cliCtx.Context)
		},
	}
}

func (c *cmdProvisioningStandalone) Flags() []cli.Flag {
	flags := c.cmdProvisioning.Flags()
	return append(flags, &cli.StringFlag{
		Name:        "site-id",
		Aliases:     []string{"s"},
		Usage:       "site ID where to craate device",
		Destination: &c.siteID,
		Required:    true,
	}, &cli.StringFlag{
		Name:        "device-name",
		Aliases:     []string{"n"},
		Usage:       "name for a new device",
		Destination: &c.deviceName,
		Required:    true,
	})
}

func (c *cmdProvisioningStandalone) do(ctx context.Context) error {
	body, err := json.Marshal(map[string]interface{}{
		"site_id": c.siteID,
		"name":    c.deviceName,
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
