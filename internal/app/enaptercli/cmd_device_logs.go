package enaptercli

import (
	"context"
	"encoding/json"
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
	deviceID   string
	follow     bool
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
		CustomHelpTemplate: cmd.CommandHelpTemplate(),
		Flags:              cmd.Flags(),
		Before:             cmd.Before,
		Action: func(cliCtx *cli.Context) error {
			return cmd.do(cliCtx.Context)
		},
	}
}

func (c *cmdDevicesLogs) Flags() []cli.Flag {
	flags := c.cmdDevices.Flags()
	return append(flags, &cli.StringFlag{
		Name:        "device-id",
		Aliases:     []string{"d"},
		Usage:       "device ID",
		Destination: &c.deviceID,
		Required:    true,
	}, &cli.BoolFlag{
		Name:        "follow",
		Aliases:     []string{"f"},
		Usage:       "follow the log output",
		Destination: &c.follow,
	}, &cli.TimestampFlag{
		Name:        "from",
		Usage:       "from timestamp in RFC 3339 format (e.g. 2006-01-02T15:04:05Z)",
		Destination: &c.from,
		Layout:      time.RFC3339,
	}, &cli.TimestampFlag{
		Name:        "to",
		Usage:       "to timestamp in RFC 3339 format (e.g. 2006-01-02T15:04:05Z)",
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
		Usage:       "number of logs to skip when retrieving",
		Destination: &c.offset,
	}, &cli.StringFlag{
		Name:        "severity",
		Aliases:     []string{"s"},
		Usage:       "filter logs by severity",
		Destination: &c.severity,
	}, &cli.StringFlag{
		Name:        "order",
		Usage:       "order logs by criteria (RECEIVED_AT_ASC[default], RECEIVED_AT_DESC)",
		Destination: &c.order,
		Action: func(_ *cli.Context, v string) error {
			if v != "RECEIVED_AT_ASC" && v != "RECEIVED_AT_DESC" {
				return fmt.Errorf("%w: should be one of [RECEIVED_AT_ASC, RECEIVED_AT_DESC]", errUnsupportedFlagValue)
			}
			return nil
		},
	}, &cli.StringFlag{
		Name:        "show",
		Usage:       "filter logs by criteria (ALL[default], PERSISTED_ONLY, TEMPORARY_ONLY)",
		Destination: &c.showFilter,
		Action: func(_ *cli.Context, v string) error {
			if v != "ALL" && v != "PERSISTED_ONLY" && v != "TEMPORARY_ONLY" {
				return fmt.Errorf("%w: should be one of [ALL, PERSISTED_ONLY, TEMPORARY_ONLY]", errUnsupportedFlagValue)
			}
			return nil
		},
	})
}

func (c *cmdDevicesLogs) do(ctx context.Context) error {
	if c.follow {
		return c.doFollow(ctx)
	}
	return c.doList(ctx)
}

func (c *cmdDevicesLogs) doFollow(ctx context.Context) error {
	if c.from.Value() != nil {
		return cli.Exit("Option received_at_from is unsupported in follow mode.", 1)
	}
	if c.to.Value() != nil {
		return cli.Exit("Option received_at_to is unsupported in follow mode.", 1)
	}
	if c.offset > 0 {
		return cli.Exit("Option offset is unsupported in follow mode.", 1)
	}
	if c.limit > 0 {
		return cli.Exit("Option limit is unsupported in follow mode.", 1)
	}
	if c.order != "" {
		return cli.Exit("Option order is unsupported in follow mode.", 1)
	}

	query := url.Values{}
	if c.severity != "" {
		query.Add("severity", c.severity)
	}
	if c.showFilter != "" {
		query.Add("show", c.showFilter)
	}

	path := fmt.Sprintf("/devices/%s/logs", c.deviceID)

	return c.runWebSocket(ctx, runWebSocketParams{
		Path:  path,
		Query: query,
		RespProcessor: func(r io.Reader) error {
			var msg struct {
				ReceivedAt string `json:"received_at"`
				Timestamp  string `json:"timestamp"`
				Severity   string `json:"severity"`
				Message    string `json:"message"`
			}
			if err := json.NewDecoder(r).Decode(&msg); err != nil {
				return fmt.Errorf("parse payload: %w", err)
			}
			fmt.Fprintf(c.writer, "%s [%s] %s\n", msg.ReceivedAt, msg.Severity, msg.Message)
			return nil
		},
	})
}

func (c *cmdDevicesLogs) doList(ctx context.Context) error {
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
		Path:   "/" + c.deviceID + "/logs",
		Query:  query,
		//nolint:bodyclose //body is closed in doHTTPRequest
		RespProcessor: okRespBodyProcessor(func(body io.Reader) error {
			var resp struct {
				Logs []struct {
					ReceivedAt string `json:"received_at"`
					Timestamp  string `json:"timestamp"`
					Severity   string `json:"severity"`
					Message    string `json:"message"`
				} `json:"logs"`
			}
			if err := json.NewDecoder(body).Decode(&resp); err != nil {
				return fmt.Errorf("parse response body: %w", err)
			}
			for _, l := range resp.Logs {
				fmt.Fprintf(c.writer, "%s [%s] %s\n", l.ReceivedAt, l.Severity, l.Message)
			}
			return nil
		}),
	})
}
