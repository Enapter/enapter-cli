package configfile

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

type Config struct {
	DefaultConn string                `json:"default_connection,omitempty"`
	Connections map[string]Connection `json:"connections,omitempty"`
}

type Connection struct {
	Gateway       bool   `json:"gateway,omitempty"`
	URL           string `json:"url"`
	SiteID        string `json:"site_id,omitempty"`
	Token         Token  `json:"token"`
	AllowInsecure bool   `json:"allow_insecure,omitempty"`
}

type Token struct {
	Value string `json:"value"`
}

const (
	dirName  = ".enapter3"
	fileName = "config.json"
)

func Load() (Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return Config{}, fmt.Errorf("get home dir: %w", err)
	}

	path := filepath.Join(home, dirName, fileName)
	f, err := os.Open(path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return Config{}, nil
		}
		return Config{}, fmt.Errorf("open config file: %w", err)
	}
	defer f.Close()

	var config Config
	if err := json.NewDecoder(f).Decode(&config); err != nil {
		return Config{}, fmt.Errorf("decode config file: %w", err)
	}

	return config, nil
}

func Save(c Config) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("get home dir: %w", err)
	}

	const perm = 0o755
	dir := filepath.Join(home, dirName)
	if err := os.MkdirAll(dir, perm); err != nil {
		return fmt.Errorf("create config dir: %w", err)
	}

	path := filepath.Join(home, dirName, fileName)
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create config file: %w", err)
	}
	defer f.Close()

	encoder := json.NewEncoder(f)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(c); err != nil {
		return fmt.Errorf("encode config file: %w", err)
	}

	return f.Sync()
}
