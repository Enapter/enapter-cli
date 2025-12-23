package enaptercli

import (
	"fmt"

	"github.com/urfave/cli/v2"

	"github.com/enapter/enapter-cli/internal/app/configfile"
)

type cmdConnectionRemove struct {
	name string
}

func buildCmdConnectionRemove() *cli.Command {
	cmd := &cmdConnectionRemove{}
	return &cli.Command{
		Name:  "remove",
		Usage: "Remove a connection",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "name",
				Usage:       "connection name",
				Destination: &cmd.name,
				Required:    true,
			},
		},
		Action: cmd.do,
	}
}

func (c *cmdConnectionRemove) do(cliCtx *cli.Context) error {
	config, err := configfile.Load()
	if err != nil {
		return err
	}

	if _, ok := config.Connections[c.name]; !ok {
		fmt.Fprintln(cliCtx.App.ErrWriter, "WARNING: unknown connection.")
		return nil
	}

	delete(config.Connections, c.name)
	if config.DefaultConn == c.name {
		fmt.Fprintln(cliCtx.App.ErrWriter, "WARNING: removed connection was set as default.")
		config.DefaultConn = ""
	}

	return configfile.Save(config)
}
