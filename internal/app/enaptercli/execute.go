package enaptercli

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"os"

	"github.com/urfave/cli/v2"
)

// NewApp creates a new Enapter CLI tool application instance.
func NewApp() *cli.App {
	app := cli.NewApp()

	app.Usage = "Command Line Interface (CLI) for Enapter services."
	app.Description = "The Enapter CLI requires an access token for authentication. " +
		"You can obtain the token in your Enapter Cloud account settings."
	app.CustomAppHelpTemplate = cli.AppHelpTemplate + enapterAPIEnvVarsHelp

	app.Commands = []*cli.Command{
		buildCmdSite(),
		buildCmdDevice(),
		buildCmdBlueprint(),
		buildCmdRuleEngine(),
		buildCmdConnection(),
	}

	return app
}

func zipDir(path string) ([]byte, error) {
	fsys := os.DirFS(path)

	buf := &bytes.Buffer{}
	zw := zip.NewWriter(buf)

	err := fs.WalkDir(fsys, ".", func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}

		f, err := fsys.Open(path)
		if err != nil {
			return fmt.Errorf("open: %w", err)
		}
		defer f.Close()

		zf, err := zw.Create(path)
		if err != nil {
			return fmt.Errorf("create: %w", err)
		}

		if _, err = io.Copy(zf, f); err != nil {
			return fmt.Errorf("copy: %w", err)
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
