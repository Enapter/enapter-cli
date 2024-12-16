package enaptercli

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/urfave/cli/v2"
)

type cmdBlueprintsDownload struct {
	cmdBlueprints
	blueprintID    string
	outputFileName string
}

func buildCmdBlueprintsDownload() *cli.Command {
	cmd := &cmdBlueprintsDownload{}
	return &cli.Command{
		Name:               "download",
		Usage:              "Download blueprint zip from Platform",
		CustomHelpTemplate: cmd.HelpTemplate(),
		Flags:              cmd.Flags(),
		Before:             cmd.Before,
		Action: func(cliCtx *cli.Context) error {
			return cmd.do(cliCtx.Context)
		},
	}
}

func (c *cmdBlueprintsDownload) Flags() []cli.Flag {
	flags := c.cmdBlueprints.Flags()
	return append(flags, &cli.StringFlag{
		Name:        "blueprint-id",
		Aliases:     []string{"b"},
		Usage:       "blueprint name or ID to download",
		Destination: &c.blueprintID,
		Required:    true,
	}, &cli.StringFlag{
		Name:        "output",
		Aliases:     []string{"o"},
		Usage:       "blueprint file name to save",
		Destination: &c.outputFileName,
	})
}

func (c *cmdBlueprintsDownload) do(ctx context.Context) error {
	if c.outputFileName == "" {
		c.outputFileName = strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(c.blueprintID,
			":", "_"), ".", "_"), "/", "_") + ".enbp"
	}

	if !isBlueprintID(c.blueprintID) {
		blueprintName, blueprintTag := parseBlueprintName(c.blueprintID)
		err := c.doHTTPRequest(ctx, doHTTPRequestParams{
			Method: http.MethodGet,
			Path:   "/blueprints/enapter/" + blueprintName + "/" + blueprintTag,
			//nolint:bodyclose //body is closed in doHTTPRequest
			RespProcessor: okRespBodyProcessor(func(body io.Reader) error {
				var resp struct {
					Blueprint struct {
						ID string `json:"id"`
					} `json:"blueprint"`
				}
				if err := json.NewDecoder(body).Decode(&resp); err != nil {
					return fmt.Errorf("parse response body: %w", err)
				}
				c.blueprintID = resp.Blueprint.ID
				return nil
			}),
		})
		if err != nil {
			return fmt.Errorf("get blueprint info by name: %w", err)
		}
	}

	return c.doHTTPRequest(ctx, doHTTPRequestParams{
		Method: http.MethodGet,
		Path:   "/blueprints/" + c.blueprintID + "/zip",
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
