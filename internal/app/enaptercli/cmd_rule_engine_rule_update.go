package enaptercli

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/urfave/cli/v3"

	"github.com/enapter/enapter-cli/internal/app/cliflags"
)

type cmdRuleEngineRuleUpdate struct {
	cmdRuleEngineRule
	ruleID string
	slug   string
}

func buildCmdRuleEngineRuleUpdate() *cli.Command {
	cmd := &cmdRuleEngineRuleUpdate{}
	return &cli.Command{
		Name:               "update",
		Usage:              "Update a rule",
		CustomHelpTemplate: cmd.CommandHelpTemplate(),
		Flags:              cmd.Flags(),
		Before:             cmd.Before,
		Action: func(ctx context.Context, cliCmd *cli.Command) error {
			return cmd.do(ctx, cliCmd)
		},
	}
}

func (c *cmdRuleEngineRuleUpdate) Flags() []cli.Flag {
	return append(c.cmdRuleEngineRule.Flags(),
		&cli.StringFlag{
			Name:        "rule-id",
			Usage:       "Rule ID or slug to update",
			Destination: &c.ruleID,
			Required:    true,
		},
		&cli.StringFlag{
			Name:        "slug",
			Usage:       "A new rule slug",
			Destination: &c.slug,
			Action:      cliflags.TrimSpaceAction(&c.slug),
		},
	)
}

func (c *cmdRuleEngineRuleUpdate) do(ctx context.Context, cliCmd *cli.Command) error {
	payload := make(map[string]any)

	if cliCmd.IsSet("slug") {
		payload["slug"] = c.slug
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("build request: %w", err)
	}

	return c.doHTTPRequest(ctx, doHTTPRequestParams{
		Method:      http.MethodPatch,
		Path:        "/" + c.ruleID,
		Body:        bytes.NewReader(body),
		ContentType: contentTypeJSON,
	})
}
