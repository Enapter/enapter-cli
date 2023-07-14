package enaptercli

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/urfave/cli/v2"

	"github.com/enapter/enapter-cli/internal/cloudapi"
)

type cmdDevicesUploadLogs struct {
	cmdDevices
	operationID string
	timeout     time.Duration
}

func buildCmdDevicesUploadLogs() *cli.Command {
	cmd := &cmdDevicesUploadLogs{}

	return &cli.Command{
		Name:               "upload-logs",
		Usage:              "Show blueprint uploading logs",
		CustomHelpTemplate: cmd.HelpTemplate(),
		Flags:              cmd.Flags(),
		Before:             cmd.Before,
		Action: func(cliCtx *cli.Context) error {
			return cmd.run(cliCtx.Context, cliCtx.App.Version)
		},
	}
}

func (c *cmdDevicesUploadLogs) Flags() []cli.Flag {
	flags := c.cmdDevices.Flags()
	flags = append(flags,
		&cli.DurationFlag{
			Name:        "timeout",
			Usage:       "Time to wait for blueprint uploading",
			Destination: &c.timeout,
			Value:       deviceUploadDefaultTimeout,
		},
		&cli.StringFlag{
			Name:        "operation-id",
			Usage:       "Uploading operation ID (optional)",
			Destination: &c.operationID,
		},
	)
	return flags
}

func (c *cmdDevicesUploadLogs) run(ctx context.Context, version string) error {
	if c.timeout != 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, c.timeout)
		defer cancel()
	}

	transport := cloudapi.NewCredentialsTransport(http.DefaultTransport, c.token, version)
	transport = cloudapi.NewCLIMessageWriterTransport(transport, &onceWriter{w: c.writer})
	client := cloudapi.NewClientWithURL(&http.Client{Transport: transport}, c.graphqlURL)

	if c.operationID != "" {
		return client.WriteOperationLogs(ctx, c.hardwareID, c.operationID, c.writeLog)
	}

	const lastOperationsNumber = 2
	return client.WriteLastOperationsLogs(ctx, c.hardwareID, lastOperationsNumber, c.writeLog)
}

func (c *cmdDevicesUploadLogs) writeLog(operationID string, l cloudapi.OperationLog) {
	fmt.Fprintf(c.writer, "[#%s] %s [%s] %s\n", operationID, l.CreatedAt, l.Severity, l.Payload)
}
