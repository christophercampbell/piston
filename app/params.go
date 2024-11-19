package app

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/urfave/cli/v2"
)

func orDefaultString(cli *cli.Context, name, defaultValue string) string {
	value := cli.String(name)
	if value == "" {
		return defaultValue
	}
	return value
}

func processPath(path string) string {
	if userHome, err := os.UserHomeDir(); err == nil {
		path = strings.Replace(path, "~", userHome, 1)
		path = filepath.Clean(path)
	}
	return path
}
