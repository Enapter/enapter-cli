package enaptercli

import (
	"github.com/urfave/cli/v3"
)

func buildCmdConnection() *cli.Command {
	return &cli.Command{
		Name:  "connection",
		Usage: "Manage connections to Enapter Cloud and Gateways",
		Commands: []*cli.Command{
			buildCmdConnectionAdd(),
			buildCmdConnectionRemove(),
			buildCmdConnectionList(),
			buildCmdConnectionSetDefault(),
		},
	}
}
