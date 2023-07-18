package enaptercli

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/urfave/cli/v2"

	"github.com/enapter/enapter-cli/internal/publichttp"
)

type cmdDevicesTelemetry struct {
	cmdDevices
	items cli.StringSlice
}

func buildCmdDevicesTelemetry() *cli.Command {
	cmd := &cmdDevicesTelemetry{}

	return &cli.Command{
		Name:               "telemetry",
		Usage:              "Receive device telemetry",
		CustomHelpTemplate: cmd.HelpTemplate(),
		Flags:              cmd.Flags(),
		Before:             cmd.Before,
		Action: func(cliCtx *cli.Context) error {
			return cmd.run(cliCtx.Context)
		},
	}
}

func (c *cmdDevicesTelemetry) Flags() []cli.Flag {
	flags := c.cmdDevices.Flags()
	flags = append(flags,
		&cli.StringSliceFlag{
			Name:        "items",
			Usage:       "Telemetry items",
			Required:    true,
			Destination: &c.items,
		},
	)
	return flags
}

func (c *cmdDevicesTelemetry) run(ctx context.Context) error {
	transport := publichttp.NewAuthTokenTransport(http.DefaultTransport, c.token)
	client, err := publichttp.NewClientWithURL(&http.Client{Transport: transport}, c.apiHost)
	if err != nil {
		return fmt.Errorf("create http client: %w", err)
	}

	query := publichttp.NowQuery{
		Devices: map[string][]string{
			c.hardwareID: c.items.Value(),
		},
	}

	resp, err := client.Telemetry.Now(ctx, query)
	if err != nil {
		return err
	}

	if len(resp.Errors) > 0 {
		for _, err := range resp.Errors {
			fmt.Fprintf(c.writer, "[ERROR] %v (Details: %v)\n", err.Message, err.Details)
		}
		return errFinishedWithError
	}

	telemetryByName := resp.Devices[c.hardwareID]
	return c.printTelemetry(telemetryByName)
}

func (c *cmdDevicesTelemetry) printTelemetry(t publichttp.TelemetryByName) error {
	s, err := json.Marshal(t)
	if err != nil {
		return fmt.Errorf("format response: %w", err)
	}
	fmt.Fprintln(c.writer, string(s))
	return nil
}
