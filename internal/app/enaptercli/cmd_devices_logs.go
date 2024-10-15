package enaptercli

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/urfave/cli/v2"
)

type cmdDevicesLogs struct {
	cmdDevices
	from       cli.Timestamp
	to         cli.Timestamp
	offset     int
	limit      int
	severity   string
	order      string
	showFilter string
}

func buildCmdDevicesLogs() *cli.Command {
	cmd := &cmdDevicesLogs{}
	return &cli.Command{
		Name:               "logs",
		Usage:              "Show device logs",
		CustomHelpTemplate: cmd.HelpTemplate(),
		Flags:              cmd.Flags(),
		Before:             cmd.Before,
		Action: func(cliCtx *cli.Context) error {
			return cmd.do(cliCtx.Context)
		},
	}
}

func (c *cmdDevicesLogs) Flags() []cli.Flag {
	flags := c.cmdDevices.Flags()
	return append(flags, &cli.TimestampFlag{
		Name:        "from",
		Aliases:     []string{"f"},
		Usage:       "from timestamp in rfc 3339 format (like 2006-01-02T15:04:05Z)",
		Destination: &c.from,
		Layout:      time.RFC3339,
	}, &cli.TimestampFlag{
		Name:        "to",
		Aliases:     []string{"t"},
		Usage:       "to timestamp in rfc 3339 format (like 2006-01-02T15:04:05Z)",
		Destination: &c.to,
		Layout:      time.RFC3339,
	}, &cli.IntFlag{
		Name:        "limit",
		Aliases:     []string{"l"},
		Usage:       "maximum number of logs to retrieve",
		Destination: &c.limit,
	}, &cli.IntFlag{
		Name:        "offset",
		Aliases:     []string{"o"},
		Usage:       "number of logs to skip on retrieve",
		Destination: &c.offset,
	}, &cli.StringFlag{
		Name:        "severity",
		Aliases:     []string{"s"},
		Usage:       "filter logs by severity",
		Destination: &c.severity,
	}, &cli.StringFlag{
		Name:        "order",
		Usage:       "order logs by criteria (received_at_asc[default], received_at_desc)",
		Destination: &c.order,
		Action: func(_ *cli.Context, v string) error {
			if v != "received_at_asc" && v != "received_at_desc" {
				return fmt.Errorf("%w: should be one of [received_at_asc, received_at_desc]", errUnsupportedFlagValue)
			}
			return nil
		},
	}, &cli.StringFlag{
		Name:        "show",
		Usage:       "filter logs by criteria (all[default], persist_only, temporary_only)",
		Destination: &c.showFilter,
		Action: func(_ *cli.Context, v string) error {
			if v != "all" && v != "persist_only" && v != "temporary_only" {
				return fmt.Errorf("%w: should be one of [all, persist_only, temporary_only]", errUnsupportedFlagValue)
			}
			return nil
		},
	})
}

func (c *cmdDevicesLogs) do(ctx context.Context) error {
	query := url.Values{}
	if c.from.Value() != nil {
		query.Add("received_at_from", c.from.Value().Format(time.RFC3339))
	}
	if c.to.Value() != nil {
		query.Add("received_at_to", c.to.Value().Format(time.RFC3339))
	}
	if c.offset > 0 {
		query.Add("offset", strconv.Itoa(c.offset))
	}
	if c.limit > 0 {
		query.Add("limit", strconv.Itoa(c.limit))
	}
	if c.severity != "" {
		query.Add("severity", c.severity)
	}
	if c.order != "" {
		query.Add("order", c.order)
	}
	if c.showFilter != "" {
		query.Add("show", c.showFilter)
	}

	return c.doHTTPRequest(ctx, doHTTPRequestParams{
		Method: http.MethodGet,
		Path:   "/logs",
		Query:  query,
		//nolint:bodyclose //body is closed in doHTTPRequest
		RespProcessor: okRespBodyProcessor(func(body io.Reader) error {
			_, err := c.parseAndDumpDeviceLogs(body)
			return err
		}),
	})
}
