package enaptercli

import "errors"

var (
	errAPITokenMissed = errors.New("API token missing. Set it up using environment " +
		"variable ENAPTER3_API_TOKEN")
	errBlueprintIDMissed    = errors.New("blueprint ID is missed")
	errBlueprintPathMissed  = errors.New("blueprint path is missed")
	errUnsupportedFlagValue = errors.New("unsupported flag value")
)
