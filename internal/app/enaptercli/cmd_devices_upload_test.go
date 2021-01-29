package enaptercli_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
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
	dirs, err := ioutil.ReadDir(testdataDir)
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
		"--api-url apiURL --blueprint-dir wrong", " ")
	app := startTestApp(args...)
	defer app.Stop()

	appErr := app.Wait()
	require.EqualError(t, appErr, `failed to zip blueprint dir "wrong": lstat wrong: no such file or directory`)
}

func testDeviceUpload(t *testing.T, dir, blueprintDir string) {
	token := faker.Word()
	hardwareID := "SIM-WTM"

	var opts deviceUploadTestSettings
	opts.Fill(t, filepath.Join(dir, "settings.json"))

	uploadReqFilename := filepath.Join(dir, "upload_req")
	uploadReq, err := ioutil.ReadFile(uploadReqFilename)
	require.NoError(t, err)

	uploadRespFilename := filepath.Join(dir, "upload_resp")
	uploadResp, err := ioutil.ReadFile(uploadRespFilename)
	require.NoError(t, err)

	resps := bytes.Split(uploadResp, []byte{'\n'})
	reqs := bytes.Split(uploadReq, []byte{'\n'})
	if update {
		reqs = nil
	}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, resp := range resps {
			if opts.CliMessage != "" {
				w.Header().Set("X-ENAPTER-CLI-MESSAGE", opts.CliMessage)
			}
			if len(resp) == 0 {
				continue
			}

			reqBody, err := ioutil.ReadAll(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte("failed to read request"))
				continue
			}

			if update {
				reqs = append(reqs, reqBody)
			} else {
				if len(reqs) == 0 {
					break
				}

				if !bytes.Equal(reqBody, reqs[0]) {
					w.WriteHeader(http.StatusBadRequest)
					_, _ = w.Write([]byte("unexpected request"))
					return
				}
				reqs = reqs[1:]
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

	expectedFileName := filepath.Join(dir, "output")
	if update {
		err := ioutil.WriteFile(expectedFileName, actual, 0600)
		require.NoError(t, err)

		err = ioutil.WriteFile(uploadReqFilename, bytes.Join(reqs, []byte{'\n'}), 0600)
		require.NoError(t, err)
	}

	expected, err := ioutil.ReadFile(expectedFileName)
	require.NoError(t, err)

	require.Equal(t, string(expected), string(actual))
}

type deviceUploadTestSettings struct {
	CliMessage string `json:"cli_message"`
}

func (s *deviceUploadTestSettings) Fill(t *testing.T, filename string) {
	t.Helper()

	data, err := ioutil.ReadFile(filename)
	if errors.Is(err, os.ErrNotExist) {
		return
	}

	require.NoError(t, json.Unmarshal(data, s))
}
