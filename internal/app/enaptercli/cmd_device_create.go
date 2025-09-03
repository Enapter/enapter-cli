package enaptercli

import (
	"github.com/urfave/cli/v2"
)

type cmdDeviceCreate struct {
	cmdBase
}

func buildCmdDeviceCreate() *cli.Command {
	cmd := &cmdDeviceCreate{}
	return &cli.Command{
		Name:               "create",
		Usage:              "Create devices of different types",
		CustomHelpTemplate: cmd.SubcommandHelpTemplate(),
		Subcommands: []*cli.Command{
			buildCmdDeviceCreateStandalone(),
			buildCmdDeviceCreateLua(),
		},
	}
}
