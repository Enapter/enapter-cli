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

	"github.com/enapter/enapter-cli/internal/app/cliflags"
)

type cmdRuleEngineRuleCreate struct {
	cmdRuleEngineRule
	slug           string
	scriptPath     string
	runtimeVersion string
	execInterval   time.Duration
	disable        bool
}

func buildCmdRuleEngineRuleCreate() *cli.Command {
	cmd := &cmdRuleEngineRuleCreate{}
	return &cli.Command{
		Name:               "create",
		Usage:              "Create a new rule",
		CustomHelpTemplate: cmd.CommandHelpTemplate(),
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
			Usage:       "Slug for the new rule",
			Destination: &c.slug,
			Required:    true,
		},
		&cli.StringFlag{
			Name:        "script",
			Usage:       "Path to the file containing the script code",
			Destination: &c.scriptPath,
			Required:    true,
		},
		&cli.StringFlag{
			Name:        "runtime-version",
			Usage:       "Version of the runtime to use for the script execution",
			Destination: &c.runtimeVersion,
			Value:       ruleRuntimeV3,
			Action: func(_ *cli.Context, v string) error {
				return c.validateRuntimeVersion(v)
			},
		},
		&cliflags.Duration{
			DurationFlag: cli.DurationFlag{
				Name:        "exec-interval",
				Usage:       "How often to execute the script (only compatible with the runtime version 1)",
				Destination: &c.execInterval,
			},
		},
		&cli.BoolFlag{
			Name:        "disable",
			Usage:       "Disable the rule upon creation",
			Destination: &c.disable,
		},
	)
}

func (c *cmdRuleEngineRuleCreate) do(ctx context.Context) error {
	if c.scriptPath == "-" {
		c.scriptPath = "/dev/stdin"
	}
	scriptBytes, err := os.ReadFile(c.scriptPath)
	if err != nil {
		return fmt.Errorf("read script code file: %w", err)
	}

	body, err := json.Marshal(map[string]any{
		"slug": c.slug,
		"script": map[string]any{
			"code":            base64.StdEncoding.EncodeToString(scriptBytes),
			"runtime_version": c.runtimeVersion,
			"exec_interval":   c.execInterval.String(),
		},
		"disable": c.disable,
	})
	if err != nil {
		return fmt.Errorf("build request: %w", err)
	}

	return c.doHTTPRequest(ctx, doHTTPRequestParams{
		Method:      http.MethodPost,
		Path:        "",
		Body:        bytes.NewReader(body),
		ContentType: contentTypeJSON,
	})
}
