package enaptercli

import "github.com/urfave/cli/v2"

type cmdDevices struct {
	cmdBase
	hardwareID string
}

func buildCmdDevices() *cli.Command {
	return &cli.Command{
		Name:  "devices",
		Usage: "Device information and management commands.",
		Subcommands: []*cli.Command{
			buildCmdDevicesUpload(),
			buildCmdDevicesLogs(),
			buildCmdDevicesUploadLogs(),
			buildCmdDevicesExecute(),
		},
	}
}

func (c *cmdDevices) Flags() []cli.Flag {
	flags := c.cmdBase.Flags()
	flags = append(flags, &cli.StringFlag{
		Name:        "hardware-id",
		Usage:       "Hardware ID of the device; can be obtained in cloud.enapter.com",
		Required:    true,
		Destination: &c.hardwareID,
	})
	return flags
}
