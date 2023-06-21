package enaptercli_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bxcodec/faker/v3"
	"github.com/stretchr/testify/require"
)

func TestDeviceUploadLogs(t *testing.T) {
	errorsDir := "testdata/device_upload_logs"
	dirs, err := os.ReadDir(errorsDir)
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

type devicesUploadLogsTestSettings struct {
	OperationID string `json:"operation_id"`
	HardwareID  string `json:"hardware_id"`
	CliMessage  string `json:"cli_message"`
	Token       string `json:"-"`
}

func (s *devicesUploadLogsTestSettings) Fill(t *testing.T, dir string) {
	settingsBytes, err := os.ReadFile(filepath.Join(dir, "settings.json"))
	require.NoError(t, err)
	require.NoError(t, json.Unmarshal(settingsBytes, s))
	s.Token = faker.Word()
}

func testDeviceUploadLogs(t *testing.T, dir string) {
	var settings devicesUploadLogsTestSettings
	settings.Fill(t, dir)

	reqs := byteSliceSliceFromFile(t, filepath.Join(dir, "requests"))
	resps := byteSliceSliceFromFile(t, filepath.Join(dir, "responses"))

	srv := startTestServer(reqs, resps, settings.CliMessage)
	defer srv.Close()

	args := strings.Split("enapter devices upload-logs", " ")
	args = append(args,
		"--token", settings.Token,
		"--hardware-id", settings.HardwareID,
		"--gql-api-url", srv.URL)
	if settings.OperationID != "" {
		args = append(args, "--operation-id", settings.OperationID)
	}

	checkTestAppOutput(t, dir, args, reqs)
}
