package enaptercli

import (
	"github.com/urfave/cli/v3"
)

type cmdBlueprintsProfiles struct {
	cmdBase
}

func buildCmdBlueprintsProfiles() *cli.Command {
	return &cli.Command{
		Name:  "profiles",
		Usage: "Manage blueprints profiles",
		Commands: []*cli.Command{
			buildCmdBlueprintsProfilesDownload(),
			buildCmdBlueprintsProfilesUpload(),
		},
	}
}
