package enaptercli_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"text/template"

	"github.com/stretchr/testify/require"
)

const testToken = "enapter_api_test_token"

func TestHelpMessages(t *testing.T) {
	files, err := os.ReadDir("testdata/helps")
	require.NoError(t, err)

	for _, fi := range files {
		fi := fi
		t.Run(fi.Name(), func(t *testing.T) {
			args := strings.Split(fi.Name(), " ")
			args = append(args, "-h")
			app := startTestApp(args...)
			appErr := app.Wait()

			actual, err := io.ReadAll(app.Stdout())
			require.NoError(t, err)

			if appErr != nil {
				actual = append(actual, []byte("app exit with error: "+appErr.Error()+"\n")...)
			}

			exepctedFileName := filepath.Join("testdata/helps", fi.Name())
			if update {
				err := os.WriteFile(exepctedFileName, actual, 0o600)
				require.NoError(t, err)
			} else {
				require.Equal(t, readFileToString(t, exepctedFileName), string(actual))
			}
		})
	}
}

func TestHTTPReqResp(t *testing.T) {
	const testdataPath = "testdata/http_req_resp"
	tests, err := os.ReadDir(testdataPath)
	require.NoError(t, err)

	for _, tc := range tests {
		tc := tc
		t.Run(tc.Name(), func(t *testing.T) {
			reqCount := 0
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				reqObj := struct {
					Method string
					URL    string
					Header http.Header
					Body   string
				}{
					r.Method,
					r.URL.String(),
					r.Header,
					readBodyAsString(t, r.Body),
				}
				expReqFileName := filepath.Join(testdataPath, tc.Name(), "req_"+strconv.Itoa(reqCount))
				if update {
					err := os.WriteFile(expReqFileName, shouldMarshalIndent(t, reqObj), 0o600)
					require.NoError(t, err)
				}

				resp := shouldReadFile(t, filepath.Join(testdataPath, tc.Name(), "resp_"+strconv.Itoa(reqCount)))
				_, _ = w.Write(resp)
			}))
			defer srv.Close()

			tmplParams := struct {
				BaseFlags string
			}{
				BaseFlags: strings.Join([]string{"--token", testToken, "--api-host", srv.URL}, " "),
			}

			cmd := executeTmpl(t, filepath.Join(testdataPath, tc.Name(), "cmd.tmpl"), tmplParams)
			args := strings.Split(cmd, " ")
			app := startTestApp(args...)
			appErr := app.Wait()

			actual, err := io.ReadAll(app.Stdout())
			require.NoError(t, err)

			if appErr != nil {
				actual = append(actual, []byte("app exit with error: "+appErr.Error()+"\n")...)
			}

			exepctedOutFileName := filepath.Join(testdataPath, tc.Name(), "out")
			if update {
				err := os.WriteFile(exepctedOutFileName, actual, 0o600)
				require.NoError(t, err)
			} else {
				require.Equal(t, readFileToString(t, exepctedOutFileName), string(actual))
			}
		})
	}
}

func executeTmpl(t *testing.T, tmplFilePath string, tmplParams interface{}) string {
	t.Helper()
	tmplData := readFileToString(t, tmplFilePath)
	tmplData = strings.TrimRight(tmplData, " \n\t")

	tmpl := template.New(tmplFilePath)
	tmpl, err := tmpl.Parse(tmplData)
	require.NoError(t, err)

	out := &bytes.Buffer{}
	require.NoError(t, tmpl.Execute(out, tmplParams))

	return out.String()
}

func readFileToString(t *testing.T, path string) string {
	t.Helper()
	return string(shouldReadFile(t, path))
}

func shouldReadFile(t *testing.T, path string) []byte {
	t.Helper()
	d, err := os.ReadFile(path)
	require.NoError(t, err)
	return d
}

func readBodyAsString(t *testing.T, r io.Reader) string {
	t.Helper()
	d, err := io.ReadAll(r)
	require.NoError(t, err)
	return string(d)
}

func shouldMarshalIndent(t *testing.T, v interface{}) []byte {
	t.Helper()
	d, err := json.MarshalIndent(v, "", "  ")
	require.NoError(t, err)
	return d
}
