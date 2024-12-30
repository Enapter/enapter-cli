package enaptercli

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/urfave/cli/v2"
)

type cmdBlueprintsUpload struct {
	cmdBlueprints
	blueprintPath string
}

func buildCmdBlueprintsUpload() *cli.Command {
	cmd := &cmdBlueprintsUpload{}
	return &cli.Command{
		Name:               "upload",
		Usage:              "Upload blueprint into Platform",
		CustomHelpTemplate: cmd.HelpTemplate(),
		Flags:              cmd.Flags(),
		Before:             cmd.Before,
		Action: func(cliCtx *cli.Context) error {
			return cmd.upload(cliCtx.Context)
		},
	}
}

func (c *cmdBlueprintsUpload) Flags() []cli.Flag {
	flags := c.cmdBlueprints.Flags()
	return append(flags, &cli.StringFlag{
		Name:        "path",
		Aliases:     []string{"p"},
		Usage:       "blueprint path (zip file or directory)",
		Destination: &c.blueprintPath,
		Required:    true,
	})
}

func (c *cmdBlueprintsUpload) upload(ctx context.Context) error {
	fi, err := os.Stat(c.blueprintPath)
	if err != nil {
		return fmt.Errorf("check blueprint path: %w", err)
	}

	var data []byte
	if fi.IsDir() {
		data, err = zipDir(c.blueprintPath)
		if err != nil {
			return fmt.Errorf("zip blueprint directory: %w", err)
		}
	} else {
		data, err = os.ReadFile(c.blueprintPath)
		if err != nil {
			return fmt.Errorf("read blueprint zip file: %w", err)
		}
	}

	return c.doHTTPRequest(ctx, doHTTPRequestParams{
		Method: http.MethodPost,
		Path:   "/blueprints/upload",
		Body:   bytes.NewReader(data),
	})
}
