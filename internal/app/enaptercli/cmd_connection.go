package enaptercli

import (
	"github.com/urfave/cli/v2"
)

func buildCmdConnection() *cli.Command {
	return &cli.Command{
		Name:  "connection",
		Usage: "Manage connections to Enapter Cloud and Gateways",
		Subcommands: []*cli.Command{
			buildCmdConnectionAdd(),
			buildCmdConnectionRemove(),
			buildCmdConnectionList(),
			buildCmdConnectionSetDefault(),
		},
	}
}
