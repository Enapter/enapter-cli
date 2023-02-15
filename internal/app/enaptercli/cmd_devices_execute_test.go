package enaptercli_test

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/bxcodec/faker/v3"
	"github.com/stretchr/testify/require"
)

func TestDeviceExecute(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		basePath := "testdata/device_execute/simple"
		showProgress := false
		testDeviceExecute(t, basePath, showProgress, http.StatusOK)
	})

	t.Run("progress", func(t *testing.T) {
		basePath := "testdata/device_execute/progress"
		showProgress := true
		testDeviceExecute(t, basePath, showProgress, http.StatusOK)
	})

	t.Run("error", func(t *testing.T) {
		basePath := "testdata/device_execute/error"
		showProgress := false
		testDeviceExecute(t, basePath, showProgress, http.StatusForbidden)
	})
}

func testDeviceExecute(
	t *testing.T, basePath string, showProgress bool, statusCode int,
) {
	resp := readFileLines(t, filepath.Join(basePath, "responses"))
	server := startExecuteTestServer(showProgress, statusCode, resp)
	defer server.Close()

	args := []string{"enapter", "devices", "execute"}
	args = append(args,
		"--token", faker.Word(),
		"--hardware-id", faker.Word(),
		"--command", faker.Word(),
		"--api-host", server.URL)
	if showProgress {
		args = append(args, "--show-progress")
	}

	checkExecuteTestAppOutput(t, basePath, args)
}

func readFileLines(t *testing.T, path string) [][]byte {
	f, err := os.ReadFile(path)
	require.NoError(t, err)
	return bytes.Split(f, []byte{'\n'})
}

func startExecuteTestServer(
	showProgress bool, statusCode int, responses [][]byte,
) *httptest.Server {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(statusCode)

		for _, r := range responses {
			_, _ = w.Write(append(r, '\n'))
			if showProgress {
				w.(http.Flusher).Flush()
			}
		}
	}
	return httptest.NewServer(http.HandlerFunc(handler))
}

func checkExecuteTestAppOutput(t *testing.T, basePath string, args []string) {
	app := startTestApp(args...)
	defer app.Stop()

	appErr := app.Wait()

	actual, err := io.ReadAll(app.Stdout())
	require.NoError(t, err)

	if appErr != nil {
		actual = append(actual, []byte("app exit with error: "+appErr.Error()+"\n")...)
	}

	expectedFileName := filepath.Join(basePath, "output")
	if update {
		err := os.WriteFile(expectedFileName, actual, 0o600)
		require.NoError(t, err)
	}

	expected, err := os.ReadFile(expectedFileName)
	require.NoError(t, err)

	require.Equal(t, string(expected), string(actual))
}
