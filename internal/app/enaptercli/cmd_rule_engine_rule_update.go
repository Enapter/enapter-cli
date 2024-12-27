package enaptercli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/urfave/cli/v2"
)

type cmdRuleEngineRuleUpdate struct {
	cmdRuleEngineRule
	ruleID string
	slug   string
	name   string
}

func buildCmdRuleEngineRuleUpdate() *cli.Command {
	cmd := &cmdRuleEngineRuleUpdate{}
	return &cli.Command{
		Name:               "update",
		Usage:              "Update a rule",
		CustomHelpTemplate: cmd.HelpTemplate(),
		Flags:              cmd.Flags(),
		Before:             cmd.Before,
		Action: func(cliCtx *cli.Context) error {
			return cmd.do(cliCtx)
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
		},
		&cli.StringFlag{
			Name:        "name",
			Usage:       "A new rule name",
			Destination: &c.name,
		},
	)
}

func (c *cmdRuleEngineRuleUpdate) do(cliCtx *cli.Context) error {
	payload := struct {
		Rule       map[string]any `json:"rule"`
		UpdateMask string         `json:"update_mask"`
	}{
		Rule:       make(map[string]any),
		UpdateMask: "",
	}

	if cliCtx.IsSet("slug") {
		payload.Rule["slug"] = c.slug
		payload.UpdateMask += "slug,"
	}
	if cliCtx.IsSet("name") {
		payload.Rule["name"] = c.name
		payload.UpdateMask += "name,"
	}

	payload.UpdateMask = strings.TrimSuffix(payload.UpdateMask, ",")

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("build request: %w", err)
	}

	return c.doHTTPRequest(cliCtx.Context, doHTTPRequestParams{
		Method:      http.MethodPatch,
		Path:        "/" + c.ruleID,
		Body:        bytes.NewReader(body),
		ContentType: contentTypeJSON,
	})
}
