package enaptercli_test

import (
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHelpMessages(t *testing.T) {
	files, err := ioutil.ReadDir("testdata/helps")
	require.NoError(t, err)

	for _, fi := range files {
		fi := fi
		t.Run(fi.Name(), func(t *testing.T) {
			args := strings.Split(fi.Name(), " ")
			args = append(args, "-h")
			app := startTestApp(args...)
			appErr := app.Wait()

			actual, err := ioutil.ReadAll(app.Stdout())
			require.NoError(t, err)

			if appErr != nil {
				actual = append(actual, []byte("app exit with error: "+appErr.Error()+"\n")...)
			}

			exepctedFileName := filepath.Join("testdata/helps", fi.Name())
			if update {
				err := ioutil.WriteFile(exepctedFileName, actual, 0600)
				require.NoError(t, err)
			}

			expected, err := ioutil.ReadFile(exepctedFileName)
			require.NoError(t, err)

			require.Equal(t, string(expected), string(actual))
		})
	}
}
