package enaptercli

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"

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
				Usage:       "Connection name",
				Destination: &cmd.name,
				Required:    true,
			},
		},
		Action: func(_ context.Context, cliCmd *cli.Command) error {
			return cmd.do(cliCmd)
		},
	}
}

func (c *cmdConnectionRemove) do(cliCmd *cli.Command) error {
	config, err := configfile.Load()
	if err != nil {
		return err
	}

	if _, ok := config.Connections[c.name]; !ok {
		fmt.Fprintln(cliCmd.Root().ErrWriter, "WARNING: unknown connection.")
		return nil
	}

	delete(config.Connections, c.name)
	if config.DefaultConn == c.name {
		fmt.Fprintln(cliCmd.Root().ErrWriter, "WARNING: removed connection was set as default.")
		config.DefaultConn = ""
	}

	return configfile.Save(config)
}
