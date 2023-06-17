package enaptercli

import (
	"io"

	"github.com/urfave/cli/v2"
)

type cmdBase struct {
	token         string
	apiHost       string
	graphqlURL    string
	websocketsURL string
	writer        io.Writer
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
			Name:        "gql-api-url",
			Usage:       "Override Cloud API endpoint",
			EnvVars:     []string{"ENAPTER_GQL_API_URL"},
			Hidden:      true,
			Value:       "https://cli.enapter.com/graphql",
			Destination: &c.graphqlURL,
		},
		&cli.StringFlag{
			Name:        "ws-api-url",
			Usage:       "Override Cloud API endpoint",
			EnvVars:     []string{"ENAPTER_WS_API_URL"},
			Hidden:      true,
			Value:       "wss://cli.enapter.com/cable",
			Destination: &c.websocketsURL,
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
