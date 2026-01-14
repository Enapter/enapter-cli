package enaptercli

import "errors"

var (
	errSiteIDMismatch = errors.New("passed site-ID must match the site-ID of the current connection")
	errSiteIDMissing  = errors.New("site ID is required, " +
		"specify --site-id or select a connection with a configured site ID")
	errUnsupportedFlagValue = errors.New("unsupported flag value")
	errOnlyOneBlueprinFlag  = errors.New("only one of --blueprint-id or --blueprint-path can be specified")
	errMissedBlueprintFlag  = errors.New("one of --blueprint-id or --blueprint-path must be specified")
)
