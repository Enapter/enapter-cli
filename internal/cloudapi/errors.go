package cloudapi

import "errors"

var (
	ErrFinishedWithError = errors.New("request execution failed")
	ErrLogStatusError    = errors.New("error during request execution")
	ErrRequestTimedOut   = errors.New("request timed out")
	ErrFinished          = errors.New("finished")
)
