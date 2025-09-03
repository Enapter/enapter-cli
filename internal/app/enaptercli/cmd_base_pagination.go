package enaptercli

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/urfave/cli/v2"
)

var errEndPagination = errors.New("end pagination")

type paginateHTTPRequestParams struct {
	BaseParams doHTTPRequestParams
	DoFn       func(ctx context.Context, p doHTTPRequestParams) error
	Limit      int
	ObjectName string
}

func (c *cmdBase) doPaginateRequest(ctx context.Context, p paginateHTTPRequestParams) error {
	const maxPageLimit = 50
	if p.BaseParams.Query == nil {
		p.BaseParams.Query = url.Values{}
	}
	if p.Limit > 0 && p.Limit < maxPageLimit {
		p.BaseParams.Query.Set("limit", strconv.Itoa(p.Limit))
		return p.DoFn(ctx, p.BaseParams)
	}

	paginateRespProcesor := &paginateRespProcesor{
		ObjectName:  p.ObjectName,
		seenObjects: make(map[string]struct{}),
	}
	for {
		reqPageParams := p.BaseParams
		reqPageParams.Query.Set("offset", strconv.Itoa(len(paginateRespProcesor.Objects)))
		reqPageParams.Query.Set("limit", strconv.Itoa(maxPageLimit))
		reqPageParams.RespProcessor = paginateRespProcesor.Process

		err := p.DoFn(ctx, reqPageParams)
		if err != nil {
			if errors.Is(err, errEndPagination) {
				break
			}
			return fmt.Errorf("failed to retrieve page: %w", err)
		}
		if p.Limit > 0 && len(paginateRespProcesor.Objects) >= p.Limit {
			break
		}
	}

	returnCount := len(paginateRespProcesor.Objects)
	if p.Limit > 0 && returnCount > p.Limit {
		returnCount = p.Limit
	}
	respBytes, err := json.Marshal(map[string]any{
		"total_count": paginateRespProcesor.TotalCount,
		p.ObjectName:  paginateRespProcesor.Objects[:returnCount],
	})
	if err != nil {
		return cli.Exit("Failed to marshal response: "+err.Error(), 1)
	}
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewReader(respBytes)),
	}
	return c.defaultRespProcessor(resp)
}

type paginateRespProcesor struct {
	TotalCount  int
	Objects     []any
	ObjectName  string
	seenObjects map[string]struct{}
}

func (p *paginateRespProcesor) Process(resp *http.Response) error {
	if resp.StatusCode != http.StatusOK {
		return cli.Exit("Unexpected response status: "+resp.Status, 1)
	}

	var pageBody map[string]json.RawMessage
	if err := json.NewDecoder(resp.Body).Decode(&pageBody); err != nil {
		return cli.Exit("Failed to parse response: "+err.Error(), 1)
	}

	if err := json.Unmarshal(pageBody["total_count"], &p.TotalCount); err != nil {
		return cli.Exit("Failed to parse total_count: "+err.Error(), 1)
	}

	var objects []json.RawMessage
	if err := json.Unmarshal(pageBody[p.ObjectName], &objects); err != nil {
		return cli.Exit("Failed to parse "+p.ObjectName+": "+err.Error(), 1)
	}

	if len(objects) == 0 {
		return errEndPagination
	}

	for _, obj := range objects {
		var objMap map[string]any
		if err := json.Unmarshal(obj, &objMap); err != nil {
			return cli.Exit("Failed to parse object: "+err.Error(), 1)
		}

		id, ok := objMap["id"].(string)
		if !ok || id == "" {
			return cli.Exit("Object ID is missing or not a string", 1)
		}

		if _, seen := p.seenObjects[id]; !seen {
			p.seenObjects[id] = struct{}{}
			p.Objects = append(p.Objects, objMap)
		}
	}
	return nil
}
