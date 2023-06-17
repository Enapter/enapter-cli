package cloudapi

import "errors"

var (
	ErrRequestTimedOut = errors.New("request timed out")
	ErrFinished        = errors.New("finished")
)
