package enaptercli

import (
	"github.com/urfave/cli/v2"
)

type cmdProvisioning struct {
	cmdBase
}

func buildCmdProvisioning() *cli.Command {
	return &cli.Command{
		Name:  "provisioning",
		Usage: "Create devices of different types",
		Subcommands: []*cli.Command{
			buildCmdProvisioningStandalone(),
			buildCmdProvisioningLua(),
		},
	}
}
