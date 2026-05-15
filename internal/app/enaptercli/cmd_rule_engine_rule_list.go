package enaptercli

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/urfave/cli/v3"
)

type cmdRuleEngineRuleList struct {
	cmdRuleEngineRule
	limit int
}

func buildCmdRuleEngineRuleList() *cli.Command {
	cmd := &cmdRuleEngineRuleList{}
	return &cli.Command{
		Name:               "list",
		Usage:              "List rules",
		CustomHelpTemplate: cmd.CommandHelpTemplate(),
		Flags:              cmd.Flags(),
		Before:             cmd.Before,
		Action: func(ctx context.Context, _ *cli.Command) error {
			return cmd.do(ctx)
		},
	}
}

func (c *cmdRuleEngineRuleList) Flags() []cli.Flag {
	flags := c.cmdRuleEngineRule.Flags()
	return append(flags, &cli.IntFlag{
		Name:        "limit",
		Usage:       "maximum number of rules to retrieve",
		Destination: &c.limit,
		DefaultText: "retrieves all",
	})
}

func (c *cmdRuleEngineRuleList) do(ctx context.Context) error {
	body, err := c.doProbeRequest(ctx)
	if err != nil {
		return err
	}

	if len(body.Rules) == 0 {
		return c.printEmptyResponse()
	}

	if body.TotalCount == 0 {
		return c.doLegacyUnpaginatedPath(body)
	}

	return c.doPaginatedPath(ctx)
}

type ruleListResponse struct {
	Rules      []json.RawMessage `json:"rules"`
	TotalCount int               `json:"total_count"`
}

func (c *cmdRuleEngineRuleList) doProbeRequest(ctx context.Context) (*ruleListResponse, error) {
	var body ruleListResponse
	err := c.doHTTPRequest(ctx, doHTTPRequestParams{
		Method: http.MethodGet,
		Path:   "",
		RespProcessor: func(resp *http.Response) error {
			if resp.StatusCode != http.StatusOK {
				return cli.Exit(parseRespErrorMessage(resp), 1)
			}
			if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
				return fmt.Errorf("decode response: %w", err)
			}
			return nil
		},
	})
	if err != nil {
		return nil, err
	}
	return &body, nil
}

func (c *cmdRuleEngineRuleList) printEmptyResponse() error {
	respBytes, err := json.Marshal(map[string]any{
		"rules":       []json.RawMessage{},
		"total_count": 0,
	})
	if err != nil {
		return fmt.Errorf("marshal empty response: %w", err)
	}

	return c.defaultRespProcessor(&http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewReader(respBytes)),
	})
}

func (c *cmdRuleEngineRuleList) doLegacyUnpaginatedPath(body *ruleListResponse) error {
	body.TotalCount = len(body.Rules)

	if c.limit > 0 && len(body.Rules) > c.limit {
		body.Rules = body.Rules[:c.limit]
	}

	respBytes, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("marshal legacy response: %w", err)
	}

	return c.defaultRespProcessor(&http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewReader(respBytes)),
	})
}

func (c *cmdRuleEngineRuleList) doPaginatedPath(ctx context.Context) error {
	return c.doPaginateRequest(ctx, paginateHTTPRequestParams{
		ObjectName: "rules",
		Limit:      c.limit,
		DoFn:       c.doHTTPRequest,
		BaseParams: doHTTPRequestParams{
			Method: http.MethodGet,
			Path:   "",
		},
	})
}
