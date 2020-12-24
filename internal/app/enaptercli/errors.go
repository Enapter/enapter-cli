package enaptercli

import "errors"

var (
	errFinishedWithError = errors.New("request execution failed")
	errLogStatusError    = errors.New("error during request execution")
	errAPITokenMissed    = errors.New("API token missing. Set it up using environment " +
		"variable ENAPTER_API_TOKEN")
	errRequestTimedOut = errors.New("request timed out")
	errFinished        = errors.New("finished")
)
