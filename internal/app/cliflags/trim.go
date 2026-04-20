package cliflags

import (
	"context"
	"strings"

	"github.com/urfave/cli/v3"
)

// TrimSpaceAction returns a cli.StringFlag Action that trims leading and
// trailing whitespace from the flag value and stores it in dest.
func TrimSpaceAction(dest *string) func(context.Context, *cli.Command, string) error {
	return func(_ context.Context, _ *cli.Command, v string) error {
		*dest = strings.TrimSpace(v)
		return nil
	}
}
