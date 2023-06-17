package cloudapi

import (
	"context"
	"errors"
	"fmt"

	"github.com/shurcooL/graphql"
)

type UpdateRuleInput struct {
	RuleID            string
	LuaCode           string
	StdlibVersion     string
	ExecutionInterval int
}

type UpdateRuleData struct {
	Code    string
	Message string
	Title   string
}

type UpdateRuleError struct {
	Code    string
	Message string
	Path    []string
	Title   string
}

func (c *Client) UpdateRule(
	ctx context.Context, input UpdateRuleInput,
) (UpdateRuleData, []UpdateRuleError, error) {
	client := graphql.NewClient(c.host, c.httpClient)

	var mutation struct {
		Rule struct {
			Update struct {
				Data   UpdateRuleData
				Errors []UpdateRuleError
			} `graphql:"update(input: $input)"`
		}
	}

	type UpdateInput struct {
		RuleID            graphql.String `json:"ruleId"`
		LuaCode           graphql.String `json:"luaCode"`
		StdlibVersion     graphql.String `json:"stdlibVersion,omitempty"`
		ExecutionInterval graphql.Int    `json:"executionInterval,omitempty"`
	}

	variables := map[string]interface{}{
		"input": UpdateInput{
			RuleID:            graphql.String(input.RuleID),
			LuaCode:           graphql.String(input.LuaCode),
			StdlibVersion:     graphql.String(input.StdlibVersion),
			ExecutionInterval: graphql.Int(input.ExecutionInterval),
		},
	}

	if err := client.Mutate(ctx, &mutation, variables); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			err = ErrRequestTimedOut
		}
		return UpdateRuleData{}, nil, fmt.Errorf("mutate: %w", err)
	}

	return mutation.Rule.Update.Data, mutation.Rule.Update.Errors, nil
}
