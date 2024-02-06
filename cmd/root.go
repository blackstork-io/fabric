package cmd

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/spf13/cobra"

	"github.com/blackstork-io/fabric/pkg/utils"
)

// Set by goreleaser
var version = "v0.0.0-dev"

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

	PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
		SourceDir, err = filepath.Abs(rawArgs.sourceDir)
		if err != nil {
			return fmt.Errorf("bad source dir '%s': %w", rawArgs.sourceDir, err)
		}

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

// exposed args
var (
	SourceDir string
)

var rawArgs = struct {
	sourceDir string
	logOutput string
	logLevel  string
	verbose   bool
}{}

func init() {
	rootCmd.PersistentFlags().StringVar(&rawArgs.sourceDir, "source-dir", ".", "a path to a directory with *.fabric files")
	rootCmd.PersistentFlags().StringVar(&rawArgs.logOutput, "log-output", "plain", "log output kind (plain or json)")
	rootCmd.PersistentFlags().StringVar(
		&rawArgs.logLevel, "logging-level", "info",
		fmt.Sprintf("logging level (%s)", validLogLevels.String()),
	)
	rootCmd.PersistentFlags().BoolVarP(&rawArgs.verbose, "verbose", "v", false, "a shortcut to --logging-level debug")
}
