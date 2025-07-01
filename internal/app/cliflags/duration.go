package cliflags

import (
	"github.com/urfave/cli/v2"
)

// Duration is a wrapper around cli.DurationFlag to implement cli.Flag interface.
// It differs from cli.DurationFlag in that it does not return a default text if the value is zero.
type Duration struct {
	cli.DurationFlag
}

var (
	_ cli.Flag              = (*Duration)(nil)
	_ cli.DocGenerationFlag = (*Duration)(nil)
)

func (d *Duration) String() string {
	return cli.FlagStringer(d)
}

func (d *Duration) GetDefaultText() string {
	if d.Value == 0 {
		return ""
	}
	return d.DurationFlag.GetDefaultText()
}
