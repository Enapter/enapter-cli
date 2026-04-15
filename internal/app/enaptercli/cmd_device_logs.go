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

	"github.com/urfave/cli/v3"
)

type cmdDeviceLogs struct {
	cmdDevice
	deviceID      string
	follow        bool
	receivedAtGte time.Time
	receivedAtLt  time.Time
	offset        int
	limit         int
	severity      string
	order         string
	retention     string
}

func buildCmdDeviceLogs() *cli.Command {
	cmd := &cmdDeviceLogs{}
	return &cli.Command{
		Name:               "logs",
		Usage:              "Show device logs",
		CustomHelpTemplate: cmd.CommandHelpTemplate(),
		Flags:              cmd.Flags(),
		Before:             cmd.Before,
		Action: func(ctx context.Context, _ *cli.Command) error {
			return cmd.do(ctx)
		},
	}
}

func (c *cmdDeviceLogs) Flags() []cli.Flag {
	flags := c.cmdDevice.Flags()
	return append(flags, &cli.StringFlag{
		Name:        "device-id",
		Aliases:     []string{"d"},
		Usage:       "Device ID",
		Destination: &c.deviceID,
		Required:    true,
	}, &cli.BoolFlag{
		Name:        "follow",
		Aliases:     []string{"f"},
		Usage:       "Follow the log output",
		Destination: &c.follow,
	}, &cli.TimestampFlag{
		Name:        "received-at-gte",
		Usage:       "From timestamp in RFC 3339 format (e.g. 2006-01-02T15:04:05Z)",
		Destination: &c.receivedAtGte,
		Config:      cli.TimestampConfig{Layouts: []string{time.RFC3339}},
	}, &cli.TimestampFlag{
		Name:        "received-at-lt",
		Usage:       "To timestamp in RFC 3339 format (e.g. 2006-01-02T15:04:05Z)",
		Destination: &c.receivedAtLt,
		Config:      cli.TimestampConfig{Layouts: []string{time.RFC3339}},
	}, &cli.IntFlag{
		Name:        "limit",
		Aliases:     []string{"l"},
		Usage:       "Maximum number of logs to retrieve",
		Destination: &c.limit,
	}, &cli.IntFlag{
		Name:        "offset",
		Aliases:     []string{"o"},
		Usage:       "Number of logs to skip when retrieving",
		Destination: &c.offset,
	}, &cli.StringFlag{
		Name:        "severity",
		Aliases:     []string{"s"},
		Usage:       "Filter logs by severity",
		Destination: &c.severity,
	}, &cli.StringFlag{
		Name:        "order",
		Usage:       "Order logs by criteria (RECEIVED_AT_ASC[default], RECEIVED_AT_DESC)",
		Destination: &c.order,
		Action: func(_ context.Context, _ *cli.Command, v string) error {
			if v != "RECEIVED_AT_ASC" && v != "RECEIVED_AT_DESC" {
				return fmt.Errorf("%w: should be one of [RECEIVED_AT_ASC, RECEIVED_AT_DESC]", errUnsupportedFlagValue)
			}
			return nil
		},
	}, &cli.StringFlag{
		Name:        "retention",
		Usage:       "Filter logs by retention (ALL[default], PERSISTENT, EPHEMERAL)",
		Destination: &c.retention,
		Action: func(_ context.Context, _ *cli.Command, v string) error {
			if v != "ALL" && v != "PERSISTENT" && v != "EPHEMERAL" {
				return fmt.Errorf("%w: should be one of [ALL, PERSISTENT, EPHEMERAL]", errUnsupportedFlagValue)
			}
			return nil
		},
	})
}

func (c *cmdDeviceLogs) do(ctx context.Context) error {
	if c.follow {
		return c.doFollow(ctx)
	}
	return c.doList(ctx)
}

func (c *cmdDeviceLogs) doFollow(ctx context.Context) error {
	if !c.receivedAtGte.IsZero() {
		return cli.Exit("Option --received-at-gte is unsupported in follow mode.", 1)
	}
	if !c.receivedAtLt.IsZero() {
		return cli.Exit("Option --received-at-lt is unsupported in follow mode.", 1)
	}
	if c.offset > 0 {
		return cli.Exit("Option --offset is unsupported in follow mode.", 1)
	}
	if c.limit > 0 {
		return cli.Exit("Option --limit is unsupported in follow mode.", 1)
	}
	if c.order != "" {
		return cli.Exit("Option --order is unsupported in follow mode.", 1)
	}

	query := url.Values{}
	if c.severity != "" {
		query.Add("severity", c.severity)
	}
	if c.retention != "" {
		query.Add("retention", c.retention)
	}

	return c.runWebSocket(ctx, runWebSocketParams{
		Path:  "/" + c.deviceID + "/logs",
		Query: query,
		RespProcessor: func(r io.Reader) error {
			var msg struct {
				Timestamp  int64  `json:"timestamp"`
				ReceivedAt string `json:"received_at"`
				Log        struct {
					Severity string `json:"severity"`
					Message  string `json:"message"`
				} `json:"log"`
			}
			if err := json.NewDecoder(r).Decode(&msg); err != nil {
				return fmt.Errorf("parse payload: %w", err)
			}

			color := c.logColor(msg.Log.Severity)
			if color != "" {
				fmt.Fprint(c.writer, color)
			}
			fmt.Fprintf(c.writer, "%s [%s] %s\n", msg.ReceivedAt, msg.Log.Severity, msg.Log.Message)
			if color != "" {
				fmt.Fprint(c.writer, colorReset)
			}
			return nil
		},
	})
}

func (c *cmdDeviceLogs) doList(ctx context.Context) error {
	query := url.Values{}
	if !c.receivedAtGte.IsZero() {
		query.Add("received_at.gte", c.receivedAtGte.Format(time.RFC3339))
	}
	if !c.receivedAtLt.IsZero() {
		query.Add("received_at.lt", c.receivedAtLt.Format(time.RFC3339))
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
	if c.retention != "" {
		query.Add("retention", c.retention)
	}

	return c.doHTTPRequest(ctx, doHTTPRequestParams{
		Method: http.MethodGet,
		Path:   "/" + c.deviceID + "/logs",
		Query:  query,
		//nolint:bodyclose //body is closed in doHTTPRequest
		RespProcessor: okRespBodyProcessor(func(body io.Reader) error {
			var resp struct {
				Logs []struct {
					Timestamp  int64  `json:"timestamp"`
					ReceivedAt string `json:"received_at"`
					Severity   string `json:"severity"`
					Message    string `json:"message"`
				} `json:"logs"`
			}
			if err := json.NewDecoder(body).Decode(&resp); err != nil {
				return fmt.Errorf("parse response body: %w", err)
			}
			for _, l := range resp.Logs {
				color := c.logColor(l.Severity)
				if color != "" {
					fmt.Fprint(c.writer, color)
				}
				fmt.Fprintf(c.writer, "%s [%s] %s\n", l.ReceivedAt, l.Severity, l.Message)
				if color != "" {
					fmt.Fprint(c.writer, colorReset)
				}
			}
			return nil
		}),
	})
}

func (c *cmdDeviceLogs) logColor(severity string) string {
	if !c.colorize {
		return ""
	}
	switch severity {
	case "warning":
		return colorYellow
	case "error":
		return colorRed
	default:
		return ""
	}
}
