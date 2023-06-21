package enaptercli_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/bxcodec/faker/v3"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/require"
)

func testLogsCommand(
	t *testing.T, inputFileName, untilLinePrefix, expectedFileName string,
	identifier map[string]string, args []string,
) {
	token := faker.Word()
	messagesBytes, err := os.ReadFile(inputFileName)
	require.NoError(t, err)
	messages := bytes.Split(messagesBytes, []byte{'\n'})

	handleErrCh := make(chan string)
	wsPath, srv := startTestLogsServer(t, token, identifier, messages, handleErrCh)
	defer srv.Close()

	args = append(args, "--token", token, "--ws-api-url", wsPath)
	app := startTestApp(args...)
	defer app.Stop()

	actual := readOutputUntilLineOrError(t, app.Stdout(), untilLinePrefix, handleErrCh)
	if update {
		err := os.WriteFile(expectedFileName, []byte(actual), 0o600)
		require.NoError(t, err)
	}

	expected, err := os.ReadFile(expectedFileName)
	require.NoError(t, err)
	require.Equal(t, string(expected), actual)

	app.Stop()
	appErr := app.Wait()
	require.NoError(t, appErr)

	restOutput, err := io.ReadAll(app.Stdout())
	require.NoError(t, err)
	require.Empty(t, string(restOutput))
}

func startTestLogsServer(
	t *testing.T, token string, identifier map[string]string, messages [][]byte,
	handleErrCh chan<- string,
) (string, *httptest.Server) {
	t.Helper()

	handler := buildTestLogsHandler(token, identifier, messages, handleErrCh)
	srv := httptest.NewServer(handler)

	u, err := url.Parse(srv.URL)
	require.NoError(t, err)
	u.Scheme = "ws"

	return u.String(), srv
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

//nolint:funlen // because contains a lot of simple logged checks.
func buildTestLogsHandler(
	token string, identifier map[string]string, messages [][]byte, handleErrCh chan<- string,
) http.Handler {
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
		defer c.Close()

		msgType, msgBytes, err := c.ReadMessage()
		if err != nil {
			handleErrCh <- fmt.Sprintf("failed to read subscribe message: %s", err)
			return
		}

		if msgType != websocket.TextMessage {
			handleErrCh <- fmt.Sprintf("subscribe message should be text type [%d], but [%d]",
				websocket.TextMessage, msgType)
			return
		}

		subMsg := struct {
			Command    string `json:"command"`
			Identifier string `json:"identifier"`
		}{}
		if err := json.Unmarshal(msgBytes, &subMsg); err != nil {
			handleErrCh <- fmt.Sprintf("failed to unmarshall subsribe message %q: %s", string(msgBytes), err.Error())
			return
		}

		if subMsg.Command != "subscribe" {
			handleErrCh <- fmt.Sprintf("this is not subscribe message, but %q", subMsg.Command)
			return
		}

		var reqIdentifier map[string]string
		if err := json.Unmarshal([]byte(subMsg.Identifier), &reqIdentifier); err != nil {
			handleErrCh <- fmt.Sprintf("failed to unmarshall subsribe message identifier %q: %s",
				subMsg.Identifier, err.Error())
			return
		}

		if !reflect.DeepEqual(identifier, reqIdentifier) {
			handleErrCh <- fmt.Sprintf("subsribe message identifier shoud be equal to %q, but %q",
				identifier, reqIdentifier)
			return
		}

		for _, m := range messages {
			if len(m) == 0 {
				continue
			}
			if err := c.WriteMessage(websocket.TextMessage, m); err != nil {
				handleErrCh <- fmt.Sprintf("failed to write message: %s", err.Error())
				return
			}
		}

		<-r.Context().Done()
	}

	return http.HandlerFunc(f)
}
