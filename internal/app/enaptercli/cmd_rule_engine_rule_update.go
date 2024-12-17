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
	slug string
	name string
}

func buildCmdRuleEngineRuleUpdate() *cli.Command {
	cmd := &cmdRuleEngineRuleUpdate{}
	return &cli.Command{
		Name:               "update",
		Usage:              "Update a rule",
		Args:               true,
		ArgsUsage:          "RULE",
		CustomHelpTemplate: cmd.HelpTemplate(),
		Flags:              cmd.Flags(),
		Before:             cmd.Before,
		Action: func(cliCtx *cli.Context) error {
			return cmd.do(cliCtx, cliCtx.Args().First())
		},
	}
}

func (c *cmdRuleEngineRuleUpdate) Before(cliCtx *cli.Context) error {
	if err := c.cmdRuleEngineRule.Before(cliCtx); err != nil {
		return err
	}

	if cliCtx.Args().Get(0) == "" {
		return errRequiresAtLeastOneArgument
	}

	return nil
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

func (c *cmdRuleEngineRuleUpdate) do(cliCtx *cli.Context, rule string) error {
	payload := struct {
		Rule       map[string]any `json:"rule"`
		UpdateMask string         `json:"update_mask"`
	}{
		Rule:       make(map[string]any),
		UpdateMask: "",
	}

	if cliCtx.IsSet("slug") {
		payload.Rule["slug"] = c.slug
		payload.UpdateMask = payload.UpdateMask + "slug,"
	}
	if cliCtx.IsSet("name") {
		payload.Rule["name"] = c.name
		payload.UpdateMask = payload.UpdateMask + "name,"
	}

	payload.UpdateMask = strings.TrimSuffix(payload.UpdateMask, ",")

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("build request: %w", err)
	}

	return c.doHTTPRequest(cliCtx.Context, doHTTPRequestParams{
		Method: http.MethodPatch,
		Path:   "/" + rule,
		Body:   bytes.NewReader(body),
	})
}
