/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/spf13/cobra"
)

// Set by goreleaser
var version = "v0.0.0-dev"

type logLevel struct {
	Name string
	Val  int
}

type logLevels []logLevel

func (ll logLevels) Find(name string) (level int, found bool) {
	name = strings.ToLower(strings.TrimSpace(rawArgs.logLevel))
	idx := slices.IndexFunc(ll, func(s logLevel) bool {
		return s.Name == name
	})
	if idx == -1 {
		return
	}
	return ll[idx].Val, true
}

// TODO: convert to appropriate type for logging framework of choice
var ValidLogLevels = logLevels{
	{"trace", 0},
	{"debug", 1},
	{"info", 2},
	{"warn", 3},
	{"error", 4},
	{"fatal", 5},
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
		if rawArgs.verbose {
			rawArgs.logLevel = "debug"
		}
		logLevel, found := ValidLogLevels.Find(rawArgs.logLevel)
		if !found {
			return fmt.Errorf("invalid log level '%s'", rawArgs.logLevel)
		}
		// TODO: set log output format
		_ = logLevel
		output := strings.ToLower(strings.TrimSpace(rawArgs.logOutput))
		switch output {
		case "json", "plain":
			// TODO: set log output format
			_ = output
		default:
			return fmt.Errorf("invalid log output '%s'", rawArgs.logOutput)
		}
		SourceDir, err = filepath.Abs(rawArgs.sourceDir)
		if err != nil {
			return fmt.Errorf("bad source dir '%s': %w", rawArgs.sourceDir, err)
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
	rootCmd.PersistentFlags().StringVar(&rawArgs.logLevel, "logging-level", "info", "logging level")
	rootCmd.PersistentFlags().BoolVarP(&rawArgs.verbose, "verbose", "v", false, "a shortcut to --logging-level debug")
}
