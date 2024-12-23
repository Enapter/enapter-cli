package enaptercli

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/urfave/cli/v3"
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
		ArgsUsage:          "<blueprint path (zip file or directory)>",
		Action: func(ctx context.Context, cm *cli.Command) error {
			return cmd.upload(ctx, cm.Args().Get(0))
		},
	}
}

func (c *cmdBlueprintsUpload) Before(
	ctx context.Context, cm *cli.Command,
) (context.Context, error) {
	ctx, err := c.cmdBlueprints.Before(ctx, cm)
	if err != nil {
		return nil, err
	}

	if cm.Args().Get(0) == "" {
		return nil, errBlueprintPathMissed
	}
	return ctx, nil
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
