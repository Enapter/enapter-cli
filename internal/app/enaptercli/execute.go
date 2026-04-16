package enaptercli

import (
	"archive/zip"
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5/plumbing/format/gitignore"
	"github.com/urfave/cli/v3"
)

// NewApp creates a new Enapter CLI tool application instance.
func NewApp() *cli.Command {
	cmd := &cli.Command{}

	cmd.Name = "enapter3"
	cmd.Usage = "Command Line Interface (CLI) for Enapter services."
	cmd.Description = "The Enapter CLI requires an access token for authentication. " +
		"You can obtain the token in your Enapter Cloud account settings."
	cmd.CustomRootCommandHelpTemplate = cli.RootCommandHelpTemplate + enapterAPIEnvVarsHelp

	cli.ShowSubcommandHelp = func(cmd *cli.Command) error {
		tmpl := cmd.CustomHelpTemplate
		if tmpl == "" {
			tmpl = cli.SubcommandHelpTemplate
		}
		cli.HelpPrinter(cmd.Root().Writer, tmpl, cmd)
		return nil
	}

	cmd.Commands = []*cli.Command{
		buildCmdSite(),
		buildCmdDevice(),
		buildCmdBlueprint(),
		buildCmdRuleEngine(),
		buildCmdConnection(),
	}

	return cmd
}

const blueprintIgnoreFile = ".blueprintignore"

func zipDir(dir string) ([]byte, error) {
	fsys := os.DirFS(dir)

	matcher, err := loadBlueprintIgnore(dir)
	if err != nil {
		return nil, fmt.Errorf("load %s: %w", blueprintIgnoreFile, err)
	}

	buf := &bytes.Buffer{}
	zw := zip.NewWriter(buf)

	err = fs.WalkDir(fsys, ".", func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if path == "." {
			return nil
		}

		if matcher.Match(path, entry.IsDir()) {
			if entry.IsDir() {
				return fs.SkipDir
			}
			return nil
		}

		if entry.IsDir() {
			return nil
		}

		f, err := fsys.Open(path)
		if err != nil {
			return fmt.Errorf("open: %w", err)
		}
		defer f.Close()

		zf, err := zw.Create(path)
		if err != nil {
			return fmt.Errorf("create: %w", err)
		}

		if _, err = io.Copy(zf, f); err != nil {
			return fmt.Errorf("copy: %w", err)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("walk dir: %w", err)
	}

	if err := zw.Close(); err != nil {
		return nil, fmt.Errorf("close zip: %w", err)
	}

	return buf.Bytes(), nil
}

type blueprintIgnoreMatcher struct {
	m      gitignore.Matcher
	noFile bool
}

func (m *blueprintIgnoreMatcher) Match(path string, isDir bool) bool {
	if m.noFile {
		return false
	}

	if path == blueprintIgnoreFile {
		return true
	}

	parts := strings.Split(path, string(filepath.Separator))
	return m.m.Match(parts, isDir)
}

func loadBlueprintIgnore(dir string) (*blueprintIgnoreMatcher, error) {
	f, err := os.Open(filepath.Join(dir, blueprintIgnoreFile))
	if os.IsNotExist(err) {
		return &blueprintIgnoreMatcher{noFile: true}, nil
	}
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var patterns []gitignore.Pattern

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		patterns = append(patterns, gitignore.ParsePattern(line, nil))
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return &blueprintIgnoreMatcher{m: gitignore.NewMatcher(patterns)}, nil
}
