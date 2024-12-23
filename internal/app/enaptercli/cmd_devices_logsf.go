package enaptercli

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/urfave/cli/v3"
)

type cmdDevicesLogsf struct {
	cmdDevices
}

func buildCmdDevicesLogsf() *cli.Command {
	cmd := &cmdDevicesLogsf{}
	return &cli.Command{
		Name:               "logsf",
		Usage:              "Follow device logs",
		CustomHelpTemplate: cmd.HelpTemplate(),
		Flags:              cmd.Flags(),
		Before:             cmd.Before,
		Action: func(ctx context.Context, cm *cli.Command) error {
			return cmd.do(ctx)
		},
	}
}

func (c *cmdDevicesLogsf) do(ctx context.Context) error {
	const singleRequestLimit = 10

	query := url.Values{}
	query.Add("received_at_from", time.Now().Add(-time.Hour).UTC().Format(time.RFC3339))
	query.Add("order", "received_at_asc")
	query.Add("limit", strconv.Itoa(singleRequestLimit))

	offset := 0

	for {
		retryNow := false
		err := c.doHTTPRequest(ctx, doHTTPRequestParams{
			Method: http.MethodGet,
			Path:   "/logs",
			Query:  query,
			//nolint:bodyclose //body is closed in doHTTPRequest
			RespProcessor: okRespBodyProcessor(func(body io.Reader) error {
				n, err := c.parseAndDumpDeviceLogs(body)
				retryNow = n == singleRequestLimit
				offset += n
				query.Set("offset", strconv.Itoa(offset))
				return err
			}),
		})
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return nil
			}
			fmt.Fprintf(c.writer, "Failed to retrieve logs: %s\n", err)
			continue
		}

		if !retryNow {
			time.Sleep(time.Second)
		}
	}
}
