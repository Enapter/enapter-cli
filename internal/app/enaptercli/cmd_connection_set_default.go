package enaptercli

import (
	"context"

	"github.com/urfave/cli/v3"

	"github.com/enapter/enapter-cli/v3/internal/app/configfile"
)

type cmdConnectionSetDefault struct {
	name string
}

func buildCmdConnectionSetDefault() *cli.Command {
	cmd := &cmdConnectionSetDefault{}
	return &cli.Command{
		Name:  "set-default",
		Usage: "Set default connection",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "name",
				Usage:       "Connection name",
				Destination: &cmd.name,
				Required:    true,
			},
		},
		Action: cmd.do,
	}
}

func (c *cmdConnectionSetDefault) do(context.Context, *cli.Command) error {
	config, err := configfile.Load()
	if err != nil {
		return err
	}

	if _, ok := config.Connections[c.name]; !ok {
		return cli.Exit("Unknown connection.", 1)
	}

	config.DefaultConn = c.name
	return configfile.Save(config)
}
