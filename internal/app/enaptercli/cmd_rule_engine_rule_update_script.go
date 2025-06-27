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

type cmdRuleEngineRuleUpdateScript struct {
	cmdRuleEngineRule
	ruleID         string
	scriptPath     string
	runtimeVersion string
	execInterval   time.Duration
}

func buildCmdRuleEngineRuleUpdateScript() *cli.Command {
	cmd := &cmdRuleEngineRuleUpdateScript{}
	return &cli.Command{
		Name:               "update-script",
		Usage:              "Update the script of a rule",
		CustomHelpTemplate: cmd.CommandHelpTemplate(),
		Flags:              cmd.Flags(),
		Before:             cmd.Before,
		Action: func(cliCtx *cli.Context) error {
			return cmd.do(cliCtx.Context)
		},
	}
}

func (c *cmdRuleEngineRuleUpdateScript) Flags() []cli.Flag {
	return append(c.cmdRuleEngineRule.Flags(),
		&cli.StringFlag{
			Name:        "rule-id",
			Usage:       "Rule ID or slug to update",
			Destination: &c.ruleID,
			Required:    true,
		},
		&cli.StringFlag{
			Name:        "script",
			Usage:       "Path to a file containing the script code",
			Destination: &c.scriptPath,
			Required:    true,
		},
		&cli.StringFlag{
			Name:        "runtime-version",
			Usage:       "Version of a runtime to use for the script execution",
			Destination: &c.runtimeVersion,
			Value:       ruleRuntimeV3,
			Action: func(_ *cli.Context, v string) error {
				return c.validateRuntimeVersion(v)
			},
		},
		&cli.DurationFlag{
			Name:        "exec-interval",
			Usage:       "How often to execute the script. This option is only compatible with the runtime version 1",
			Destination: &c.execInterval,
		},
	)
}

func (c *cmdRuleEngineRuleUpdateScript) do(ctx context.Context) error {
	if c.scriptPath == "-" {
		c.scriptPath = "/dev/stdin"
	}
	scriptBytes, err := os.ReadFile(c.scriptPath)
	if err != nil {
		return fmt.Errorf("read script file: %w", err)
	}

	body, err := json.Marshal(map[string]any{
		"new_rule_script": map[string]any{
			"code":            base64.StdEncoding.EncodeToString(scriptBytes),
			"runtime_version": c.runtimeVersion,
			"exec_interval":   c.execInterval.String(),
		},
	})
	if err != nil {
		return fmt.Errorf("build request: %w", err)
	}

	return c.doHTTPRequest(ctx, doHTTPRequestParams{
		Method:      http.MethodPost,
		Path:        "/" + c.ruleID + "/update_script",
		Body:        bytes.NewReader(body),
		ContentType: contentTypeJSON,
	})
}
