package enaptercli

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v2"

	"github.com/enapter/enapter-cli/internal/cloudapi"
)

type cmdDevicesLogs struct {
	cmdDevices
}

func buildCmdDevicesLogs() *cli.Command {
	cmd := &cmdDevicesLogs{}

	return &cli.Command{
		Name:               "logs",
		Usage:              "Stream logs from a device",
		CustomHelpTemplate: cmd.HelpTemplate(),
		Flags:              cmd.Flags(),
		Before:             cmd.Before,
		Action: func(cliCtx *cli.Context) error {
			return cmd.run(cliCtx.Context, cliCtx.App.Version)
		},
	}
}

func (c *cmdDevicesLogs) run(ctx context.Context, version string) error {
	writer, err := cloudapi.NewDeviceLogsWriter(c.websocketsURL, c.token,
		version, c.hardwareID, c.writeLog)
	if err != nil {
		return fmt.Errorf("create writer: %w", err)
	}
	return writer.Run(ctx)
}

func (c *cmdDevicesLogs) writeLog(topic, msg string) {
	fmt.Fprintf(c.writer, "[%s] %s\n", topic, msg)
}
