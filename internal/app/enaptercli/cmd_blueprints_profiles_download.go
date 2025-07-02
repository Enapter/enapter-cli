package enaptercli

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/urfave/cli/v2"
)

type cmdBlueprintsProfilesDownload struct {
	cmdBlueprintsProfiles
	outputFileName string
}

func buildCmdBlueprintsProfilesDownload() *cli.Command {
	cmd := &cmdBlueprintsProfilesDownload{}
	return &cli.Command{
		Name:               "download",
		Usage:              "Download profiles zip from the Platform",
		CustomHelpTemplate: cmd.CommandHelpTemplate(),
		Flags:              cmd.Flags(),
		Before:             cmd.Before,
		Action: func(cliCtx *cli.Context) error {
			return cmd.do(cliCtx.Context)
		},
	}
}

func (c *cmdBlueprintsProfilesDownload) Flags() []cli.Flag {
	flags := c.cmdBlueprintsProfiles.Flags()
	return append(flags, &cli.StringFlag{
		Name:        "output",
		Aliases:     []string{"o"},
		Usage:       "file name to save the downloaded profiles",
		Destination: &c.outputFileName,
	})
}

func (c *cmdBlueprintsProfilesDownload) do(ctx context.Context) error {
	if c.outputFileName == "" {
		c.outputFileName = "profiles.zip"
	}

	return c.doHTTPRequest(ctx, doHTTPRequestParams{
		Method: http.MethodGet,
		Path:   "/blueprints/download_device_profiles",
		//nolint:bodyclose //body is closed in doHTTPRequest
		RespProcessor: okRespBodyProcessor(func(body io.Reader) error {
			outFile, err := os.Create(c.outputFileName)
			if err != nil {
				return fmt.Errorf("create output file %q: %w", c.outputFileName, err)
			}
			if _, err := io.Copy(outFile, body); err != nil {
				return fmt.Errorf("write output file %q: %w", c.outputFileName, err)
			}
			return nil
		}),
	})
}
