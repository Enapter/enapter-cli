package enaptercli

import (
	"github.com/urfave/cli/v3"
)

// NewApp creates a new Enapter CLI tool application instance.
func NewApp() *cli.Command {
	cmd := &cli.Command{}

	cmd.Name = "enapter3"
	cmd.Usage = "Command Line Interface (CLI) for Enapter services."
	cmd.Description = "The Enapter CLI requires an access token for authentication. " +
		"You can obtain the token in your Enapter Cloud account settings."
	cmd.CustomRootCommandHelpTemplate = cli.RootCommandHelpTemplate + enapterAPIEnvVarsHelp

	cli.ShowSubcommandHelp = func(cmd *cli.Command) error {
		tmpl := cmd.CustomHelpTemplate
		if tmpl == "" {
			tmpl = cli.SubcommandHelpTemplate
		}
		cli.HelpPrinter(cmd.Root().Writer, tmpl, cmd)
		return nil
	}

	cmd.Commands = []*cli.Command{
		buildCmdSite(),
		buildCmdDevice(),
		buildCmdBlueprint(),
		buildCmdRuleEngine(),
		buildCmdConnection(),
	}

	return cmd
}
