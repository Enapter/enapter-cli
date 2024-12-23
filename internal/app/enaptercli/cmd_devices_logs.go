package enaptercli

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/urfave/cli/v3"
)

type cmdDevicesLogs struct {
	cmdDevices
	from       time.Time
	to         time.Time
	offset     int64
	limit      int64
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
		Action: func(ctx context.Context, cm *cli.Command) error {
			return cmd.do(ctx)
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
		Config: cli.TimestampConfig{
			Layouts: []string{time.RFC3339},
		},
	}, &cli.TimestampFlag{
		Name:        "to",
		Aliases:     []string{"t"},
		Usage:       "to timestamp in rfc 3339 format (like 2006-01-02T15:04:05Z)",
		Destination: &c.to,
		Config: cli.TimestampConfig{
			Layouts: []string{time.RFC3339},
		},
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
		Action: func(_ context.Context, cm *cli.Command, v string) error {
			if v != "received_at_asc" && v != "received_at_desc" {
				return fmt.Errorf("%w: should be one of [received_at_asc, received_at_desc]", errUnsupportedFlagValue)
			}
			return nil
		},
	}, &cli.StringFlag{
		Name:        "show",
		Usage:       "filter logs by criteria (all[default], persist_only, temporary_only)",
		Destination: &c.showFilter,
		Action: func(_ context.Context, cm *cli.Command, v string) error {
			if v != "all" && v != "persist_only" && v != "temporary_only" {
				return fmt.Errorf("%w: should be one of [all, persist_only, temporary_only]", errUnsupportedFlagValue)
			}
			return nil
		},
	})
}

func (c *cmdDevicesLogs) do(ctx context.Context) error {
	query := url.Values{}
	if !c.from.IsZero() {
		query.Add("received_at_from", c.from.Format(time.RFC3339))
	}
	if !c.to.IsZero() {
		query.Add("received_at_to", c.to.Format(time.RFC3339))
	}
	if c.offset > 0 {
		query.Add("offset", strconv.FormatInt(c.offset, 10))
	}
	if c.limit > 0 {
		query.Add("limit", strconv.FormatInt(c.limit, 10))
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
