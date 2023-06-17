package enaptercli

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/shurcooL/graphql"
	"github.com/urfave/cli/v2"
)

const deviceUploadDefaultTimeout = 30 * time.Second

type cmdDevicesUpload struct {
	cmdDevicesUploadCommon
	blueprintDir string
}

type cmdDevicesUploadCommon struct {
	cmdDevices
	timeout time.Duration
}

func (c *cmdDevicesUploadCommon) Flags() []cli.Flag {
	flags := c.cmdDevices.Flags()
	flags = append(flags,
		&cli.DurationFlag{
			Name:        "timeout",
			Usage:       "Time to wait for blueprint uploading",
			Destination: &c.timeout,
			Value:       deviceUploadDefaultTimeout,
		},
	)
	return flags
}

func buildCmdDevicesUpload() *cli.Command {
	cmd := &cmdDevicesUpload{}

	flags := cmd.Flags()
	flags = append(flags, &cli.StringFlag{
		Name:        "blueprint-dir",
		Usage:       "Directory which contains blueprint file",
		Required:    true,
		Destination: &cmd.blueprintDir,
	})

	return &cli.Command{
		Name:  "upload",
		Usage: "Upload blueprint to a device",
		Description: "Blueprint combines device capabilities declaration and Lua firmware for Enapter UCM. " +
			"The command updates device blueprint and uploads the firmware to the UCM. Learn more " +
			"about Enapter Blueprints at https://handbook.enapter.com/blueprints.",
		CustomHelpTemplate: cmd.HelpTemplate(),
		Flags:              flags,
		Before:             cmd.Before,
		Action: func(cliCtx *cli.Context) error {
			return cmd.upload(cliCtx.Context, cliCtx.App.Version)
		},
	}
}

// UploadBlueprintInput contains mutation input variables.
type UploadBlueprintInput struct {
	Blueprint  graphql.String `json:"blueprint"`
	HardwareID graphql.ID     `json:"hardwareId"`
}

type mutation struct {
	Device *struct {
		UploadBlueprint struct {
			Data   uploadBlueprintData
			Errors []uploadBlueprintError
		} `graphql:"uploadBlueprint(input: $input)"`
	}
}

type uploadBlueprintError struct {
	Code    string
	Message string
	Path    []string
	Title   string
}

type uploadBlueprintData struct {
	Code        string
	Message     string
	Title       string
	OperationID string
}

func (c *cmdDevicesUpload) upload(ctx context.Context, version string) error {
	if c.timeout != 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, c.timeout)
		defer cancel()
	}

	zipBuf, err := c.blueprintZipBuf()
	if err != nil {
		return err
	}

	onceWriter := &onceWriter{w: c.writer}
	m, err := c.sendRequest(ctx, zipBuf, onceWriter, version)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}

	if m.Device == nil {
		return fmt.Errorf("device not found: %w", errFinishedWithError)
	}

	if len(m.Device.UploadBlueprint.Errors) != 0 {
		c.dumpUploadErrors(m.Device.UploadBlueprint.Errors)
		return errFinishedWithError
	}

	fmt.Fprintln(c.writer, "upload started with operation id", m.Device.UploadBlueprint.Data.OperationID)

	op := cmdDevicesUploadLogs{
		cmdDevicesUploadCommon: c.cmdDevicesUploadCommon,
	}
	if err := op.logs(ctx, m.Device.UploadBlueprint.Data.OperationID, onceWriter, version); err != nil {
		return err
	}

	fmt.Fprintln(c.writer, "Done!")
	return nil
}

func (c *cmdDevicesUpload) blueprintZipBuf() (*bytes.Buffer, error) {
	bpBytes, err := zipDir(c.blueprintDir)
	if err != nil {
		return nil, fmt.Errorf("failed to zip blueprint dir %q: %w", c.blueprintDir, err)
	}

	zipBuf := &bytes.Buffer{}
	zipBuf.WriteString("data:application/gzip;base64,")
	enc := base64.NewEncoder(base64.StdEncoding, zipBuf)
	_, err = enc.Write(bpBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to encode blueprint as base64: %w", err)
	}

	if err := enc.Close(); err != nil {
		return nil, fmt.Errorf("failed to encode blueprint as base64: %w", err)
	}

	return zipBuf, nil
}

func (c *cmdDevicesUpload) sendRequest(
	ctx context.Context, blueprintBuf *bytes.Buffer, onceWriter *onceWriter, version string,
) (*mutation, error) {
	extraHeaders := map[string][]string{
		"Authorization":         {"Bearer " + c.token},
		"X-ENAPTER-CLI-VERSION": {version},
	}
	httpClient := &http.Client{
		Transport: extraHeaderRoundTripper{
			tripper: cliMessageRoundTripper{
				tripper: http.DefaultTransport,
				writer:  onceWriter,
			},
			extraHeaders: extraHeaders,
		},
	}

	client := graphql.NewClient(c.graphqlURL, httpClient)

	var m mutation
	variables := map[string]interface{}{
		"input": UploadBlueprintInput{
			Blueprint:  graphql.String(blueprintBuf.String()),
			HardwareID: graphql.String(c.hardwareID),
		},
	}

	if err := client.Mutate(ctx, &m, variables); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			err = errRequestTimedOut
		}
		return nil, err
	}

	return &m, nil
}

func (c *cmdDevicesUpload) dumpUploadErrors(errs []uploadBlueprintError) {
	for _, e := range errs {
		fmt.Fprintln(c.writer, "[ERROR]", e.Message)
	}
}

type extraHeaderRoundTripper struct {
	tripper      http.RoundTripper
	extraHeaders map[string][]string
}

func (e extraHeaderRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	newReq := new(http.Request)
	*newReq = *r

	newReq.Header = make(http.Header, len(r.Header)+len(e.extraHeaders))
	for k, s := range r.Header {
		newReq.Header[k] = s
	}
	for k, s := range e.extraHeaders {
		newReq.Header[k] = s
	}

	return e.tripper.RoundTrip(newReq)
}

type cliMessageRoundTripper struct {
	tripper http.RoundTripper
	writer  io.Writer
}

func (c cliMessageRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	resp, err := c.tripper.RoundTrip(r)
	if err != nil {
		return nil, err
	}

	if msg := resp.Header.Get("X-ENAPTER-CLI-MESSAGE"); msg != "" {
		fmt.Fprintln(c.writer, msg)
	}

	return resp, nil
}
