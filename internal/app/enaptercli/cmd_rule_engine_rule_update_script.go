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
	code           string
	runtimeVersion int
	execInterval   time.Duration
}

func buildCmdRuleEngineRuleUpdateScript() *cli.Command {
	cmd := &cmdRuleEngineRuleUpdateScript{}
	return &cli.Command{
		Name:               "update-script",
		Usage:              "Update the script of a rule",
		Args:               true,
		ArgsUsage:          "RULE",
		CustomHelpTemplate: cmd.HelpTemplate(),
		Flags:              cmd.Flags(),
		Before:             cmd.Before,
		Action: func(cliCtx *cli.Context) error {
			return cmd.do(cliCtx.Context, cliCtx.Args().First())
		},
	}
}

func (c *cmdRuleEngineRuleUpdateScript) Before(cliCtx *cli.Context) error {
	if err := c.cmdRuleEngineRule.Before(cliCtx); err != nil {
		return err
	}

	if cliCtx.Args().Get(0) == "" {
		return errRequiresAtLeastOneArgument
	}

	return nil
}

func (c *cmdRuleEngineRuleUpdateScript) Flags() []cli.Flag {
	return append(c.cmdRuleEngineRule.Flags(),
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
		},
		&cli.DurationFlag{
			Name:        "exec-interval",
			Usage:       "How often to execute the script. This option is only compatible with the runtime version 1",
			Destination: &c.execInterval,
		},
	)
}

func (c *cmdRuleEngineRuleUpdateScript) do(ctx context.Context, rule string) error {
	if c.code == "-" {
		c.code = "/dev/stdin"
	}
	codeBytes, err := os.ReadFile(c.code)
	if err != nil {
		return fmt.Errorf("read script code file: %w", err)
	}
	code := base64.StdEncoding.EncodeToString(codeBytes)

	body, err := json.Marshal(map[string]any{
		"new_rule_script": map[string]any{
			"code":            code,
			"runtime_version": c.runtimeVersion,
			"exec_interval":   c.execInterval.String(),
		},
	})
	if err != nil {
		return fmt.Errorf("build request: %w", err)
	}

	return c.doHTTPRequest(ctx, doHTTPRequestParams{
		Method: http.MethodPost,
		Path:   "/" + rule + "/update_script",
		Body:   bytes.NewReader(body),
	})
}
