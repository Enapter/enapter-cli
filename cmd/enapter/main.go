package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"runtime"

	"github.com/urfave/cli/v2"

	"github.com/enapter/enapter-cli/internal/app/enaptercli"
)

//nolint:gochecknoglobals // because sets up via ldflags
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Printf("Enapter CLI %s, commit %s, built at %s, Go version %s\n",
			c.App.Version, commit, date, runtime.Version())
	}

	app := enaptercli.NewApp()
	app.Version = version

	interruptCh := make(chan os.Signal, 1)
	signal.Notify(interruptCh, os.Interrupt)
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		<-interruptCh
		cancel()
		fmt.Println("Stopping... Press Ctrl+C second time to force stop.")
		<-interruptCh
		os.Exit(1)
	}()

	if err := app.RunContext(ctx, os.Args); err != nil {
		fmt.Println("")
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
