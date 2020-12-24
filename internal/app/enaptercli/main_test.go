package enaptercli_test

import (
	"flag"
	"os"
	"testing"
)

//nolint:gochecknoglobals // because it needs for golden files updates.
var update bool

func TestMain(m *testing.M) {
	flag.BoolVar(&update, "update", false, "update golden testfile")
	os.Exit(m.Run())
}
