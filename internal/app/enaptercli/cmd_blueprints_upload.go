package enaptercli

import (
	"bytes"
	"context"
	"encoding/json"
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
		Usage:              "Upload the blueprint to the Platform",
		CustomHelpTemplate: cmd.CommandHelpTemplate(),
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
	return uploadBlueprint(ctx, c.blueprintPath, c.doHTTPRequest)
}

func uploadBlueprintAndReturnBlueprintID(ctx context.Context, blueprintPath string,
	doHTTPRequest func(context.Context, doHTTPRequestParams) error,
) (string, error) {
	var blueprintID string
	err := uploadBlueprint(ctx, blueprintPath, func(ctx context.Context, reqParams doHTTPRequestParams) error {
		reqParams.RespProcessor = func(resp *http.Response) error {
			if resp.StatusCode != http.StatusOK {
				return cli.Exit(parseRespErrorMessage(resp), 1)
			}

			var respBlueprint struct {
				Blueprint struct {
					ID string `json:"id"`
				} `json:"blueprint"`
			}
			if err := json.NewDecoder(resp.Body).Decode(&respBlueprint); err != nil {
				return fmt.Errorf("decode blueprint response: %w", err)
			}
			blueprintID = respBlueprint.Blueprint.ID
			return nil
		}
		return doHTTPRequest(ctx, reqParams)
	})
	return blueprintID, err
}

func uploadBlueprint(
	ctx context.Context, blueprintPath string,
	doHTTPRequest func(context.Context, doHTTPRequestParams) error,
) error {
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

	return doHTTPRequest(ctx, doHTTPRequestParams{
		Method: http.MethodPost,
		Path:   "/blueprints/upload",
		Body:   bytes.NewReader(data),
	})
}
