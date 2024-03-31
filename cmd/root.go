package cmd

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"runtime/debug"
	"slices"
	"strings"

	"github.com/golang-cz/devslog"
	"github.com/lmittmann/tint"
	"github.com/mattn/go-colorable"
	"github.com/spf13/cobra"
	"golang.org/x/term"

	"github.com/blackstork-io/fabric/pkg/fabctx"
	"github.com/blackstork-io/fabric/pkg/utils"
)

// Overridden by goreleaser.
var (
	version = ""
	builtBy = "golang"
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

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "fabric",
	Version: version,

	PersistentPreRunE: func(_ *cobra.Command, _ []string) (err error) {
		err = validateDir("source dir", rawArgs.sourceDir)
		if err != nil {
			return
		}
		cliArgs.sourceDir = rawArgs.sourceDir

		cliArgs.colorize = rawArgs.colorize && term.IsTerminal(int(os.Stderr.Fd()))

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

		var handler slog.Handler

		switch strings.ToLower(strings.TrimSpace(rawArgs.logOutput)) {
		case "plain":
			if cliArgs.colorize && level <= slog.LevelDebug {
				handler = devslog.NewHandler(os.Stderr, &devslog.Options{
					HandlerOptions: opts,
				})
			} else {
				var output io.Writer
				if cliArgs.colorize {
					// only affects windows, noop on *nix
					output = colorable.NewColorable(os.Stderr)
				} else {
					output = os.Stderr
				}

				handler = tint.NewHandler(
					output,
					&tint.Options{
						AddSource:   opts.AddSource,
						Level:       opts.Level,
						ReplaceAttr: opts.ReplaceAttr,
						NoColor:     !cliArgs.colorize,
					},
				)
			}
		case "json":
			handler = slog.NewJSONHandler(os.Stderr, opts)
		default:
			return fmt.Errorf("unknown log output '%s'", rawArgs.logOutput)
		}
		slog.SetDefault(slog.New(handler))
		slog.Debug("Logging enabled")
		if strings.Contains(version, "-dev") {
			slog.Warn("This is a dev version of the software!", "version", version)
		}
		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.ExecuteContext(fabctx.New())
	if err != nil {
		os.Exit(1)
	}
}

var cliArgs = struct {
	sourceDir string
	colorize  bool
}{}

var rawArgs = struct {
	sourceDir string
	logOutput string
	logLevel  string
	verbose   bool
	colorize  bool
}{}

func init() {
	rootCmd.PersistentFlags().StringVar(&rawArgs.sourceDir, "source-dir", ".", "a path to a directory with *.fabric files")
	rootCmd.PersistentFlags().StringVar(&rawArgs.logOutput, "log-format", "plain", "format of the logs (plain or json)")
	rootCmd.PersistentFlags().StringVar(
		&rawArgs.logLevel, "log-level", "info",
		fmt.Sprintf("logging level (%s)", validLogLevels.String()),
	)
	rootCmd.PersistentFlags().BoolVar(&rawArgs.colorize, "color", true, "enables colorizing the logs and diagnostics (if supported by the terminal and log format)")
	rootCmd.PersistentFlags().BoolVarP(&rawArgs.verbose, "verbose", "v", false, "a shortcut to --log-level debug")

	if builtBy == "goreleaser" {
		// version is set by goreleaser
		return
	}

	version = fmt.Sprintf(
		"%s+builtBy.%s",
		versionFromBuildInfo(),
		builtBy,
	)
}

func versionFromBuildInfo() (version string) {
	version = "v0.0.0-dev"
	info, ok := debug.ReadBuildInfo()
	if !ok || info == nil {
		return
	}
	if info.Main.Version != "(devel)" {
		version = info.Main.Version
		if !strings.HasPrefix(version, "v") {
			version = "v" + version
		}
		return
	}
	// It's a dev version not built by goreleaser, add extra info
	dirtyIdx := slices.IndexFunc(info.Settings, func(s debug.BuildSetting) bool {
		return s.Key == "vcs.modified"
	})
	if dirtyIdx != -1 && info.Settings[dirtyIdx].Value == "true" {
		version += "+dirty"
	}

	shaIdx := slices.IndexFunc(info.Settings, func(s debug.BuildSetting) bool {
		return s.Key == "vcs.revision"
	})
	if shaIdx != -1 {
		version = fmt.Sprintf("%s+rev.%s", version, info.Settings[shaIdx].Value)
	}
	return
}
