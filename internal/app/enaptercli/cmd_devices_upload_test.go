package enaptercli_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/bxcodec/faker/v3"
	"github.com/stretchr/testify/require"
)

func TestDeviceUpload(t *testing.T) {
	const simpleBlueprintDir = "./testdata/device_upload/simple/blueprint"

	t.Run("simple", func(t *testing.T) {
		uploadRespFileName := "testdata/device_upload/simple/upload_resp"
		blueprintDir := "testdata/device_upload/simple/blueprint"
		expectedFileName := "testdata/device_upload/simple/output"
		testDeviceUpload(t, uploadRespFileName, blueprintDir, expectedFileName, "")
	})

	t.Run("wrong directory", func(t *testing.T) {
		uploadRespFileName := "testdata/device_upload/with_wrong_directory/upload_resp"
		blueprintDir := "testdata/device_upload/testdir"
		expectedFileName := "testdata/device_upload/with_wrong_directory/output"
		testDeviceUpload(t, uploadRespFileName, blueprintDir, expectedFileName, "")
	})

	t.Run("with dot in the paths", func(t *testing.T) {
		uploadRespFileName := "./testdata/device_upload/simple/upload_resp"
		expectedFileName := "./testdata/device_upload/simple/output"
		testDeviceUpload(t, uploadRespFileName, simpleBlueprintDir, expectedFileName, "")
	})

	t.Run("invalid hardware id", func(t *testing.T) {
		uploadRespFileName := "./testdata/device_upload/invalid_hardware_id/upload_resp"
		expectedFileName := "./testdata/device_upload/invalid_hardware_id/output"
		testDeviceUpload(t, uploadRespFileName, simpleBlueprintDir, expectedFileName, "")
	})

	t.Run("with errors in upload response", func(t *testing.T) {
		uploadRespFileName := "./testdata/device_upload/upload_errors/upload_resp"
		expectedFileName := "./testdata/device_upload/upload_errors/output"
		testDeviceUpload(t, uploadRespFileName, simpleBlueprintDir, expectedFileName, "")
	})

	t.Run("cli message", func(t *testing.T) {
		uploadRespFileName := "./testdata/device_upload/cli_message/upload_resp"
		expectedFileName := "./testdata/device_upload/cli_message/output"
		testDeviceUpload(t, uploadRespFileName, simpleBlueprintDir, expectedFileName, "VERSION IS OUTDATED\n")
	})
}

func testDeviceUpload(t *testing.T, uploadRespFilename, blueprintDir, expectedFileName, cliMessage string) {
	token := faker.Word()
	hardwareID := "SIM-WTM"

	uploadResp, err := ioutil.ReadFile(uploadRespFilename)
	require.NoError(t, err)

	resps := bytes.Split(uploadResp, []byte{'\n'})
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, resp := range resps {
			if cliMessage != "" {
				w.Header().Set("X-ENAPTER-CLI-MESSAGE", cliMessage)
			}
			if len(resp) == 0 {
				continue
			}

			_, _ = w.Write(resp)
			resps = resps[1:]
			return
		}

		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("to much requests for test"))
	}))
	defer srv.Close()

	args := strings.Split("enapter devices upload --token", " ")
	args = append(args, token, "--hardware-id", hardwareID, "--api-url", srv.URL,
		"--blueprint-dir", blueprintDir)
	app := startTestApp(args...)
	defer app.Stop()

	appErr := app.Wait()

	actual, err := ioutil.ReadAll(app.Stdout())
	require.NoError(t, err)

	if appErr != nil {
		actual = append(actual, []byte("app exit with error: "+appErr.Error()+"\n")...)
	}

	if update {
		err := ioutil.WriteFile(expectedFileName, actual, 0600)
		require.NoError(t, err)
	}

	expected, err := ioutil.ReadFile(expectedFileName)
	require.NoError(t, err)

	require.Equal(t, string(expected), string(actual))
}
