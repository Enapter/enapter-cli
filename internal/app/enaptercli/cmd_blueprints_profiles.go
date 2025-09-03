package enaptercli

import (
	"github.com/urfave/cli/v2"
)

type cmdBlueprintsProfiles struct {
	cmdBase
}

func buildCmdBlueprintsProfiles() *cli.Command {
	cmd := &cmdBlueprintsProfiles{}
	return &cli.Command{
		Name:               "profiles",
		Usage:              "Manage blueprint profiles",
		CustomHelpTemplate: cmd.SubcommandHelpTemplate(),
		Subcommands: []*cli.Command{
			buildCmdBlueprintsProfilesDownload(),
			buildCmdBlueprintsProfilesUpload(),
		},
	}
}
