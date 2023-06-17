package enaptercli

import (
	"io"

	"github.com/urfave/cli/v2"
)

type cmdBase struct {
	token        string
	apiHost      string
	cloudAPIHost string
	writer       io.Writer
}

func (c *cmdBase) Flags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "token",
			Usage:       "Enapter API token",
			EnvVars:     []string{"ENAPTER_API_TOKEN"},
			Hidden:      true,
			Destination: &c.token,
		},
		&cli.StringFlag{
			Name:        "api-host",
			Usage:       "Override API endpoint",
			EnvVars:     []string{"ENAPTER_API_HOST"},
			Hidden:      true,
			Value:       "https://api.enapter.com",
			Destination: &c.apiHost,
		},
		&cli.StringFlag{
			Name:        "cloud-api-host",
			Usage:       "Override Cloud API endpoint",
			EnvVars:     []string{"ENAPTER_CLOUD_API_HOST"},
			Hidden:      true,
			Value:       "cli.enapter.com",
			Destination: &c.cloudAPIHost,
		},
	}
}

func (c *cmdBase) Before(cliCtx *cli.Context) error {
	if cliCtx.String("token") == "" {
		return errAPITokenMissed
	}
	c.writer = cliCtx.App.Writer
	return nil
}

func (c *cmdBase) HelpTemplate() string {
	return cli.CommandHelpTemplate + `ENVIRONMENT VARIABLES:
   ENAPTER_API_TOKEN  Enapter API access token

`
}
