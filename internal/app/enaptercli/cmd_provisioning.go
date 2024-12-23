package enaptercli

import (
	"github.com/urfave/cli/v3"
)

type cmdProvisioning struct {
	cmdBase
}

func buildCmdProvisioning() *cli.Command {
	return &cli.Command{
		Name:  "provisioning",
		Usage: "Create devices of different types",
		Commands: []*cli.Command{
			buildCmdProvisioningStandalone(),
			buildCmdProvisioningLua(),
		},
	}
}
