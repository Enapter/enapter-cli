package enaptercli

import (
	"strings"

	"github.com/urfave/cli/v2"
)

type cmdBlueprints struct {
	cmdBase
}

func buildCmdBlueprints() *cli.Command {
	cmd := &cmdBlueprints{}
	return &cli.Command{
		Name:               "blueprint",
		Usage:              "Manage blueprints",
		CustomHelpTemplate: cmd.SubcommandHelpTemplate(),
		Subcommands: []*cli.Command{
			buildCmdBlueprintsProfiles(),
			buildCmdBlueprintsUpload(),
			buildCmdBlueprintsDownload(),
			buildCmdBlueprintsGet(),
		},
	}
}

func isBlueprintID(s string) bool {
	const blueprintIDLen = 36
	if len(s) != blueprintIDLen {
		return false
	}

	isDashPos := func(i int) bool { return i == 8 || i == 13 || i == 18 || i == 23 }
	for i := 0; i < blueprintIDLen; i++ {
		if isDashPos(i) {
			if s[i] != '-' {
				return false
			}
		} else {
			isHexDigit := (s[i] >= '0' && s[i] <= '9') || (s[i] >= 'a' && s[i] <= 'f')
			if !isHexDigit {
				return false
			}
		}
	}
	return true
}

func parseBlueprintName(n string) (name, tag string) {
	const blueprintNameParts = 2
	nameTag := strings.SplitN(n, ":", blueprintNameParts)
	name = nameTag[0]
	tag = "latest"
	if len(nameTag) > 1 {
		tag = nameTag[1]
	}
	return name, tag
}
