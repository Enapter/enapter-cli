package enaptercli_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/bxcodec/faker/v3"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/require"
)

func TestDeviceLogs(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		inputFileName := "testdata/device_logs/simple/input"
		untilLinePrefix := "[telemetry]"
		expectedFileName := "testdata/device_logs/simple/output"
		testDeviceLogs(t, inputFileName, untilLinePrefix, expectedFileName)
	})

	t.Run("invalid token", func(t *testing.T) {
		inputFileName := "testdata/device_logs/disconnect/invalid_token/input"
		untilLinePrefix := "[connection]"
		expectedFileName := "testdata/device_logs/disconnect/invalid_token/output"
		testDeviceLogs(t, inputFileName, untilLinePrefix, expectedFileName)
	})

	t.Run("device not found", func(t *testing.T) {
		inputFileName := "testdata/device_logs/disconnect/device_not_found/input"
		untilLinePrefix := "[connection] disconnected"
		expectedFileName := "testdata/device_logs/disconnect/device_not_found/output"
		testDeviceLogs(t, inputFileName, untilLinePrefix, expectedFileName)
	})
}

func testDeviceLogs(t *testing.T, inputFileName, untilLinePrefix, expectedFileName string) {
	token := faker.Word()
	const hardwareID = "SIM-WTM"
	handleErrCh := make(chan string)

	wsPath, srv := startWsServer(t, inputFileName, token, hardwareID, handleErrCh)
	defer srv.Close()

	args := strings.Split("enapter devices logs --token", " ")
	args = append(args, token, "--hardware-id", hardwareID, "--ws-api-url", wsPath)
	app := startTestApp(args...)
	defer app.Stop()

	actual := readOutputUntilLineOrError(t, app.Stdout(), untilLinePrefix, handleErrCh)
	if update {
		err := ioutil.WriteFile(expectedFileName, []byte(actual), 0600)
		require.NoError(t, err)
	}

	expected, err := ioutil.ReadFile(expectedFileName)
	require.NoError(t, err)
	require.Equal(t, string(expected), actual)

	app.Stop()
	appErr := app.Wait()
	require.NoError(t, appErr)

	restOutput, err := ioutil.ReadAll(app.Stdout())
	require.NoError(t, err)
	require.Empty(t, string(restOutput))
}

func readOutputUntilLineOrError(
	t *testing.T, r *lineBuffer, prefix string, wsHandleErrCh <-chan string,
) string {
	t.Helper()

	readStr, readErr := startBackgroundReadUntilLine(r, prefix)

	timer := time.NewTimer(5 * time.Second)
	select {
	case <-timer.C:
		require.Fail(t, "read output timed out")
	case errStr := <-wsHandleErrCh:
		require.Failf(t, "ws handler finished with error", errStr)
	case err := <-readErr:
		require.Failf(t, "read finished with error", err.Error())
	case s := <-readStr:
		return s
	}

	return ""
}

func startBackgroundReadUntilLine(r *lineBuffer, prefix string) (<-chan string, <-chan error) {
	readStr := make(chan string, 1)
	readErr := make(chan error, 1)

	go func() {
		buf := strings.Builder{}

		for {
			s, err := r.ReadLine()
			if err != nil {
				readErr <- err
				return
			}

			buf.WriteString(s)

			if strings.HasPrefix(s, prefix) {
				readStr <- buf.String()
				return
			}
		}
	}()

	return readStr, readErr
}

type closer interface {
	Close()
}

func startWsServer(
	t *testing.T, inputFileName, token, hardwareID string, handleErrCh chan<- string,
) (string, closer) {
	t.Helper()

	msgsBytes, err := ioutil.ReadFile(inputFileName)
	require.NoError(t, err)

	srv := httptest.NewServer(deviceLogsHandler(token, hardwareID, msgsBytes, handleErrCh))

	u, err := url.Parse(srv.URL)
	require.NoError(t, err)

	u.Scheme = "ws"
	return u.String(), srv
}

//nolint:funlen // because contains a lot of simple logged checks.
func deviceLogsHandler(token, hardwareID string, msgs []byte, handleErrCh chan<- string) http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {
		reqToken := r.URL.Query().Get("token")
		if reqToken != token {
			w.WriteHeader(http.StatusBadRequest)
			handleErrCh <- fmt.Sprintf("unexpected token %q, should be %q", reqToken, token)
			return
		}

		upgrader := websocket.Upgrader{}
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			handleErrCh <- fmt.Sprintf("failed to upgrade: %s", err)
			return
		}

		msgType, msgBytes, err := c.ReadMessage()
		if err != nil {
			handleErrCh <- fmt.Sprintf("failed to read subscribe message: %s", err)
			return
		}

		if msgType != websocket.TextMessage {
			handleErrCh <- fmt.Sprintf("subscribe message should be text type [%d], but [%d]",
				websocket.TextMessage, msgType)
			c.Close()
			return
		}

		subMsg := struct {
			Command    string `json:"command"`
			Identifier string `json:"identifier"`
		}{}
		if err := json.Unmarshal(msgBytes, &subMsg); err != nil {
			handleErrCh <- fmt.Sprintf("faild to unmarshall subsribe message %q: %s", string(msgBytes), err.Error())
			c.Close()
			return
		}

		if subMsg.Command != "subscribe" {
			handleErrCh <- fmt.Sprintf("this is not subscribe message, but %q", subMsg.Command)
			c.Close()
			return
		}

		msgIdentifier := struct {
			Channel    string `json:"channel"`
			HardwareID string `json:"hardware_id"`
		}{}
		if err := json.Unmarshal([]byte(subMsg.Identifier), &msgIdentifier); err != nil {
			handleErrCh <- fmt.Sprintf("faild to unmarshall subsribe message identifier %q: %s",
				subMsg.Identifier, err.Error())
			c.Close()
			return
		}

		if msgIdentifier.Channel != "DeviceChannel" {
			handleErrCh <- fmt.Sprintf("subsribe message identifier shoud have channel %q, but %q",
				"DeviceChannel", msgIdentifier.Channel)
			c.Close()
			return
		}

		if msgIdentifier.HardwareID != hardwareID {
			handleErrCh <- fmt.Sprintf("subsribe message identifier shoud have hardware_id %q, but %q",
				hardwareID, msgIdentifier.HardwareID)
			c.Close()
			return
		}

		msgss := bytes.Split(msgs, []byte{'\n'})
		for _, m := range msgss {
			if len(m) == 0 {
				continue
			}
			if err := c.WriteMessage(websocket.TextMessage, m); err != nil {
				handleErrCh <- fmt.Sprintf("faild to write message: %s", err.Error())
				c.Close()
				return
			}
		}

		<-r.Context().Done()
		c.Close()
	}

	return http.HandlerFunc(f)
}
