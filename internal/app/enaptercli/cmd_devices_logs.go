package enaptercli

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/urfave/cli/v2"
)

type cmdDevicesLogs struct {
	cmdDevices

	urlStr    string
	urlStrLog string
	wsConnMu  sync.Mutex
	wsConn    *websocket.Conn
}

func buildCmdDevicesLogs() *cli.Command {
	cmd := &cmdDevicesLogs{}

	var wsAPIURL string
	flags := cmd.Flags()
	flags = append(flags, &cli.StringFlag{
		Name:        "ws-api-url",
		EnvVars:     []string{"ENAPTER_WS_API_URL"},
		Hidden:      true,
		Destination: &wsAPIURL,
	})

	return &cli.Command{
		Name:               "logs",
		Usage:              "Stream logs from a device",
		CustomHelpTemplate: cmd.DevicesCmdHelpTemplate(),
		Flags:              flags,
		Before: func(cliCtx *cli.Context) error {
			if err := cmd.Before(cliCtx); err != nil {
				return err
			}

			u := &url.URL{}
			if wsAPIURL != "" {
				var err error
				u, err = url.Parse(wsAPIURL)
				if err != nil {
					return fmt.Errorf("failed to parse url path %q: %w", wsAPIURL, err)
				}
			} else {
				u.Host = cmd.apiHost
				u.Scheme = "wss"
				u.Path = "/cable"
			}

			q := url.Values{
				"token":               {cmd.token},
				"enapter_api_version": {cliCtx.App.Version},
			}
			u.RawQuery = q.Encode()

			cmd.urlStr = u.String()

			q["token"] = []string{"---"}
			u.RawQuery = q.Encode()
			cmd.urlStrLog = u.String()
			return nil
		},
		Action: func(cliCtx *cli.Context) error {
			return cmd.dumpLogs(cliCtx.Context)
		},
	}
}

func (c *cmdDevicesLogs) dumpLogs(ctx context.Context) error {
	defer c.closeConnect()

	go func() {
		<-ctx.Done()
		c.closeConnect()
	}()

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		if err := c.connect(ctx); err != nil {
			c.writeLog("connection", fmt.Sprintf("failed to connect: %s", err.Error()))
			continue
		}

		if err := c.subscribe(); err != nil {
			c.writeLog("connection", fmt.Sprintf("failed to subscribe: %s", err.Error()))
			continue
		}

		if err := c.readAndWriteLogs(ctx); err != nil {
			if errors.Is(err, errFinished) {
				return nil
			}
			c.writeLog("read_error", fmt.Sprintf("failed to read msg: %s", err.Error()))
			return err
		}
	}
}

func (c *cmdDevicesLogs) connect(ctx context.Context) error {
	wsConn, resp, err := websocket.DefaultDialer.DialContext(ctx, c.urlStr, nil)
	if err != nil {
		return fmt.Errorf("failed to dial %q: %w", c.urlStrLog, err)
	}
	defer resp.Body.Close()

	c.wsConnMu.Lock()
	defer c.wsConnMu.Unlock()
	c.wsConn = wsConn

	return nil
}

func (c *cmdDevicesLogs) closeConnect() {
	c.wsConnMu.Lock()
	defer c.wsConnMu.Unlock()

	if c.wsConn != nil {
		c.wsConn.Close()
	}
}

type identifierMsg struct {
	Channel    string `json:"channel"`
	HardwareID string `json:"hardware_id"`
}

func (c *cmdDevicesLogs) subscribe() error {
	identifier := identifierMsg{
		Channel:    "DeviceChannel",
		HardwareID: c.hardwareID,
	}

	identifierBytes, err := json.Marshal(&identifier)
	if err != nil {
		return fmt.Errorf("failed to marshal sm: %w", err)
	}

	msg := struct {
		Command    string `json:"command"`
		Identifier string `json:"identifier"`
	}{
		Command:    "subscribe",
		Identifier: string(identifierBytes),
	}

	msgBytes, err := json.Marshal(&msg)
	if err != nil {
		return fmt.Errorf("failed to marshal m: %w", err)
	}

	err = c.wsConn.WriteMessage(websocket.TextMessage, msgBytes)
	if err != nil {
		return fmt.Errorf("failed to subscribe: %w", err)
	}

	return nil
}

func (c *cmdDevicesLogs) readAndWriteLogs(ctx context.Context) error {
	for {
		msgType, msgBytes, err := c.wsConn.ReadMessage()
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		if err != nil {
			return err
		}

		if msgType != websocket.TextMessage {
			c.writeLog("read_error",
				fmt.Sprintf("skip unsupported message type [%d] (only text type [%d] supported)",
					msgType, websocket.TextMessage))
			continue
		}

		if err := c.process(msgBytes); err != nil {
			return err
		}
	}
}

func (c *cmdDevicesLogs) writeLog(topic, msg string) {
	fmt.Fprintf(c.writer, "[%s] %s\n", topic, msg)
}

type logsBaseMessage struct {
	Type       string          `json:"type"`
	Identifier string          `json:"identifier"`
	Message    json.RawMessage `json:"message"`
}

type logsMessage struct {
	Topic   string `json:"topic"`
	Payload string `json:"payload"`
}

type disconnectMsg struct {
	Type      string `json:"type"`
	Reason    string `json:"reason"`
	Reconnect bool   `json:"reconnect"`
}

func (c *cmdDevicesLogs) process(msgBytes []byte) error {
	baseMsg := logsBaseMessage{}
	if err := json.Unmarshal(msgBytes, &baseMsg); err != nil {
		errMsg := fmt.Sprintf("skip invalid message %s: %s", string(msgBytes), err.Error())
		c.writeLog("read_error", errMsg)
		return err
	}

	switch baseMsg.Type {
	case "ping":
	case "welcome", "confirm_subscription":
		c.writeLog("connection", baseMsg.Type)
	case "reject_subscription":
		c.writeLog("connection", "disconnected")
		return errFinished
	case "disconnect":
		var disconnectMsg disconnectMsg
		if err := json.Unmarshal(msgBytes, &disconnectMsg); err != nil {
			return nil
		}
		if disconnectMsg.Reconnect {
			c.writeLog("connection",
				fmt.Sprintf("disconnected with reason: %s. Reconnecting...", disconnectMsg.Reason))
			return nil
		}
		c.writeLog("connection",
			fmt.Sprintf("disconnected with reason: %s", disconnectMsg.Reason))
		return errFinished
	case "message":
		var identifier identifierMsg
		if err := json.Unmarshal([]byte(baseMsg.Identifier), &identifier); err != nil {
			c.writeLog("read_error",
				fmt.Sprintf("skip message with invalid identifier %s: %s", baseMsg.Identifier, err.Error()))
			return nil
		}
		if identifier.Channel != "DeviceChannel" || identifier.HardwareID != c.hardwareID {
			c.writeLog("read_error",
				fmt.Sprintf("skip message with unknown identifier %+v", identifier))
			return nil
		}
		logsMsg := logsMessage{}
		if err := json.Unmarshal(baseMsg.Message, &logsMsg); err != nil {
			c.writeLog("read_error",
				fmt.Sprintf("skip invalid log message %s: %s", string(baseMsg.Message), err.Error()))
			return nil
		}
		c.writeLog(logsMsg.Topic, logsMsg.Payload)
	default:
		c.writeLog("unknown", string(msgBytes))
		return nil
	}

	return nil
}
