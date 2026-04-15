package enaptercli

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/urfave/cli/v3"
)

type cmdBlueprintProfilesDownload struct {
	cmdBlueprintProfiles
	outputFileName string
}

func buildCmdBlueprintProfilesDownload() *cli.Command {
	cmd := &cmdBlueprintProfilesDownload{}
	return &cli.Command{
		Name:               "download",
		Usage:              "Download profiles zip from the Platform",
		CustomHelpTemplate: cmd.CommandHelpTemplate(),
		Flags:              cmd.Flags(),
		Before:             cmd.Before,
		Action: func(ctx context.Context, _ *cli.Command) error {
			return cmd.do(ctx)
		},
	}
}

func (c *cmdBlueprintProfilesDownload) Flags() []cli.Flag {
	flags := c.cmdBlueprintProfiles.Flags()
	return append(flags, &cli.StringFlag{
		Name:        "output",
		Aliases:     []string{"o"},
		Usage:       "File name to save the downloaded profiles",
		Destination: &c.outputFileName,
	})
}

func (c *cmdBlueprintProfilesDownload) do(ctx context.Context) error {
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
			fmt.Fprintln(c.writer, c.outputFileName)
			return nil
		}),
	})
}
