package enaptercli_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bxcodec/faker/v3"
	"github.com/stretchr/testify/require"
)

func TestRulesUpdate(t *testing.T) {
	testdataDir := "testdata/rules_update"
	dirs, err := os.ReadDir(testdataDir)
	require.NoError(t, err)

	for _, dir := range dirs {
		if !dir.IsDir() {
			continue
		}

		dir := dir
		t.Run(dir.Name(), func(t *testing.T) {
			testRulesUpdate(t, filepath.Join(testdataDir, dir.Name()))
		})
	}
}

func TestRulesUpdateWrongFilePath(t *testing.T) {
	args := strings.Split("enapter rules update --token token --rule-id ruleID "+
		"--gql-api-url apiURL --rule-path wrong", " ")
	app := startTestApp(args...)
	defer app.Stop()

	appErr := app.Wait()
	require.EqualError(t, appErr, "read rule file: open wrong: no such file or directory")
}

type rulesUpdateTestSettings struct {
	RuleID   string `json:"rule_id"`
	RulePath string `json:"rule_path"`
	Token    string `json:"-"`
}

func (s *rulesUpdateTestSettings) Fill(t *testing.T, dir string) {
	settingsBytes, err := os.ReadFile(filepath.Join(dir, "settings.json"))
	require.NoError(t, err)
	require.NoError(t, json.Unmarshal(settingsBytes, s))

	s.RulePath = filepath.Join(dir, s.RulePath)
	s.Token = faker.Word()
}

func testRulesUpdate(t *testing.T, dir string) {
	var settings rulesUpdateTestSettings
	settings.Fill(t, dir)

	reqs := byteSliceSliceFromFile(t, filepath.Join(dir, "requests"))
	resps := byteSliceSliceFromFile(t, filepath.Join(dir, "responses"))

	srv := startTestServer(reqs, resps, "")
	defer srv.Close()

	args := strings.Split("enapter rules update", " ")
	args = append(args,
		"--token", settings.Token,
		"--rule-id", settings.RuleID,
		"--rule-path", settings.RulePath,
		"--gql-api-url", srv.URL)

	checkTestAppOutput(t, dir, args, reqs)
}
