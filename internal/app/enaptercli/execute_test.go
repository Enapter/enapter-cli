package enaptercli_test

import (
	"bytes"
	"encoding/base64"
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
	"unicode/utf8"

	"github.com/stretchr/testify/require"
)

const testToken = "enapter_api_test_token"

func TestHelpMessages(t *testing.T) {
	files, err := os.ReadDir("testdata/helps")
	require.NoError(t, err)

	for _, fi := range files {
		t.Run(fi.Name(), func(t *testing.T) {
			args := strings.Split(fi.Name(), " ")
			args = append(args, "-h")
			app := startTestApp(args...)
			appErr := app.Wait()

			actual, err := io.ReadAll(app.Output())
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
		t.Run(tc.Name(), func(t *testing.T) {
			path := filepath.Join(testdataPath, tc.Name())
			testExecute(t, path)
		})
	}
}

func testExecute(t *testing.T, path string) {
	srv := newTestServer(t, path)

	cmd := executeTmpl(t, filepath.Join(path, "cmd.tmpl"), struct {
		Token string
		URL   string
	}{
		Token: testToken,
		URL:   srv.URL,
	})

	t.Setenv("ENAPTER3_CONFIG", t.TempDir())
	output := executeCommands(t, cmd)

	exepctedOutFileName := filepath.Join(path, "out")
	if update {
		err := os.WriteFile(exepctedOutFileName, output, 0o600)
		require.NoError(t, err)
	} else {
		expected := readFileToString(t, exepctedOutFileName)
		require.Equal(t, expected, string(output))
	}
}

func newTestServer(t *testing.T, path string) *httptest.Server {
	t.Helper()

	reqCount := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqPath := filepath.Join(path, "req_"+strconv.Itoa(reqCount))
		respPath := filepath.Join(path, "resp_"+strconv.Itoa(reqCount))

		info := struct {
			Method string
			URL    string
			Header http.Header
			Body   string
		}{
			Method: r.Method,
			URL:    r.URL.String(),
			Header: r.Header,
			Body:   readBodyAsString(t, r.Body),
		}
		if update {
			err := os.WriteFile(reqPath, shouldMarshalIndent(t, info), 0o600)
			require.NoError(t, err)
		} else {
			expected := readFileToString(t, reqPath)
			actual := string(shouldMarshalIndent(t, info))
			require.Equal(t, expected, actual)
		}

		resp := shouldReadFile(t, respPath)
		_, _ = w.Write(resp)

		reqCount++
	}))
	t.Cleanup(func() { srv.Close() })

	return srv
}

func executeCommands(t *testing.T, cmd string) []byte {
	t.Helper()

	var output []byte
	for cmd := range strings.Lines(cmd) {
		cmd := strings.Trim(cmd, "\n")
		args := strings.Split(cmd, " ")

		app := startTestApp(args...)
		appErr := app.Wait()

		out, err := io.ReadAll(app.Output())
		require.NoError(t, err)

		output = append(output, out...)
		if appErr != nil {
			output = append(output, []byte("app exit with error: "+appErr.Error()+"\n")...)
			break
		}
	}
	return output
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
	if !utf8.Valid(d) {
		return "base64:" + base64.StdEncoding.EncodeToString(d)
	}
	return string(d)
}

func shouldMarshalIndent(t *testing.T, v interface{}) []byte {
	t.Helper()
	d, err := json.MarshalIndent(v, "", "  ")
	require.NoError(t, err)
	return d
}
