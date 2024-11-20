package app

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/pkg/errors"
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

func validateDuration(cli *cli.Context, name string, dflt, min time.Duration) (*time.Duration, error) {
	if !cli.IsSet(name) {
		return &dflt, nil
	}
	p := cli.String(name)
	period, err := time.ParseDuration(p)
	if err != nil {
		return nil, err
	}
	if period < min {
		return nil, errors.New(fmt.Sprintf("%s should be greater than %s", name, min))
	}
	return &period, nil
}
