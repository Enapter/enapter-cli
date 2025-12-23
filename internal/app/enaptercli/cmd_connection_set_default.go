package enaptercli

import (
	"github.com/urfave/cli/v2"

	"github.com/enapter/enapter-cli/internal/app/configfile"
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
				Usage:       "connection name",
				Destination: &cmd.name,
				Required:    true,
			},
		},
		Action: cmd.do,
	}
}

func (c *cmdConnectionSetDefault) do(*cli.Context) error {
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
