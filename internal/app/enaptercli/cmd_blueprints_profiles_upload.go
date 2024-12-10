package enaptercli

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/urfave/cli/v2"
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
		Args:               true,
		ArgsUsage:          "profiles zip file path",
		Action: func(cliCtx *cli.Context) error {
			return cmd.upload(cliCtx.Context, cliCtx.Args().Get(0))
		},
	}
}

func (c *cmdBlueprintsProfilesUpload) Before(cliCtx *cli.Context) error {
	if err := c.cmdBlueprintsProfiles.Before(cliCtx); err != nil {
		return err
	}

	if cliCtx.Args().Get(0) == "" {
		return errProfilesPathMissed
	}
	return nil
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
