package cmd

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"slices"
	"strings"

	"github.com/blackstork-io/fabric/pkg/utils"
)

type logLevels struct {
	Names []string
	Vals  []slog.Level
}

func (ll *logLevels) Find(name string) (level slog.Level, err error) {
	nameKey := strings.ToLower(strings.TrimSpace(rawArgs.logLevel))
	idx := slices.Index(ll.Names, nameKey)
	if idx == -1 {
		err = fmt.Errorf("unknown log level '%s'", name)
		return
	}
	return ll.Vals[idx], nil
}

func (ll *logLevels) String() string {
	return utils.JoinSurround(", ", "'", ll.Names...)
}

var validLogLevels = logLevels{
	Names: []string{
		"debug",
		"info",
		"warn",
		"error",
	},
	Vals: []slog.Level{
		slog.LevelDebug,
		slog.LevelInfo,
		slog.LevelWarn,
		slog.LevelError,
	},
}

func validateDir(what, dir string) error {
	info, err := os.Stat(dir)
	switch {
	case err == nil:
	case errors.Is(err, os.ErrNotExist):
		return fmt.Errorf("failed to open %s: path '%s' doesn't exist", what, dir)
	case errors.Is(err, os.ErrPermission):
		return fmt.Errorf("failed to open %s: permission to access path '%s' denied", what, dir)
	default:
		return fmt.Errorf("failed to open %s: path '%s': %w", what, dir, err)
	}
	if !info.IsDir() {
		return fmt.Errorf("failed to open %s: path '%s' is not a directory", what, dir)
	}
	return nil
}
