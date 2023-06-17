package enaptercli

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/urfave/cli/v2"

	"github.com/enapter/enapter-cli/internal/cloudapi"
)

const ruleUpdateDefaultTimeout = 30 * time.Second

type cmdRulesUpdate struct {
	cmdRules
	path              string
	executionInterval int
	stdlibVersion     string
	timeout           time.Duration
}

func buildCmdRulesUpdate() *cli.Command {
	cmd := &cmdRulesUpdate{}

	return &cli.Command{
		Name:               "update",
		Usage:              "Update rule.",
		CustomHelpTemplate: cmd.HelpTemplate(),
		Flags:              cmd.Flags(),
		Before:             cmd.Before,
		Action: func(cliCtx *cli.Context) error {
			return cmd.run(cliCtx.Context, cliCtx.App.Version)
		},
	}
}

func (c *cmdRulesUpdate) Flags() []cli.Flag {
	flags := c.cmdRules.Flags()
	flags = append(flags,
		&cli.StringFlag{
			Name:        "rule-path",
			Usage:       "Path to file with rule Lua code",
			Destination: &c.path,
		},
		&cli.IntFlag{
			Name:        "execution-interval",
			Usage:       "Rule execution interval in milliseconds",
			DefaultText: "chosen by the server",
			Destination: &c.executionInterval,
		},
		&cli.StringFlag{
			Name:        "stdlib-version",
			Usage:       "Version of standard library used by the rule",
			DefaultText: "chosen by the server",
			Destination: &c.stdlibVersion,
		},
		&cli.DurationFlag{
			Name:        "timeout",
			Usage:       "Time to wait for rule update",
			Destination: &c.timeout,
			Value:       ruleUpdateDefaultTimeout,
		},
	)
	return flags
}

func (c *cmdRulesUpdate) run(ctx context.Context, version string) error {
	if c.timeout != 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, c.timeout)
		defer cancel()
	}

	luaCode, err := os.ReadFile(c.path)
	if err != nil {
		return fmt.Errorf("read rule file: %w", err)
	}

	transport := cloudapi.NewCredentialsTransport(http.DefaultTransport, c.token, version)
	transport = cloudapi.NewCLIMessageWriterTransport(transport, &onceWriter{w: c.writer})
	client := cloudapi.NewClientWithURL(&http.Client{Transport: transport}, c.graphqlURL)

	input := cloudapi.UpdateRuleInput{
		RuleID:            c.ruleID,
		LuaCode:           string(luaCode),
		StdlibVersion:     c.stdlibVersion,
		ExecutionInterval: c.executionInterval,
	}

	updateData, updateErrors, err := client.UpdateRule(ctx, input)
	if err != nil {
		return fmt.Errorf("do update: %w", err)
	}

	if len(updateErrors) != 0 {
		for _, e := range updateErrors {
			fmt.Fprintln(c.writer, "[ERROR]", e.Message)
		}
		return errFinishedWithError
	}

	fmt.Fprintln(c.writer, updateData.Message)
	return nil
}
