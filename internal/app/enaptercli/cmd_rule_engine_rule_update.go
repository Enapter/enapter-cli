package enaptercli

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/urfave/cli/v3"
)

type cmdRuleEngineRuleUpdate struct {
	cmdRuleEngineRule
	slug string
	name string
}

func buildCmdRuleEngineRuleUpdate() *cli.Command {
	cmd := &cmdRuleEngineRuleUpdate{}
	return &cli.Command{
		Name:               "update",
		Usage:              "Update a rule",
		ArgsUsage:          "RULE",
		CustomHelpTemplate: cmd.HelpTemplate(),
		Flags:              cmd.Flags(),
		Before:             cmd.Before,
		Action: func(ctx context.Context, cm *cli.Command) error {
			return cmd.do(ctx, cm, cm.Args().First())
		},
	}
}

func (c *cmdRuleEngineRuleUpdate) Before(
	ctx context.Context, cm *cli.Command,
) (context.Context, error) {
	ctx, err := c.cmdRuleEngineRule.Before(ctx, cm)
	if err != nil {
		return nil, err
	}

	if cm.Args().Get(0) == "" {
		return nil, errRequiresAtLeastOneArgument
	}

	return ctx, nil
}

func (c *cmdRuleEngineRuleUpdate) Flags() []cli.Flag {
	return append(c.cmdRuleEngineRule.Flags(),
		&cli.StringFlag{
			Name:        "slug",
			Usage:       "A new rule slug",
			Destination: &c.slug,
		},
		&cli.StringFlag{
			Name:        "name",
			Usage:       "A new rule name",
			Destination: &c.name,
		},
	)
}

func (c *cmdRuleEngineRuleUpdate) do(
	ctx context.Context, cm *cli.Command, rule string,
) error {
	payload := struct {
		Rule       map[string]any `json:"rule"`
		UpdateMask string         `json:"update_mask"`
	}{
		Rule:       make(map[string]any),
		UpdateMask: "",
	}

	if cm.IsSet("slug") {
		payload.Rule["slug"] = c.slug
		payload.UpdateMask = payload.UpdateMask + "slug,"
	}
	if cm.IsSet("name") {
		payload.Rule["name"] = c.name
		payload.UpdateMask = payload.UpdateMask + "name,"
	}

	payload.UpdateMask = strings.TrimSuffix(payload.UpdateMask, ",")

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("build request: %w", err)
	}

	return c.doHTTPRequest(ctx, doHTTPRequestParams{
		Method: http.MethodPatch,
		Path:   "/" + rule,
		Body:   bytes.NewReader(body),
	})
}
