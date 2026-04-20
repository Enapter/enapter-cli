// Package blueprint provides functions for packaging Enapter blueprints
// into zip archives, with support for .blueprintignore files.
package blueprint

import (
	"archive/zip"
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"strings"

	"github.com/go-git/go-git/v5/plumbing/format/gitignore"
)

const ignoreFileName = ".blueprintignore"

// Zip creates a zip archive from the blueprint directory at the given
// filesystem root. It respects .blueprintignore patterns (gitignore syntax).
// The .blueprintignore file itself is excluded from the archive.
func Zip(fsys fs.FS) ([]byte, error) {
	matcher, err := loadIgnoreMatcher(fsys)
	if err != nil {
		return nil, fmt.Errorf("load %s: %w", ignoreFileName, err)
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

type ignoreMatcher struct {
	m      gitignore.Matcher
	noFile bool
}

func (m *ignoreMatcher) Match(path string, isDir bool) bool {
	if m.noFile {
		return false
	}

	if path == ignoreFileName {
		return true
	}

	parts := strings.Split(path, "/")
	return m.m.Match(parts, isDir)
}

func loadIgnoreMatcher(fsys fs.FS) (*ignoreMatcher, error) {
	f, err := fsys.Open(ignoreFileName)
	if errors.Is(err, fs.ErrNotExist) {
		return &ignoreMatcher{noFile: true}, nil
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

	return &ignoreMatcher{m: gitignore.NewMatcher(patterns)}, nil
}
