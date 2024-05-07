package e2e_test

import (
	"fmt"
	"log/slog"
	"os"
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/assert"

	"github.com/blackstork-io/fabric/cmd"
	"github.com/blackstork-io/fabric/pkg/diagnostics"
	"github.com/blackstork-io/fabric/pkg/diagnostics/diagtest"
	"github.com/blackstork-io/fabric/pkg/fabctx"
	"github.com/blackstork-io/fabric/plugin"
)

func dataTest(t *testing.T, testName string, files []string, target string, expectedResult plugin.MapData, diagAsserts diagtest.Asserts) {
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
		eval := cmd.NewEvaluator()
		defer func() {
			eval.Cleanup(nil)
		}()

		var res plugin.Data
		diags := eval.ParseFabricFiles(sourceDir)
		ctx := fabctx.New(fabctx.NoSignals)
		if !diags.HasErrors() {
			if !diags.Extend(eval.LoadPluginResolver(false)) && !diags.Extend(eval.LoadPluginRunner(ctx)) {
				var diag diagnostics.Diag
				res, diag = cmd.Data(ctx, eval.Blocks, eval.PluginCaller(), target)
				diags.Extend(diag)
			}
		}

		assert.Equal(t, expectedResult, res)
		diagAsserts.AssertMatch(t, diags, eval.FileMap)
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
		diagtest.Asserts{},
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
		diagtest.Asserts{},
	)
}
