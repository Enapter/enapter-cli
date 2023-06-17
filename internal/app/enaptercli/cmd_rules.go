package enaptercli

import "github.com/urfave/cli/v2"

type cmdRules struct {
	cmdBase
	ruleID string
}

func buildCmdRules() *cli.Command {
	return &cli.Command{
		Name:  "rules",
		Usage: "Rules information and management commands.",
		Subcommands: []*cli.Command{
			buildCmdRulesUpdate(),
			buildCmdRulesLogs(),
		},
	}
}

func (c *cmdRules) Flags() []cli.Flag {
	flags := c.cmdBase.Flags()
	flags = append(flags, &cli.StringFlag{
		Name:        "rule-id",
		Usage:       "Rule ID; can be obtained in cloud.enapter.com",
		Required:    true,
		Destination: &c.ruleID,
	})
	return flags
}
