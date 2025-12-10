package enaptercli

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/urfave/cli/v2"
)

type cmdBlueprintProfilesUpload struct {
	cmdBlueprintProfiles
	profilesPath string
}

func buildCmdBlueprintProfilesUpload() *cli.Command {
	cmd := &cmdBlueprintProfilesUpload{}
	return &cli.Command{
		Name:               "upload",
		Usage:              "Upload profiles to the Platform",
		CustomHelpTemplate: cmd.CommandHelpTemplate(),
		Flags:              cmd.Flags(),
		Before:             cmd.Before,
		Action: func(cliCtx *cli.Context) error {
			return cmd.upload(cliCtx.Context)
		},
	}
}

func (c *cmdBlueprintProfilesUpload) Flags() []cli.Flag {
	flags := c.cmdBlueprintProfiles.Flags()
	return append(flags, &cli.StringFlag{
		Name:        "path",
		Aliases:     []string{"p"},
		Usage:       "profiles zip file path",
		Destination: &c.profilesPath,
		Required:    true,
	})
}

func (c *cmdBlueprintProfilesUpload) upload(ctx context.Context) error {
	data, err := os.ReadFile(c.profilesPath)
	if err != nil {
		return fmt.Errorf("read  zip file: %w", err)
	}

	return c.doHTTPRequest(ctx, doHTTPRequestParams{
		Method: http.MethodPost,
		Path:   "/blueprints/upload_device_profiles",
		Body:   bytes.NewReader(data),
	})
}
