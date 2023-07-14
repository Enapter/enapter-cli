package enaptercli

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/urfave/cli/v2"

	"github.com/enapter/enapter-cli/internal/cloudapi"
)

const deviceUploadDefaultTimeout = 30 * time.Second

type cmdDevicesUpload struct {
	cmdDevices
	blueprintDir string
	timeout      time.Duration
}

func buildCmdDevicesUpload() *cli.Command {
	cmd := &cmdDevicesUpload{}

	return &cli.Command{
		Name:  "upload",
		Usage: "Upload blueprint to a device",
		Description: "Blueprint combines device capabilities declaration and Lua firmware for Enapter UCM. " +
			"The command updates device blueprint and uploads the firmware to the UCM. Learn more " +
			"about Enapter Blueprints at https://handbook.enapter.com/blueprints.",
		CustomHelpTemplate: cmd.HelpTemplate(),
		Flags:              cmd.Flags(),
		Before:             cmd.Before,
		Action: func(cliCtx *cli.Context) error {
			return cmd.upload(cliCtx.Context, cliCtx.App.Version)
		},
	}
}

func (c *cmdDevicesUpload) Flags() []cli.Flag {
	flags := c.cmdDevices.Flags()
	flags = append(flags,
		&cli.DurationFlag{
			Name:        "timeout",
			Usage:       "Time to wait for blueprint uploading",
			Destination: &c.timeout,
			Value:       deviceUploadDefaultTimeout,
		},
		&cli.StringFlag{
			Name:        "blueprint-dir",
			Usage:       "Directory which contains blueprint file",
			Required:    true,
			Destination: &c.blueprintDir,
		},
	)
	return flags
}

func (c *cmdDevicesUpload) upload(ctx context.Context, version string) error {
	if c.timeout != 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, c.timeout)
		defer cancel()
	}

	files, err := c.blueprintFilesList()
	if err != nil {
		return err
	}

	fmt.Fprintln(c.writer, "Blueprint files to be uploaded:")
	for _, name := range files {
		fmt.Fprintln(c.writer, "*", name)
	}

	zipBytes, err := c.blueprintZip()
	if err != nil {
		return err
	}

	onceWriter := &onceWriter{w: c.writer}
	transport := cloudapi.NewCredentialsTransport(http.DefaultTransport, c.token, version)
	transport = cloudapi.NewCLIMessageWriterTransport(transport, onceWriter)
	client := cloudapi.NewClientWithURL(&http.Client{Transport: transport}, c.graphqlURL)

	uploadData, uploadErrors, err := client.UploadBlueprint(ctx, c.hardwareID, zipBytes)
	if err != nil {
		return fmt.Errorf("do update: %w", err)
	}

	if len(uploadErrors) != 0 {
		for _, e := range uploadErrors {
			fmt.Fprintln(c.writer, "[ERROR]", e.Message)
		}
		return errFinishedWithError
	}

	fmt.Fprintln(c.writer, "upload started with operation id", uploadData.OperationID)

	err = client.WriteOperationLogs(ctx, c.hardwareID, uploadData.OperationID, c.writeLog)
	if err != nil {
		return fmt.Errorf("receive operation logs: %w", err)
	}

	fmt.Fprintln(c.writer, "Done!")
	return nil
}

func (c *cmdDevicesUpload) blueprintFilesList() ([]string, error) {
	var files []string

	err := filepath.Walk(c.blueprintDir,
		func(name string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				files = append(files, name)
			}
			return nil
		})
	if err != nil {
		return nil, err
	}

	return files, nil
}

func (c *cmdDevicesUpload) blueprintZip() ([]byte, error) {
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

	return zipBuf.Bytes(), nil
}

func (c *cmdDevicesUpload) writeLog(operationID string, l cloudapi.OperationLog) {
	fmt.Fprintf(c.writer, "[#%s] %s [%s] %s\n", operationID, l.CreatedAt, l.Severity, l.Payload)
}
