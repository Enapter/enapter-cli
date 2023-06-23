package enaptercli

import (
	"archive/zip"
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/urfave/cli/v2"
)

// NewApp creates a new Enapter CLI tool application instance.
func NewApp() *cli.App {
	app := cli.NewApp()

	app.Usage = "Command line interface for Enapter services."
	app.Description = "Enapter CLI requires access token for authentication. " +
		"The token can be obtained in your Enapter Cloud account settings.\n\n" +
		"Configure API token using ENAPTER_API_TOKEN environment variable or using --token global option."

	app.Commands = []*cli.Command{
		buildCmdDevices(),
		buildCmdRules(),
	}

	return app
}

func zipDir(path string) ([]byte, error) {
	buf := &bytes.Buffer{}
	myZip := zip.NewWriter(buf)

	path = strings.TrimPrefix(path, "./")

	err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		relPath := strings.TrimPrefix(filePath, path)
		relPath = strings.TrimPrefix(relPath, "/")
		zipFile, err := myZip.Create(relPath)
		if err != nil {
			return err
		}
		fsFile, err := os.Open(filePath)
		if err != nil {
			return err
		}
		_, err = io.Copy(zipFile, fsFile)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	if err := myZip.Close(); err != nil {
		return nil, err
	}

	return buf.Bytes(), err
}
