package enaptercli

import (
	"context"
	"fmt"
	"maps"
	"slices"
	"text/tabwriter"

	"github.com/urfave/cli/v3"

	"github.com/enapter/enapter-cli/v3/internal/app/configfile"
)

type cmdConnectionList struct{}

func buildCmdConnectionList() *cli.Command {
	cmd := &cmdConnectionList{}
	return &cli.Command{
		Name:  "list",
		Usage: "List all connections",
		Action: func(_ context.Context, cliCmd *cli.Command) error {
			return cmd.do(cliCmd)
		},
	}
}

func (c *cmdConnectionList) do(cliCmd *cli.Command) error {
	config, err := configfile.Load()
	if err != nil {
		return err
	}

	const padding = 3
	w := tabwriter.NewWriter(cliCmd.Root().Writer, 0, 0, padding, ' ', 0)

	fmt.Fprintln(w, "NAME\tTYPE\tURL\tALLOW INSECURE\tSITE ID")

	names := slices.Sorted(maps.Keys(config.Connections))
	for _, name := range names {
		conn := config.Connections[name]

		displayName := name
		if name == config.DefaultConn {
			displayName += " *"
		}

		typ := "cloud"
		if conn.Gateway {
			typ = "gateway"
		}

		allowInsecure := "no"
		if conn.AllowInsecure {
			allowInsecure = "yes"
		}

		fmt.Fprintf(w, "%s\t%s\t%v\t%v\t%v\n",
			displayName, typ, conn.URL, allowInsecure, conn.SiteID)
	}

	return w.Flush()
}
