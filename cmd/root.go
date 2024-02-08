package cmd

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"slices"
	"strings"

	"github.com/spf13/cobra"

	"github.com/blackstork-io/fabric/pkg/utils"
)

const devVersion = "v0.0.0-dev"

// Overriden by goreleaser
var version = devVersion

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

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "fabric",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Version: version,

	PersistentPreRunE: func(_ *cobra.Command, _ []string) (err error) {
		err = validateDir("source dir", rawArgs.sourceDir)
		if err != nil {
			return
		}
		cliArgs.sourceDir = rawArgs.sourceDir

		// TODO: make optional after #5 is implemented
		err = validateDir("plugins dir", rawArgs.pluginsDir)
		if err != nil {
			return
		}
		cliArgs.pluginsDir = rawArgs.pluginsDir

		var level slog.Level
		if rawArgs.verbose {
			level = slog.LevelDebug
		} else {
			level, err = validLogLevels.Find(rawArgs.logLevel)
			if err != nil {
				return
			}
		}
		opts := &slog.HandlerOptions{
			Level: level,
			// add source if in debug mode
			AddSource: level == slog.LevelDebug,
		}

		var logger *slog.Logger

		switch strings.ToLower(strings.TrimSpace(rawArgs.logOutput)) {
		case "plain":
			logger = slog.New(slog.NewTextHandler(os.Stderr, opts))
		case "json":
			logger = slog.New(slog.NewJSONHandler(os.Stderr, opts))
		default:
			return fmt.Errorf("unknown log output '%s'", rawArgs.logOutput)
		}
		slog.SetDefault(logger)
		slog.Debug("Logging enabled")
		if version == devVersion {
			slog.Warn("This is a dev version of the software")
		}
		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var cliArgs = struct {
	sourceDir  string
	pluginsDir string
}{}

var rawArgs = struct {
	sourceDir  string
	logOutput  string
	logLevel   string
	pluginsDir string
	verbose    bool
}{}

func init() {
	rootCmd.PersistentFlags().StringVar(&rawArgs.sourceDir, "source-dir", ".", "a path to a directory with *.fabric files")
	rootCmd.PersistentFlags().StringVar(&rawArgs.logOutput, "log-format", "plain", "format of the logs (plain or json)")
	rootCmd.PersistentFlags().StringVar(
		&rawArgs.logLevel, "log-level", "info",
		fmt.Sprintf("logging level (%s)", validLogLevels.String()),
	)
	rootCmd.PersistentFlags().BoolVarP(&rawArgs.verbose, "verbose", "v", false, "a shortcut to --log-level debug")
	// TODO: after #5 is implemented - make optional
	rootCmd.PersistentFlags().StringVar(
		&rawArgs.pluginsDir, "plugins-dir", "", "override for plugins dir from fabric configuration (required)",
	)
	err := rootCmd.MarkPersistentFlagRequired("plugins-dir")
	if err != nil {
		panic(err)
	}
}
