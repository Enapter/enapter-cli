package enaptercli_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bxcodec/faker/v3"
	"github.com/stretchr/testify/require"
)

func TestDeviceUploadLogs(t *testing.T) {
	errorsDir := "testdata/device_upload_logs"
	dirs, err := ioutil.ReadDir(errorsDir)
	require.NoError(t, err)

	for _, dir := range dirs {
		if !dir.IsDir() {
			continue
		}

		dir := dir
		t.Run(dir.Name(), func(t *testing.T) {
			testDeviceUploadLogs(t, filepath.Join(errorsDir, dir.Name()))
		})
	}
}

func testDeviceUploadLogs(t *testing.T, dir string) {
	testSettings := parseTestSettings(t, dir)
	srv := startTestServer(testSettings)
	defer srv.Close()

	args := strings.Split("enapter devices upload-logs", " ")
	args = append(args,
		"--token", testSettings.Token,
		"--hardware-id", testSettings.HardwareID,
		"--api-url", srv.URL)
	if testSettings.OperationID != "" {
		args = append(args, "--operation-id", testSettings.OperationID)
	}
	app := startTestApp(args...)
	defer app.Stop()

	appErr := app.Wait()

	actual, err := ioutil.ReadAll(app.Stdout())
	require.NoError(t, err)

	if appErr != nil {
		actual = append(actual, []byte("app exit with error: "+appErr.Error()+"\n")...)
	}

	expectedFileName := filepath.Join(dir, "output")
	if update {
		err := ioutil.WriteFile(expectedFileName, actual, 0600)
		require.NoError(t, err)
	}

	expected, err := ioutil.ReadFile(expectedFileName)
	require.NoError(t, err)

	require.Equal(t, string(expected), string(actual))
}

type testSettings struct {
	OperationID string   `json:"operation_id"`
	HardwareID  string   `json:"hardware_id"`
	CliMessage  string   `json:"cli_message"`
	Token       string   `json:"-"`
	Responses   [][]byte `json:"-"`
}

func parseTestSettings(t *testing.T, dir string) testSettings {
	var testSettings testSettings
	settingsBytes, err := ioutil.ReadFile(filepath.Join(dir, "settings.json"))
	require.NoError(t, err)
	require.NoError(t, json.Unmarshal(settingsBytes, &testSettings))

	uploadResp, err := ioutil.ReadFile(filepath.Join(dir, "upload_logs_resp"))
	require.NoError(t, err)

	testSettings.Responses = bytes.Split(uploadResp, []byte{'\n'})
	testSettings.Token = faker.Word()

	return testSettings
}

func startTestServer(settings testSettings) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, resp := range settings.Responses {
			if settings.CliMessage != "" {
				w.Header().Set("X-ENAPTER-CLI-MESSAGE", settings.CliMessage)
			}

			if len(resp) == 0 {
				continue
			}

			_, _ = w.Write(resp)
			settings.Responses = settings.Responses[1:]
			return
		}

		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("to much requests for test"))
	}))
}
