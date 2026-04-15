package enaptercli

import (
	"github.com/urfave/cli/v3"
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
		Commands: []*cli.Command{
			buildCmdBlueprintProfilesDownload(),
			buildCmdBlueprintProfilesUpload(),
		},
	}
}
