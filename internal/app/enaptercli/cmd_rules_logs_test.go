//nolint:dupl // not a duplicate of `devices logs` command tests
package enaptercli_test

import (
	"strings"
	"testing"
)

func TestRuleLogs(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		inputFileName := "testdata/rules_logs/simple/input"
		untilLinePrefix := "[info]"
		expectedFileName := "testdata/rules_logs/simple/output"
		testRuleLogs(t, inputFileName, untilLinePrefix, expectedFileName)
	})

	t.Run("invalid token", func(t *testing.T) {
		inputFileName := "testdata/rules_logs/disconnect/invalid_token/input"
		untilLinePrefix := "[connection]"
		expectedFileName := "testdata/rules_logs/disconnect/invalid_token/output"
		testRuleLogs(t, inputFileName, untilLinePrefix, expectedFileName)
	})

	t.Run("rule not found", func(t *testing.T) {
		inputFileName := "testdata/rules_logs/disconnect/rule_not_found/input"
		untilLinePrefix := "[connection] disconnected"
		expectedFileName := "testdata/rules_logs/disconnect/rule_not_found/output"
		testRuleLogs(t, inputFileName, untilLinePrefix, expectedFileName)
	})
}

func testRuleLogs(t *testing.T, inputFileName, untilLinePrefix, expectedFileName string) {
	const hardwareID = "SIM-RULE"

	identifier := map[string]string{"channel": "RuleChannel", "rule_id": hardwareID}

	command := strings.Split("enapter rules logs", " ")
	command = append(command, "--rule-id", hardwareID)

	testLogsCommand(t, inputFileName, untilLinePrefix, expectedFileName, identifier, command)
}
