package enaptercli

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/shurcooL/graphql"
	"github.com/urfave/cli/v2"
)

type cmdDevicesUploadLogs struct {
	cmdDevicesUploadCommon
}

func buildCmdDevicesUploadLogs() *cli.Command {
	cmd := &cmdDevicesUploadLogs{}

	var operationID string

	flags := cmd.Flags()
	flags = append(flags, &cli.StringFlag{
		Name:        "operation-id",
		Usage:       "Uploading operation ID (optional)",
		Destination: &operationID,
	})

	return &cli.Command{
		Name:               "upload-logs",
		Usage:              "Show blueprint uploading logs",
		CustomHelpTemplate: cmd.DevicesCmdHelpTemplate(),
		Flags:              flags,
		Before:             cmd.Before,
		Action: func(cliCtx *cli.Context) error {
			onceWriter := &onceWriter{w: cmd.writer}
			return cmd.logs(cliCtx.Context, operationID, onceWriter, cliCtx.App.Version)
		},
	}
}

func (c *cmdDevicesUploadLogs) logs(
	ctx context.Context, operationID string, onceWriter *onceWriter, version string,
) error {
	if c.timeout != 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, c.timeout)
		defer cancel()
	}

	extraHeaders := map[string][]string{
		"Authorization":         {"Bearer " + c.token},
		"X-ENAPTER-CLI-VERSION": {version},
	}
	httpClient := &http.Client{
		Transport: extraHeaderRoundTripper{
			tripper: cliMessageRoundTripper{
				tripper: http.DefaultTransport,
				writer:  onceWriter,
			},
			extraHeaders: extraHeaders,
		},
	}

	client := graphql.NewClient(c.apiURL, httpClient)

	if operationID != "" {
		return c.logsOperationID(ctx, client, operationID)
	}
	return c.logsLastTwo(ctx, client)
}

type blueprintUpdateOperationQuery struct {
	Device *struct {
		BlueprintUpdateOperation *struct {
			Status string
			Logs   struct {
				Edges []struct {
					Cursor graphql.String
					Node   logNode
				}
			} `graphql:"logs(after: $after_cursor)"`
		} `graphql:"blueprintUpdateOperation(id: $operation_id)"`
	} `graphql:"device(hardwareId: $hardware_id)"`
}

type logNode struct {
	Payload   string
	CreatedAt string
	Severity  string
}

func (c *cmdDevicesUploadLogs) logsOperationID(
	ctx context.Context, client *graphql.Client, operationID string,
) error {
	v := map[string]interface{}{
		"after_cursor": graphql.String(""),
		"hardware_id":  c.hardwareID,
		"operation_id": operationID,
	}

	for {
		var q blueprintUpdateOperationQuery
		if err := client.Query(ctx, &q, v); err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				err = errRequestTimedOut
			}
			return fmt.Errorf("failed to send request: %w", err)
		}

		if q.Device == nil {
			return fmt.Errorf("%w: device not found", errFinishedWithError)
		}

		if q.Device.BlueprintUpdateOperation == nil {
			return fmt.Errorf("%w: operation not found", errFinishedWithError)
		}

		status := q.Device.BlueprintUpdateOperation.Status
		logs := q.Device.BlueprintUpdateOperation.Logs

		for _, e := range logs.Edges {
			v["after_cursor"] = e.Cursor
			c.writeLog(operationID, e.Node)
		}

		if len(logs.Edges) == 0 {
			switch status {
			case "SUCCEEDED":
				return nil
			case "ERROR":
				return errLogStatusError
			}
		}

		const logRequestPeriod = 100 * time.Millisecond
		time.Sleep(logRequestPeriod)
	}
}

type blueprintUpdateOperationsQuery struct {
	Device *struct {
		BlueprintUpdateOperations struct {
			Nodes []struct {
				ID string
			}
		} `graphql:"blueprintUpdateOperations(last: $last_int)"`
	} `graphql:"device(hardwareId: $hardware_id)"`
}

func (c *cmdDevicesUploadLogs) logsLastTwo(ctx context.Context, client *graphql.Client) error {
	lastOperationsNumber := 2
	v := map[string]interface{}{
		"hardware_id": c.hardwareID,
		"last_int":    graphql.Int(lastOperationsNumber),
	}

	var q blueprintUpdateOperationsQuery
	if err := client.Query(ctx, &q, v); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			err = errRequestTimedOut
		}
		return fmt.Errorf("failed to send request: %w", err)
	}

	if q.Device == nil {
		return fmt.Errorf("%w: device not found", errFinishedWithError)
	}

	for _, op := range q.Device.BlueprintUpdateOperations.Nodes {
		if err := c.logsOperationID(ctx, client, op.ID); err != nil {
			if errors.Is(err, errLogStatusError) {
				continue
			}
			return err
		}
	}

	return nil
}

func (c *cmdDevicesUploadLogs) writeLog(operationID string, n logNode) {
	fmt.Fprintf(c.writer, "[#%s] %s [%s] %s\n", operationID, n.CreatedAt, n.Severity, n.Payload)
}
