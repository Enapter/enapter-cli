package enaptercli

import "errors"

var (
	errUnsupportedFlagValue = errors.New("unsupported flag value")
	errOnlyOneBlueprinFlag  = errors.New("only one of --blueprint-id or --blueprint-path can be specified")
	errMissedBlueprintFlag  = errors.New("one of --blueprint-id or --blueprint-path must be specified")
)
