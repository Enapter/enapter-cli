package enaptercli

import "errors"

var (
	errAPITokenMissed = errors.New("API token missing. Set it up using environment " +
		"variable ENAPTER3_API_TOKEN")
	errUnsupportedFlagValue = errors.New("unsupported flag value")
)
