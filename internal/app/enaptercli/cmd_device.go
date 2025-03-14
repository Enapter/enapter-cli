package enaptercli

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/url"

	"github.com/urfave/cli/v2"
)

type cmdDevices struct {
	cmdBase
}

func buildCmdDevices() *cli.Command {
	return &cli.Command{
		Name:  "device",
		Usage: "Manage devices",
		Subcommands: []*cli.Command{
			buildCmdDevicesList(),
			buildCmdDevicesInspect(),
			buildCmdDevicesAssignBlueprint(),
			buildCmdDevicesLogs(),
			buildCmdDevicesLogsf(),
			buildCmdDevicesDelete(),
			buildCmdDevicesExecuteCommand(),
			buildCmdDeviceExecution(),
		},
	}
}

func (c *cmdDevices) doHTTPRequest(ctx context.Context, p doHTTPRequestParams) error {
	path, err := url.JoinPath("/devices", p.Path)
	if err != nil {
		return fmt.Errorf("join path: %w", err)
	}
	p.Path = path
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

func (c *cmdDevices) validateExpandFlag(cliCtx *cli.Context) error {
	supportedFields := []string{"connectivity", "manifest", "properties", "communication_info", "site"}
	return validateExpandFlag(cliCtx, supportedFields)
}
