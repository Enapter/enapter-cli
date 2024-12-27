package enaptercli

import (
	"archive/zip"
	"bytes"
	"fmt"
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
		"The token can be obtained in your Enapter Cloud account settings."

	app.Commands = []*cli.Command{
		buildCmdDevices(),
		buildCmdBlueprints(),
		buildCmdProvisioning(),
		buildCmdRuleEngine(),
	}

	return app
}

func zipDir(path string) ([]byte, error) {
	buf := &bytes.Buffer{}
	zw := zip.NewWriter(buf)

	path = filepath.Clean(path)
	err := filepath.WalkDir(path, func(filePath string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}

		relPath := strings.TrimPrefix(filePath, path)
		relPath = strings.TrimPrefix(relPath, "/")
		zipFile, err := zw.Create(relPath)
		if err != nil {
			return err
		}

		fsFile, err := os.Open(filePath)
		if err != nil {
			return err
		}
		defer fsFile.Close()

		_, err = io.Copy(zipFile, fsFile)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("walk dir: %w", err)
	}

	if err := zw.Close(); err != nil {
		return nil, fmt.Errorf("close zip: %w", err)
	}

	return buf.Bytes(), nil
}
