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
}

func buildCmdBlueprintsUpload() *cli.Command {
	cmd := &cmdBlueprintsUpload{}
	return &cli.Command{
		Name:               "upload",
		Usage:              "Upload blueprint directory into Platform",
		CustomHelpTemplate: cmd.HelpTemplate(),
		Flags:              cmd.Flags(),
		Before:             cmd.Before,
		Args:               true,
		ArgsUsage:          "<blueprint path (zip file or directory)>",
		Action: func(cliCtx *cli.Context) error {
			return cmd.upload(cliCtx.Context, cliCtx.Args().Get(0))
		},
	}
}

func (c *cmdBlueprintsUpload) Before(cliCtx *cli.Context) error {
	if err := c.cmdBlueprints.Before(cliCtx); err != nil {
		return err
	}

	if cliCtx.Args().Get(0) == "" {
		return errBlueprintPathMissed
	}
	return nil
}

func (c *cmdBlueprintsUpload) upload(ctx context.Context, blueprintPath string) error {
	fi, err := os.Stat(blueprintPath)
	if err != nil {
		return fmt.Errorf("check blueprint path: %w", err)
	}

	var data []byte
	if fi.IsDir() {
		data, err = zipDir(blueprintPath)
		if err != nil {
			return fmt.Errorf("zip blueprint directory: %w", err)
		}
	} else {
		data, err = os.ReadFile(blueprintPath)
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
