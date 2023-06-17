package enaptercli_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
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

func testDeviceUpload(t *testing.T, dir, blueprintDir string) {
	token := faker.Word()
	hardwareID := "SIM-WTM"

	var opts deviceUploadTestSettings
	opts.Fill(t, filepath.Join(dir, "settings.json"))

	uploadReqFilename := filepath.Join(dir, "upload_req")
	uploadReq, err := os.ReadFile(uploadReqFilename)
	require.NoError(t, err)

	uploadRespFilename := filepath.Join(dir, "upload_resp")
	uploadResp, err := os.ReadFile(uploadRespFilename)
	require.NoError(t, err)

	reqs := &sliceSliceBytes{bytes.Split(uploadReq, []byte{'\n'})}
	if update {
		reqs.buf = nil
	}
	resps := &sliceSliceBytes{bytes.Split(uploadResp, []byte{'\n'})}
	srv := startDeviceUploadTestServer(opts, reqs, resps)
	defer srv.Close()

	args := strings.Split("enapter devices upload --token", " ")
	args = append(args, token, "--hardware-id", hardwareID, "--gql-api-url", srv.URL,
		"--blueprint-dir", blueprintDir)
	app := startTestApp(args...)
	defer app.Stop()

	appErr := app.Wait()

	actual, err := io.ReadAll(app.Stdout())
	require.NoError(t, err)

	if appErr != nil {
		actual = append(actual, []byte("app exit with error: "+appErr.Error()+"\n")...)
	}

	expectedFileName := filepath.Join(dir, "output")
	if update {
		err := os.WriteFile(expectedFileName, actual, 0o600)
		require.NoError(t, err)

		err = os.WriteFile(uploadReqFilename, bytes.Join(reqs.buf, []byte{'\n'}), 0o600)
		require.NoError(t, err)
	}

	expected, err := os.ReadFile(expectedFileName)
	require.NoError(t, err)

	require.Equal(t, string(expected), string(actual))
}

type deviceUploadTestSettings struct {
	CliMessage string `json:"cli_message"`
}

func (s *deviceUploadTestSettings) Fill(t *testing.T, filename string) {
	t.Helper()

	data, err := os.ReadFile(filename)
	if errors.Is(err, os.ErrNotExist) {
		return
	}

	require.NoError(t, json.Unmarshal(data, s))
}

func startDeviceUploadTestServer(
	opts deviceUploadTestSettings, reqs, resps *sliceSliceBytes,
) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := resps.Next()
		if len(resp) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte("to much requests for test (not enough responses)"))
			return
		}

		var req []byte
		if !update {
			req = reqs.Next()
			if len(req) == 0 {
				w.WriteHeader(http.StatusBadRequest)
				_, _ = w.Write([]byte("to much requests for test (not enough requests)"))
				return
			}
		}

		if opts.CliMessage != "" {
			w.Header().Set("X-ENAPTER-CLI-MESSAGE", opts.CliMessage)
		}

		reqBody, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte("failed to read request"))
			return
		}

		if update {
			reqs.Append(reqBody)
		} else {
			reqBody := bytes.TrimRight(reqBody, "\n")
			if !bytes.Equal(reqBody, req) {
				w.WriteHeader(http.StatusBadRequest)
				_, _ = w.Write([]byte("unexpected request\nActual\n"))
				_, _ = w.Write(reqBody)
				_, _ = w.Write([]byte("\nExpected\n"))
				_, _ = w.Write(req)
				return
			}
		}

		_, _ = w.Write(resp)
	}))
}

type sliceSliceBytes struct {
	buf [][]byte
}

func (b *sliceSliceBytes) Next() []byte {
	for i, s := range b.buf {
		if len(s) != 0 {
			b.buf = b.buf[i+1:]
			return s
		}
	}
	return nil
}

func (b *sliceSliceBytes) Append(d []byte) {
	b.buf = append(b.buf, d)
}
