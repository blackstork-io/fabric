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
	"github.com/blackstork-io/fabric/test/e2e/diag_test"
)

func lintTest(t *testing.T, testName string, files []string, diagAsserts [][]diag_test.Assert) {
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

		diags := eval.ParseFabricFiles(sourceDir)
		ctx := context.Background()
		if !diags.HasErrors() {
			if !diags.Extend(eval.LoadPluginResolver(false)) && !diags.Extend(eval.LoadPluginRunner(ctx)) {
				diag := cmd.Lint(ctx, eval.Blocks, eval.PluginCaller())
				diags.Extend(diag)
			}
		}

		if !diag_test.MatchBiject(diags, diagAsserts) {
			assert.Fail(t, "Diagnostics do not match", diags)
		}
	})
}

func TestE2ELint(t *testing.T) {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelWarn,
	})))

	t.Parallel()
	lintTest(
		t, "Ref loop",
		[]string{
			`
			document "test-doc" {
				content ref {
					base = content.ref.actual_block
				}
			}
			`,
			`
			content ref "ref_a" {
				base = content.ref.ref_b
			}

			content ref "ref_b" {
				base = content.ref.ref_c
			}

			content ref "ref_c" {
				base = content.ref.ref_a
			}

			content ref "actual_block" {
				base = content.ref.ref_a
			}
			`,
		},
		[][]diag_test.Assert{
			{diag_test.IsError, diag_test.SummaryContains("Circular reference detected")},
		},
	)
	lintTest(
		t, "Data ref name warning",
		[]string{
			`
			data inline "name" {
				inline {
					a = "1"
				}
			}
			document "test-doc" {
				data ref {
					base = data.inline.name
				}
				data ref {
					base = data.inline.name
				}
			}
			`,
		},
		[][]diag_test.Assert{
			{diag_test.IsWarning, diag_test.SummaryContains("Data conflict")},
		},
	)
}
