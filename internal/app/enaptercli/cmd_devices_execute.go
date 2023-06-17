package enaptercli

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/urfave/cli/v2"

	"github.com/enapter/enapter-cli/internal/publichttp"
)

type cmdDevicesExecute struct {
	cmdDevices
	commandName  string
	arguments    string
	showProgress bool
}

func buildCmdDevicesExecute() *cli.Command {
	cmd := &cmdDevicesExecute{}

	return &cli.Command{
		Name:               "execute",
		Usage:              "Execute command on device",
		CustomHelpTemplate: cmd.HelpTemplate(),
		Flags:              cmd.Flags(),
		Before:             cmd.Before,
		Action: func(cliCtx *cli.Context) error {
			return cmd.execute(cliCtx.Context)
		},
	}
}

func (c *cmdDevicesExecute) Flags() []cli.Flag {
	flags := c.cmdDevices.Flags()
	flags = append(flags,
		&cli.StringFlag{
			Name:        "command",
			Usage:       "Command name",
			Required:    true,
			Destination: &c.commandName,
		},
		&cli.StringFlag{
			Name:        "arguments",
			Usage:       "Command arguments as JSON object",
			Destination: &c.arguments,
		},
		&cli.BoolFlag{
			Name:        "show-progress",
			Usage:       "Enable in-progress responses streaming",
			Destination: &c.showProgress,
		},
	)
	return flags
}

func (c *cmdDevicesExecute) execute(ctx context.Context) error {
	transport := publichttp.NewAuthTokenTransport(http.DefaultTransport, c.token)
	client, err := publichttp.NewClientWithURL(&http.Client{Transport: transport}, c.apiHost)
	if err != nil {
		return fmt.Errorf("create http client: %w", err)
	}

	var arguments map[string]interface{}
	if c.arguments != "" {
		if err := json.Unmarshal([]byte(c.arguments), &arguments); err != nil {
			return fmt.Errorf("parse arguments: %w", err)
		}
	}

	query := publichttp.CommandQuery{
		HardwareID:  c.hardwareID,
		CommandName: c.commandName,
		Arguments:   arguments,
	}

	if c.showProgress {
		return c.executeWithProgress(ctx, client, query)
	}

	response, err := client.Commands.Execute(ctx, query)
	if err != nil {
		return err
	}
	return c.print(response)
}

func (c *cmdDevicesExecute) executeWithProgress(
	ctx context.Context, client *publichttp.Client,
	query publichttp.CommandQuery,
) error {
	progressCh, err := client.Commands.ExecuteWithProgress(ctx, query)
	if err != nil {
		return err
	}

	for progress := range progressCh {
		if progress.Error != nil {
			return progress.Error
		}
		err := c.print(progress.CommandResponse)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *cmdDevicesExecute) print(r publichttp.CommandResponse) error {
	s, err := json.Marshal(r)
	if err != nil {
		return fmt.Errorf("format response: %w", err)
	}
	fmt.Fprintln(c.writer, string(s))
	return nil
}
