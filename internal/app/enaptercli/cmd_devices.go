package enaptercli

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/urfave/cli/v2"
)

type cmdDevices struct {
	cmdBase
	deviceID string
}

func buildCmdDevices() *cli.Command {
	return &cli.Command{
		Name:  "devices",
		Usage: "Manage devices",
		Subcommands: []*cli.Command{
			buildCmdDevicesAssignBlueprint(),
			buildCmdDevicesLogs(),
			buildCmdDevicesLogsf(),
		},
	}
}

func (c *cmdDevices) Flags() []cli.Flag {
	flags := c.cmdBase.Flags()
	return append(flags, &cli.StringFlag{
		Name:        "device_id",
		Aliases:     []string{"d"},
		Usage:       "device ID",
		Destination: &c.deviceID,
		Required:    true,
	})
}

func (c *cmdDevices) doHTTPRequest(ctx context.Context, p doHTTPRequestParams) error {
	p.Path = "/devices/" + c.deviceID
	return c.cmdBase.doHTTPRequest(ctx, p)
}

func (c *cmdDevices) parseAndDumpDeviceLogs(body io.Reader) (int, error) {
	var resp struct {
		Logs []struct {
			ReceivedAt string `json:"received_at"`
			Timestamp  string `json:"timestamp"`
			Severity   string `json:"severity"`
			Message    string `json:"message"`
		} `json:"logs"`
	}
	if err := json.NewDecoder(body).Decode(&resp); err != nil {
		return 0, fmt.Errorf("parse response body: %w", err)
	}
	for _, l := range resp.Logs {
		fmt.Fprintf(c.writer, "%s [%s] %s\n", l.ReceivedAt, l.Severity, l.Message)
	}
	return len(resp.Logs), nil
}
