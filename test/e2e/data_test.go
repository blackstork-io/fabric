package e2e_test

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/assert"

	"github.com/blackstork-io/fabric/cmd"
	"github.com/blackstork-io/fabric/pkg/diagnostics"
	"github.com/blackstork-io/fabric/plugin"
	"github.com/blackstork-io/fabric/test/e2e/diag_test"
)

func dataTest(t *testing.T, testName string, files []string, target string, expectedResult plugin.MapData, diagAsserts [][]diag_test.Assert) {
	t.Helper()
	t.Run(testName, func(t *testing.T) {
		t.Parallel()
		t.Helper()

		sourceDir := fstest.MapFS{}
		for i, content := range files {
			sourceDir[fmt.Sprintf("file_%d.fabric", i)] = &fstest.MapFile{
				Data: []byte(content),
				Mode: 0o777,
			}
		}
		eval := cmd.NewEvaluator("")
		defer func() {
			eval.Cleanup(nil)
		}()

		var res plugin.MapData
		diags := eval.ParseFabricFiles(sourceDir)
		if !diags.HasErrors() {
			if !diags.Extend(eval.LoadRunner()) {
				var diag diagnostics.Diag
				res, diag = cmd.Data(context.Background(), eval.Blocks, eval.PluginCaller(), target)
				diags.Extend(diag)
			}
		}

		assert.Equal(t, expectedResult, res)
		if !diag_test.MatchBiject(diags, diagAsserts) {
			assert.Fail(t, "Diagnostics do not match", diags)
		}
	})
}

func TestE2EData(t *testing.T) {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelWarn,
	})))

	t.Parallel()
	dataTest(
		t, "Basic",
		[]string{
			`
			document "hello" {
				data inline "test" {
					hello = "world"
				}

				content text {
					text = "hello"
				}
			}
			`,
		},
		"document.hello.data.inline.test",
		plugin.MapData{
			"hello": plugin.StringData("world"),
		},
		[][]diag_test.Assert{},
	)
	dataTest(
		t, "Basic",
		[]string{
			`
			data inline "test" {
				hello = "world"
			}
			document "hello" {
				content text {
					text = "hello"
				}
			}
			`,
		},
		"data.inline.test",
		plugin.MapData{
			"hello": plugin.StringData("world"),
		},
		[][]diag_test.Assert{},
	)
}
