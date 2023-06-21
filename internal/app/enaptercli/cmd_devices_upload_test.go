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

const blueprintDir = "testdata/device_upload/simple/blueprint"

func TestDeviceUpload(t *testing.T) {
	testdataDir := "testdata/device_upload"
	dirs, err := os.ReadDir(testdataDir)
	require.NoError(t, err)

	for _, dir := range dirs {
		if !dir.IsDir() {
			continue
		}

		dir := dir
		t.Run(dir.Name(), func(t *testing.T) {
			testDeviceUpload(t, filepath.Join(testdataDir, dir.Name()), blueprintDir)
		})
	}
}

func TestDeviceUploadBlueprintDirWithDot(t *testing.T) {
	testDeviceUpload(t, "testdata/device_upload/simple", "./"+blueprintDir)
}

func TestDeviceUploadWrongBlueprintDir(t *testing.T) {
	args := strings.Split("enapter devices upload --token token --hardware-id hardwareID "+
		"--gql-api-url apiURL --blueprint-dir wrong", " ")
	app := startTestApp(args...)
	defer app.Stop()

	appErr := app.Wait()
	require.EqualError(t, appErr, `failed to zip blueprint dir "wrong": lstat wrong: no such file or directory`)
}

type devicesUploadTestSettings struct {
	HardwareID string `json:"hardware_id"`
	CliMessage string `json:"cli_message"`
	Token      string `json:"-"`
}

func (s *devicesUploadTestSettings) Fill(t *testing.T, dir string) {
	settingsBytes, err := os.ReadFile(filepath.Join(dir, "settings.json"))
	require.NoError(t, err)
	require.NoError(t, json.Unmarshal(settingsBytes, s))
	s.Token = faker.Word()
}

func testDeviceUpload(t *testing.T, dir, blueprintDir string) {
	var settings devicesUploadTestSettings
	settings.Fill(t, dir)

	reqs := byteSliceSliceFromFile(t, filepath.Join(dir, "requests"))
	resps := byteSliceSliceFromFile(t, filepath.Join(dir, "responses"))

	srv := startTestServer(reqs, resps, settings.CliMessage)
	defer srv.Close()

	args := strings.Split("enapter devices upload", " ")
	args = append(args,
		"--token", settings.Token,
		"--hardware-id", settings.HardwareID,
		"--blueprint-dir", blueprintDir,
		"--gql-api-url", srv.URL)

	checkTestAppOutput(t, dir, args, reqs)
}
