package enaptercli_test

import (
	"strings"
	"testing"
)

func TestDeviceLogs(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		inputFileName := "testdata/device_logs/simple/input"
		untilLinePrefix := "[telemetry]"
		expectedFileName := "testdata/device_logs/simple/output"
		testDeviceLogs(t, inputFileName, untilLinePrefix, expectedFileName)
	})

	t.Run("invalid token", func(t *testing.T) {
		inputFileName := "testdata/device_logs/disconnect/invalid_token/input"
		untilLinePrefix := "[connection]"
		expectedFileName := "testdata/device_logs/disconnect/invalid_token/output"
		testDeviceLogs(t, inputFileName, untilLinePrefix, expectedFileName)
	})

	t.Run("device not found", func(t *testing.T) {
		inputFileName := "testdata/device_logs/disconnect/device_not_found/input"
		untilLinePrefix := "[connection] disconnected"
		expectedFileName := "testdata/device_logs/disconnect/device_not_found/output"
		testDeviceLogs(t, inputFileName, untilLinePrefix, expectedFileName)
	})
}

func testDeviceLogs(t *testing.T, inputFileName, untilLinePrefix, expectedFileName string) {
	const hardwareID = "SIM-WTM"

	identifier := map[string]string{"channel": "DeviceChannel", "hardware_id": hardwareID}

	command := strings.Split("enapter devices logs", " ")
	command = append(command, "--hardware-id", hardwareID)

	testLogsCommand(t, inputFileName, untilLinePrefix, expectedFileName, identifier, command)
}
