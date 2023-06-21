package enaptercli_test

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

type byteSliceSlice struct {
	lines [][]byte
}

func byteSliceSliceFromFile(t *testing.T, path string) *byteSliceSlice {
	f, err := os.ReadFile(path)
	require.NoError(t, err)
	lines := bytes.Split(f, []byte{'\n'})

	n := 0
	for _, line := range lines {
		if len(line) > 0 {
			lines[n] = line
			n++
		}
	}

	return &byteSliceSlice{lines: lines[:n]}
}

func (b *byteSliceSlice) Next() []byte {
	for i, s := range b.lines {
		b.lines = b.lines[i+1:]
		return s
	}
	return nil
}

func (b *byteSliceSlice) Append(d []byte) {
	b.lines = append(b.lines, d)
}

func (b *byteSliceSlice) Buffer() [][]byte {
	return b.lines
}

func (b *byteSliceSlice) Clear() {
	b.lines = nil
}

func startTestServer(reqs, resps *byteSliceSlice, cliMessage string) *httptest.Server {
	if update {
		reqs.Clear()
	}

	handler := func(w http.ResponseWriter, r *http.Request) {
		resp := resps.Next()
		if resp == nil {
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

		if cliMessage != "" {
			w.Header().Set("X-ENAPTER-CLI-MESSAGE", cliMessage)
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
	}

	return httptest.NewServer(http.HandlerFunc(handler))
}

func checkTestAppOutput(t *testing.T, basePath string, args []string, requests *byteSliceSlice) {
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

		requestsFileName := filepath.Join(basePath, "requests")
		requestsBytes := bytes.Join(requests.Buffer(), []byte{'\n'})
		err = os.WriteFile(requestsFileName, requestsBytes, 0o600)
		require.NoError(t, err)
	}

	expected, err := os.ReadFile(expectedFileName)
	require.NoError(t, err)

	require.Equal(t, string(expected), string(actual))
}
