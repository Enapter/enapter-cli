package enaptercli

import (
	"github.com/urfave/cli/v3"
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
		Commands: []*cli.Command{
			buildCmdDeviceCreateStandalone(),
			buildCmdDeviceCreateLua(),
		},
	}
}
