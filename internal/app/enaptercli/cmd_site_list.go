package enaptercli

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/urfave/cli/v2"
)

type cmdSiteList struct {
	cmdSite
	mySites bool
	limit   int
}

func buildCmdSiteList() *cli.Command {
	cmd := &cmdSiteList{}
	return &cli.Command{
		Name:               "list",
		Usage:              "List user sites",
		CustomHelpTemplate: cmd.CommandHelpTemplate(),
		Flags:              cmd.Flags(),
		Before:             cmd.Before,
		Action: func(cliCtx *cli.Context) error {
			return cmd.do(cliCtx.Context)
		},
	}
}

func (c *cmdSiteList) Flags() []cli.Flag {
	flags := c.cmdSite.Flags()
	return append(flags, &cli.BoolFlag{
		Name:        "my-sites",
		Usage:       "Returns only sites where user is owner or installer",
		Destination: &c.mySites,
	}, &cli.IntFlag{
		Name:        "limit",
		Usage:       "Maximum number of sites to retrieve",
		Destination: &c.limit,
		DefaultText: "retrieves all",
	})
}

func (c *cmdSiteList) do(ctx context.Context) error {
	if siteID, _ := c.chooseSiteID(""); siteID != "" {
		fmt.Fprintln(c.errWriter, "WARNING: trying to get sites list when site ID "+
			"is set for current connection, result will contain only one site.")

		var resp struct {
			Site json.RawMessage `json:"site"`
		}
		if err := c.doHTTPRequest(ctx, doHTTPRequestParams{
			Method: http.MethodGet,
			Path:   "/" + c.siteID,
			RespProcessor: func(r *http.Response) error {
				return json.NewDecoder(r.Body).Decode(&resp)
			},
		}); err != nil {
			return err
		}

		return json.NewEncoder(c.writer).Encode(struct {
			Sites      []json.RawMessage `json:"sites"`
			TotalCount int               `json:"total_count"`
		}{
			Sites:      []json.RawMessage{resp.Site},
			TotalCount: 1,
		})
	}

	doPaginateRequestParams := paginateHTTPRequestParams{
		ObjectName: "sites",
		Limit:      c.limit,
		DoFn:       c.doHTTPRequest,
		BaseParams: doHTTPRequestParams{
			Method: http.MethodGet,
			Path:   "",
		},
	}

	if c.mySites {
		doPaginateRequestParams.BaseParams.Path = "/users/me/sites"
	}

	return c.doPaginateRequest(ctx, doPaginateRequestParams)
}
