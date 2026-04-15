package enaptercli

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/urfave/cli/v3"
)

type cmdDeviceCommunicationConfigGenerate struct {
	cmdDeviceCommunicationConfig
	protocol string
}

func buildCmdDeviceCommunicationConfigGenerate() *cli.Command {
	cmd := &cmdDeviceCommunicationConfigGenerate{}
	return &cli.Command{
		Name:               "generate",
		Usage:              "Generate a new communication config for device",
		CustomHelpTemplate: cmd.CommandHelpTemplate(),
		Flags:              cmd.Flags(),
		Before:             cmd.Before,
		Action: func(ctx context.Context, _ *cli.Command) error {
			return cmd.do(ctx)
		},
	}
}

func (c *cmdDeviceCommunicationConfigGenerate) Flags() []cli.Flag {
	flags := c.cmdDeviceCommunicationConfig.Flags()
	return append(flags,
		&cli.StringFlag{
			Name:        "protocol",
			Usage:       "Connection protocol (supported values: MQTT, MQTTS)",
			Destination: &c.protocol,
			Required:    true,
		},
	)
}

func (c *cmdDeviceCommunicationConfigGenerate) do(ctx context.Context) error {
	reqBody := struct {
		Protocol string `json:"protocol"`
	}{
		Protocol: c.protocol,
	}
	data, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("build request: %w", err)
	}

	return c.doHTTPRequest(ctx, doHTTPRequestParams{
		Method: http.MethodPost,
		Path:   "/generate_config",
		Body:   bytes.NewReader(data),
	})
}
