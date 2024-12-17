package enaptercli

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/urfave/cli/v2"
)

type cmdRuleEngineRuleCreate struct {
	cmdRuleEngineRule
	slug           string
	name           string
	code           string
	runtimeVersion int
	execInterval   time.Duration
	disable        bool
}

func buildCmdRuleEngineRuleCreate() *cli.Command {
	cmd := &cmdRuleEngineRuleCreate{}
	return &cli.Command{
		Name:               "create",
		Usage:              "Create a new rule",
		CustomHelpTemplate: cmd.HelpTemplate(),
		Flags:              cmd.Flags(),
		Before:             cmd.Before,
		Action: func(cliCtx *cli.Context) error {
			return cmd.do(cliCtx.Context)
		},
	}
}

func (c *cmdRuleEngineRuleCreate) Flags() []cli.Flag {
	return append(c.cmdRuleEngineRule.Flags(),
		&cli.StringFlag{
			Name:        "slug",
			Usage:       "Slug of a new rule",
			Destination: &c.slug,
			Required:    true,
		},
		&cli.StringFlag{
			Name:        "name",
			Usage:       "Name of a new rule",
			Destination: &c.name,
		},
		&cli.StringFlag{
			Name:        "code",
			Usage:       "Path to a file containing the script code",
			Destination: &c.code,
			Required:    true,
		},
		&cli.IntFlag{
			Name:        "runtime-version",
			Usage:       "Version of a runtime to use for the script execution",
			Destination: &c.runtimeVersion,
			Value:       3,
		},
		&cli.DurationFlag{
			Name:        "exec-interval",
			Usage:       "How often to execute the script. This option is only compatible with the runtime version 1",
			Destination: &c.execInterval,
		},
		&cli.BoolFlag{
			Name:        "disable",
			Usage:       "Whether to disable a rule upon creation",
			Destination: &c.disable,
		},
	)
}

func (c *cmdRuleEngineRuleCreate) do(ctx context.Context) error {
	if c.code == "-" {
		c.code = "/dev/stdin"
	}
	scriptCodeBytes, err := os.ReadFile(c.code)
	if err != nil {
		return fmt.Errorf("read script code file: %w", err)
	}
	scriptCode := base64.StdEncoding.EncodeToString(scriptCodeBytes)

	body, err := json.Marshal(map[string]any{
		"rule": map[string]any{
			"slug": c.slug,
			"name": c.name,
			"script": map[string]any{
				"code":            scriptCode,
				"runtime_version": c.runtimeVersion,
				"exec_interval":   c.execInterval.String(),
			},
		},
		"disable_rule": c.disable,
	})
	if err != nil {
		return fmt.Errorf("build request: %w", err)
	}

	return c.doHTTPRequest(ctx, doHTTPRequestParams{
		Method: http.MethodPost,
		Path:   "",
		Body:   bytes.NewReader(body),
	})
}
