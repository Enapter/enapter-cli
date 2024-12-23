package enaptercli

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/urfave/cli/v3"
)

type cmdBlueprintsProfilesUpload struct {
	cmdBlueprintsProfiles
}

func buildCmdBlueprintsProfilesUpload() *cli.Command {
	cmd := &cmdBlueprintsProfilesUpload{}
	return &cli.Command{
		Name:               "upload",
		Usage:              "Upload profiles into Platform",
		CustomHelpTemplate: cmd.HelpTemplate(),
		Flags:              cmd.Flags(),
		Before:             cmd.Before,
		ArgsUsage:          "profiles zip file path",
		Action: func(ctx context.Context, cm *cli.Command) error {
			return cmd.upload(ctx, cm.Args().Get(0))
		},
	}
}

func (c *cmdBlueprintsProfilesUpload) Before(
	ctx context.Context, cm *cli.Command,
) (context.Context, error) {
	ctx, err := c.cmdBlueprintsProfiles.Before(ctx, cm)
	if err != nil {
		return nil, err
	}

	if cm.Args().Get(0) == "" {
		return nil, errProfilesPathMissed
	}
	return ctx, nil
}

func (c *cmdBlueprintsProfilesUpload) upload(ctx context.Context, blueprintPath string) error {
	data, err := os.ReadFile(blueprintPath)
	if err != nil {
		return fmt.Errorf("read  zip file: %w", err)
	}

	return c.doHTTPRequest(ctx, doHTTPRequestParams{
		Method: http.MethodPost,
		Path:   "/blueprints/upload_device_profiles",
		Body:   bytes.NewReader(data),
	})
}
