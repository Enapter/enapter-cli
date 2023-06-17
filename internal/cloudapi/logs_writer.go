package cloudapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"sync"

	"github.com/gorilla/websocket"
)

const (
	fieldChannel    = "channel"
	fieldHardwareID = "hardware_id"
	fieldRuleID     = "rule_id"
)

type LogsWriter struct {
	url        string
	identifier map[string]string
	writeLog   func(topic, message string)
	wsConnMu   sync.Mutex
	wsConn     *websocket.Conn
}

func NewDeviceLogsWriter(
	host, token, apiVersion, hardwareID string,
	writer func(topic, message string),
) (*LogsWriter, error) {
	identifier := map[string]string{
		fieldChannel:    "DeviceChannel",
		fieldHardwareID: hardwareID,
	}
	return newLogsWriter(host, token, apiVersion, identifier, writer)
}

func NewRuleLogsWriter(
	host, token, apiVersion, ruleID string,
	writer func(topic, message string),
) (*LogsWriter, error) {
	identifier := map[string]string{
		fieldChannel: "RuleChannel",
		fieldRuleID:  ruleID,
	}
	return newLogsWriter(host, token, apiVersion, identifier, writer)
}

func newLogsWriter(
	host, token, apiVersion string, identifier map[string]string,
	writer func(topic, message string),
) (*LogsWriter, error) {
	u, err := url.Parse(host)
	if err != nil {
		return nil, fmt.Errorf("parse url: %w", err)
	}

	q := u.Query()
	q.Set("token", token)
	q.Set("enapter_api_version", apiVersion)
	u.RawQuery = q.Encode()

	return &LogsWriter{
		url:        u.String(),
		identifier: identifier,
		writeLog:   writer,
	}, nil
}

func (l *LogsWriter) Run(ctx context.Context) error {
	defer l.close()

	go func() {
		<-ctx.Done()
		l.close()
	}()

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		if err := l.connect(ctx); err != nil {
			l.writeLog("connection", fmt.Sprintf("failed to connect: %s", err.Error()))
			continue
		}

		if err := l.subscribe(); err != nil {
			l.writeLog("connection", fmt.Sprintf("failed to subscribe: %s", err.Error()))
			continue
		}

		if err := l.readAndWriteLogs(ctx); err != nil {
			if errors.Is(err, ErrFinished) {
				return nil
			}
			l.writeLog("read_error", fmt.Sprintf("failed to read msg: %s", err.Error()))
			return err
		}
	}
}

func (l *LogsWriter) connect(ctx context.Context) error {
	wsConn, resp, err := websocket.DefaultDialer.DialContext(ctx, l.url, nil)
	if err != nil {
		return fmt.Errorf("websockets dial: %w", err)
	}
	defer resp.Body.Close()

	l.wsConnMu.Lock()
	defer l.wsConnMu.Unlock()
	l.wsConn = wsConn

	return nil
}

func (l *LogsWriter) close() {
	l.wsConnMu.Lock()
	defer l.wsConnMu.Unlock()

	if l.wsConn != nil {
		l.wsConn.Close()
	}
}

func (l *LogsWriter) subscribe() error {
	identifierBytes, err := json.Marshal(l.identifier)
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

	err = l.wsConn.WriteMessage(websocket.TextMessage, msgBytes)
	if err != nil {
		return fmt.Errorf("failed to subscribe: %w", err)
	}

	return nil
}

func (l *LogsWriter) readAndWriteLogs(ctx context.Context) error {
	for {
		msgType, msgBytes, err := l.wsConn.ReadMessage()
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		if err != nil {
			return err
		}

		if msgType != websocket.TextMessage {
			l.writeLog("read_error",
				fmt.Sprintf("skip unsupported message type [%d] (only text type [%d] supported)",
					msgType, websocket.TextMessage))
			continue
		}

		if err := l.process(msgBytes); err != nil {
			return err
		}
	}
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

type disconnectMessage struct {
	Type      string `json:"type"`
	Reason    string `json:"reason"`
	Reconnect bool   `json:"reconnect"`
}

func (l *LogsWriter) process(msgBytes []byte) error {
	baseMsg := logsBaseMessage{}
	if err := json.Unmarshal(msgBytes, &baseMsg); err != nil {
		errMsg := fmt.Sprintf("skip invalid message %s: %s", string(msgBytes), err.Error())
		l.writeLog("read_error", errMsg)
		return err
	}

	switch baseMsg.Type {
	case "ping":
	case "welcome", "confirm_subscription":
		l.writeLog("connection", baseMsg.Type)
	case "reject_subscription":
		l.writeLog("connection", "disconnected")
		return ErrFinished
	case "disconnect":
		return l.processDisconnect(msgBytes)
	case "message":
		l.processMessage(baseMsg)
	default:
		l.writeLog("unknown", string(msgBytes))
	}

	return nil
}

func (l *LogsWriter) processDisconnect(msgBytes []byte) error {
	var msg disconnectMessage
	if err := json.Unmarshal(msgBytes, &msg); err != nil {
		return nil
	}
	if msg.Reconnect {
		l.writeLog("connection",
			fmt.Sprintf("disconnected with reason: %s. Reconnecting...", msg.Reason))
		return nil
	}
	l.writeLog("connection",
		fmt.Sprintf("disconnected with reason: %s", msg.Reason))
	return ErrFinished
}

func (l *LogsWriter) processMessage(baseMsg logsBaseMessage) {
	var identifier map[string]string
	if err := json.Unmarshal([]byte(baseMsg.Identifier), &identifier); err != nil {
		l.writeLog("read_error",
			fmt.Sprintf("skip message with invalid identifier %s: %s", baseMsg.Identifier, err.Error()))
		return
	}
	if !mapsEqual(l.identifier, identifier) {
		l.writeLog("read_error",
			fmt.Sprintf("skip message with unknown identifier %+v", identifier))
		return
	}

	var msg logsMessage
	if err := json.Unmarshal(baseMsg.Message, &msg); err != nil {
		l.writeLog("read_error",
			fmt.Sprintf("skip invalid log message %s: %s", string(baseMsg.Message), err.Error()))
		return
	}
	l.writeLog(msg.Topic, msg.Payload)
}

func mapsEqual(m1, m2 map[string]string) bool {
	if len(m1) != len(m2) {
		return false
	}
	for k, v1 := range m1 {
		if v2, ok := m2[k]; !ok || v1 != v2 {
			return false
		}
	}
	return true
}
