package enaptercli

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"slices"
	"time"

	"github.com/urfave/cli/v2"
)

type cmdDeviceMonitor struct {
	cmdDevice
	deviceID       string
	includeRuntime bool
}

func buildCmdDeviceStream() *cli.Command {
	cmd := &cmdDeviceMonitor{}
	return &cli.Command{
		Name:               "monitor",
		Usage:              "Monitor device traffic",
		CustomHelpTemplate: cmd.CommandHelpTemplate(),
		Flags:              cmd.Flags(),
		Before:             cmd.Before,
		Action: func(cliCtx *cli.Context) error {
			return cmd.do(cliCtx.Context)
		},
	}
}

func (c *cmdDeviceMonitor) Flags() []cli.Flag {
	flags := c.cmdDevice.Flags()
	return append(flags, &cli.StringFlag{
		Name:        "device-id",
		Aliases:     []string{"d"},
		Usage:       "Device ID",
		Destination: &c.deviceID,
		Required:    true,
	}, &cli.BoolFlag{
		Name:        "include-runtime",
		Usage:       "Monitor device's runtime traffic too",
		Destination: &c.includeRuntime,
	})
}

func (c *cmdDeviceMonitor) do(ctx context.Context) error {
	deviceID, runtimeID, err := c.resolveDeviceIDs(ctx)
	if err != nil {
		return err
	}

	query := make(url.Values)
	query.Add("id.in", deviceID)
	if runtimeID != "" {
		query.Add("id.in", runtimeID)
	}

	return c.runWebSocket(ctx, runWebSocketParams{
		Path:  "/",
		Query: query,
		RespProcessor: func(r io.Reader) error {
			return c.process(r, deviceID, runtimeID)
		},
	})
}

func (c *cmdDeviceMonitor) resolveDeviceIDs(ctx context.Context) (string, string, error) {
	var resp struct {
		Device struct {
			ID            string `json:"id"`
			Type          string `json:"type"`
			Communication struct {
				Type       string `json:"type"`
				UpstreamID string `json:"upstream_id"`
			} `json:"communication"`
		} `json:"device"`
	}

	var query url.Values
	if c.includeRuntime {
		query = url.Values{"expand": {"communication"}}
	}

	if err := c.doHTTPRequest(ctx, doHTTPRequestParams{
		Method: http.MethodGet,
		Path:   c.deviceID,
		Query:  query,
		RespProcessor: func(r *http.Response) error {
			if r.StatusCode != http.StatusOK {
				return cli.Exit(parseRespErrorMessage(r), 1)
			}
			return json.NewDecoder(r.Body).Decode(&resp)
		},
	}); err != nil {
		return "", "", err
	}

	if !c.includeRuntime {
		return resp.Device.ID, "", nil
	}

	runtimeCommTypes := []string{"UCM_LUA"}
	if !slices.Contains(runtimeCommTypes, resp.Device.Communication.Type) {
		fmt.Fprintln(c.errWriter,
			"WARNING: device does not run on a runtime, --include-runtime is ignored")
		return resp.Device.ID, "", nil
	}

	return resp.Device.ID, resp.Device.Communication.UpstreamID, nil
}

type streamMessage struct {
	DeviceID   string          `json:"device_id"`
	ReceivedAt time.Time       `json:"received_at"`
	Timestamp  int64           `json:"timestamp"`
	Telemetry  json.RawMessage `json:"telemetry,omitempty"`
	Properties json.RawMessage `json:"properties,omitempty"`
	Log        *struct {
		Severity string `json:"severity"`
		Message  string `json:"message"`
	} `json:"log,omitempty"`
}

func (c *cmdDeviceMonitor) process(r io.Reader, deviceID, runtimeID string) error {
	var m streamMessage
	if err := json.NewDecoder(r).Decode(&m); err != nil {
		return fmt.Errorf("parse payload: %w", err)
	}

	color := c.messageColor(m)
	if color != "" {
		fmt.Fprint(c.writer, color)
	}

	fmt.Fprint(c.writer, c.messageTimestamp(m))
	fmt.Fprint(c.writer, " ")

	switch m.DeviceID {
	case runtimeID:
		fmt.Fprint(c.writer, "runtime ")
	case deviceID:
		fmt.Fprint(c.writer, "device  ")
	}

	switch {
	case m.Telemetry != nil:
		fmt.Fprint(c.writer, "telemetry:  ")
		fmt.Fprint(c.writer, string(m.Telemetry))
	case m.Properties != nil:
		fmt.Fprint(c.writer, "properties: ")
		fmt.Fprint(c.writer, string(m.Properties))
	case m.Log != nil:
		fmt.Fprint(c.writer, "logs:       ")
		fmt.Fprintf(c.writer, "[%s] %s", m.Log.Severity, m.Log.Message)
	}

	if color != "" {
		fmt.Fprint(c.writer, colorReset)
	}
	fmt.Fprintln(c.writer)

	return nil
}

func (c *cmdDeviceMonitor) messageColor(m streamMessage) string {
	if !c.colorize {
		return ""
	}
	if m.Log != nil {
		switch m.Log.Severity {
		case "warning":
			return colorYellow
		case "error":
			return colorRed
		}
	}
	return ""
}

func (c *cmdDeviceMonitor) messageTimestamp(m streamMessage) string {
	ts := m.ReceivedAt
	if ts.IsZero() {
		ts = time.Unix(m.Timestamp, 0)
	}
	return ts.UTC().Format(time.RFC3339)
}
