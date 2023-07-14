package cloudapi

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/shurcooL/graphql"
)

type UploadBlueprintData struct {
	Code        string
	Message     string
	Title       string
	OperationID string
}

type UploadBlueprintError struct {
	Code    string
	Message string
	Path    []string
	Title   string
}

type uploadBlueprintMutation struct {
	Device struct {
		UploadBlueprint struct {
			Data   UploadBlueprintData
			Errors []UploadBlueprintError
		} `graphql:"uploadBlueprint(input: $input)"`
	}
}

func (c *Client) UploadBlueprint(
	ctx context.Context, hardwareID string, blueprint []byte,
) (UploadBlueprintData, []UploadBlueprintError, error) {
	type UploadBlueprintInput struct {
		Blueprint  graphql.String `json:"blueprint"`
		HardwareID graphql.ID     `json:"hardwareId"`
	}

	variables := map[string]interface{}{
		"input": UploadBlueprintInput{
			Blueprint:  graphql.String(blueprint),
			HardwareID: graphql.String(hardwareID),
		},
	}

	var mutation uploadBlueprintMutation
	if err := c.client.Mutate(ctx, &mutation, variables); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			err = ErrRequestTimedOut
		}
		return UploadBlueprintData{}, nil, fmt.Errorf("mutate: %w", err)
	}

	uploadInfo := mutation.Device.UploadBlueprint
	return uploadInfo.Data, uploadInfo.Errors, nil
}

type OperationLog struct {
	Payload   string
	CreatedAt string
	Severity  string
}

type blueprintUpdateOperationQuery struct {
	Device *struct {
		BlueprintUpdateOperation *struct {
			Status string
			Logs   struct {
				Edges []struct {
					Cursor graphql.String
					Node   OperationLog
				}
			} `graphql:"logs(after: $after_cursor)"`
		} `graphql:"blueprintUpdateOperation(id: $operation_id)"`
	} `graphql:"device(hardwareId: $hardware_id)"`
}

func (c *Client) WriteOperationLogs(
	ctx context.Context, hardwareID, operationID string,
	writeLog func(operationID string, log OperationLog),
) error {
	v := map[string]interface{}{
		"after_cursor": graphql.String(""),
		"hardware_id":  hardwareID,
		"operation_id": operationID,
	}

	for {
		var q blueprintUpdateOperationQuery
		if err := c.client.Query(ctx, &q, v); err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				err = ErrRequestTimedOut
			}
			return fmt.Errorf("failed to send request: %w", err)
		}

		if q.Device == nil {
			return fmt.Errorf("%w: device not found", ErrFinishedWithError)
		}

		if q.Device.BlueprintUpdateOperation == nil {
			return fmt.Errorf("%w: operation not found", ErrFinishedWithError)
		}

		status := q.Device.BlueprintUpdateOperation.Status
		logs := q.Device.BlueprintUpdateOperation.Logs

		for _, e := range logs.Edges {
			v["after_cursor"] = e.Cursor
			writeLog(operationID, e.Node)
		}

		if len(logs.Edges) == 0 {
			switch status {
			case "SUCCEEDED":
				return nil
			case "ERROR":
				return ErrLogStatusError
			}
		}

		const logRequestPeriod = 100 * time.Millisecond
		time.Sleep(logRequestPeriod)
	}
}

type blueprintUpdateOperationsQuery struct {
	Device *struct {
		BlueprintUpdateOperations struct {
			Nodes []struct {
				ID string
			}
		} `graphql:"blueprintUpdateOperations(last: $last_int)"`
	} `graphql:"device(hardwareId: $hardware_id)"`
}

func (c *Client) WriteLastOperationsLogs(
	ctx context.Context, hardwareID string, lastOperationsNumber int,
	writeLog func(operationID string, log OperationLog),
) error {
	v := map[string]interface{}{
		"hardware_id": hardwareID,
		"last_int":    graphql.Int(lastOperationsNumber),
	}

	var q blueprintUpdateOperationsQuery
	if err := c.client.Query(ctx, &q, v); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			err = ErrRequestTimedOut
		}
		return fmt.Errorf("failed to send request: %w", err)
	}

	if q.Device == nil {
		return fmt.Errorf("%w: device not found", ErrFinishedWithError)
	}

	for _, op := range q.Device.BlueprintUpdateOperations.Nodes {
		if err := c.WriteOperationLogs(ctx, hardwareID, op.ID, writeLog); err != nil {
			if errors.Is(err, ErrLogStatusError) {
				continue
			}
			return err
		}
	}

	return nil
}
