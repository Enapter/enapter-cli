package enaptercli

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/gorilla/websocket"
	"github.com/urfave/cli/v3"
	"golang.org/x/term"
)

type cmdDeviceRunTerminal struct {
	cmdDevice
	deviceID  string
	winWidth  int
	winHeight int
}

func buildCmdDeviceRunTerminal() *cli.Command {
	cmd := &cmdDeviceRunTerminal{}
	return &cli.Command{
		Name:  "run-terminal",
		Usage: "Run new remote terminal session",
		Description: "Remote terminal feature should be enabled in gateway settings. " +
			"Use Ctrl+] sequence to force connection close.",
		CustomHelpTemplate: cmd.CommandHelpTemplate(),
		Flags:              cmd.Flags(),
		Before:             cmd.Before,
		Action: func(ctx context.Context, _ *cli.Command) error {
			return cmd.do(ctx)
		},
	}
}

func (c *cmdDeviceRunTerminal) Flags() []cli.Flag {
	flags := c.cmdDevice.Flags()
	return append(flags, &cli.StringFlag{
		Name:        "device-id",
		Aliases:     []string{"d"},
		Usage:       "Gateway device ID",
		Destination: &c.deviceID,
		Required:    true,
	})
}

func (c *cmdDeviceRunTerminal) do(ctx context.Context) error {
	fin := os.Stdin
	fd := int(fin.Fd())
	if !term.IsTerminal(fd) {
		return cli.Exit("Standard input should be a terminal.", 1)
	}

	var credentials struct {
		ChannelID    string `json:"channel_id"`
		Token        string `json:"token"`
		WebSocketURL string `json:"websocket_url"`
	}
	if err := c.doHTTPRequest(ctx, doHTTPRequestParams{
		Method: http.MethodPost,
		Path:   "/" + c.deviceID + "/run_terminal",
		RespProcessor: func(r *http.Response) error {
			if r.StatusCode != http.StatusOK {
				return cli.Exit(parseRespErrorMessage(r), 1)
			}
			return json.NewDecoder(r.Body).Decode(&credentials)
		},
	}); err != nil {
		return err
	}

	url, err := url.Parse(credentials.WebSocketURL)
	if err != nil {
		return fmt.Errorf("parse url: %w", err)
	}
	headers := make(http.Header)
	headers.Set("Authorization", "Bearer "+credentials.Token)

	conn, err := c.dialWebSocket(ctx, url, headers)
	if err != nil {
		return fmt.Errorf("dial websocket: %w", err)
	}
	defer conn.Close()

	oldState, err := term.MakeRaw(fd)
	if err != nil {
		return fmt.Errorf("make raw terminal: %w", err)
	}
	defer func() { _ = term.Restore(fd, oldState) }()

	// TODO: wait for pong?
	if err := c.writePing(conn, credentials.ChannelID); err != nil {
		return fmt.Errorf("ping: %w", err)
	}

	fmt.Fprint(c.writer, "Use Ctrl+] to terminate the session.\r\n\r\n")

	errCh := make(chan error)
	stdinCh := make(chan byte)
	go func() { errCh <- c.runFileReader(ctx, fin, stdinCh) }()
	go func() { errCh <- c.runConnReader(ctx, conn) }()

	return c.run(ctx, conn, credentials.ChannelID, stdinCh, errCh, fd)
}

func (c *cmdDeviceRunTerminal) runFileReader(
	ctx context.Context, f *os.File, ch chan<- byte,
) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		var buf [1]byte
		n, err := f.Read(buf[:])
		if err != nil {
			return err
		}

		if n > 0 {
			select {
			case <-ctx.Done():
				return nil
			case ch <- buf[0]:
			}
		}
	}
}

func (c *cmdDeviceRunTerminal) runConnReader(
	ctx context.Context, conn *websocket.Conn,
) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		var msg struct {
			Data string `json:"data"`
		}
		if err := conn.ReadJSON(&msg); err != nil {
			return fmt.Errorf("read: %w", err)
		}

		var data []string
		if err := json.Unmarshal([]byte(msg.Data), &data); err != nil {
			return fmt.Errorf("unmarhal data: %w", err)
		}

		if len(data) == 0 {
			return cli.Exit("Unexpected payload from server.", 1)
		}
		if data[0] == "exit" {
			return nil
		}
		if data[0] == "stdin" && len(data) > 1 {
			fmt.Fprint(c.writer, data[1])
		}
	}
}

func (c *cmdDeviceRunTerminal) run(
	ctx context.Context, conn *websocket.Conn, channelID string,
	stdinCh <-chan byte, errCh <-chan error, fd int,
) error {
	const keepAliveInterval = 30 * time.Second
	const updateSizeInterval = time.Second

	keepAliveTicker := time.NewTicker(keepAliveInterval)
	updateSizeTicker := time.NewTicker(updateSizeInterval)

	for {
		select {
		case <-ctx.Done():
			return nil
		case err := <-errCh:
			return err
		default:
		}

		select {
		case <-ctx.Done():
			return nil
		case err := <-errCh:
			return err
		case <-keepAliveTicker.C:
			if err := c.writeKeepalive(conn, channelID); err != nil {
				return err
			}
		case <-updateSizeTicker.C:
			if err := c.writeSetSize(conn, channelID, fd); err != nil {
				return err
			}
		case b := <-stdinCh:
			const GS = 29 // ^]
			if b == GS {
				return cli.Exit("Exiting session.", 0)
			}
			if err := c.writeStdin(conn, channelID, b); err != nil {
				return err
			}
		}
	}
}

func (c *cmdDeviceRunTerminal) writePing(conn *websocket.Conn, channelID string) error {
	return conn.WriteJSON(map[string]any{
		"channel": channelID,
		"data":    `["ping"]`,
	})
}

func (c *cmdDeviceRunTerminal) writeStdin(conn *websocket.Conn, channelID string, b byte) error {
	data, err := json.Marshal([]string{"stdin", string(b)})
	if err != nil {
		return err
	}
	return conn.WriteJSON(map[string]any{
		"channel": channelID,
		"data":    string(data),
	})
}

func (c *cmdDeviceRunTerminal) writeKeepalive(conn *websocket.Conn, channelID string) error {
	return conn.WriteJSON(map[string]any{
		"channel": channelID,
		"data":    `["keepalive_ping"]`,
	})
}

func (c *cmdDeviceRunTerminal) writeSetSize(conn *websocket.Conn, channelID string, fd int) error {
	w, h, err := term.GetSize(fd)
	if err != nil {
		// FIXME: error on Windows
		return nil
	}
	if c.winWidth == w && c.winHeight == h {
		return nil
	}

	if err := conn.WriteJSON(map[string]any{
		"channel": channelID,
		"data":    fmt.Sprintf("[%q, %d, %d]", "set_size", h, w),
	}); err != nil {
		return err
	}

	c.winWidth = w
	c.winHeight = h
	return nil
}
