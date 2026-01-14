package enaptercli

import (
	"github.com/urfave/cli/v2"
)

type cmdBlueprintProfiles struct {
	cmdBase
}

func buildCmdBlueprintProfiles() *cli.Command {
	cmd := &cmdBlueprintProfiles{}
	return &cli.Command{
		Name:               "profiles",
		Usage:              "Manage blueprint profiles",
		CustomHelpTemplate: cmd.SubcommandHelpTemplate(),
		Subcommands: []*cli.Command{
			buildCmdBlueprintProfilesDownload(),
			buildCmdBlueprintProfilesUpload(),
		},
	}
}
