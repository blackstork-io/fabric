package e2e_test

import (
	"fmt"
	"log/slog"
	"os"
	"testing"
	"testing/fstest"

	"github.com/blackstork-io/fabric/cmd"
	"github.com/blackstork-io/fabric/pkg/diagnostics/diagtest"
	"github.com/blackstork-io/fabric/pkg/fabctx"
)

func lintTest(t *testing.T, fullLint bool, testName string, files []string, diagAsserts diagtest.Asserts) {
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
		ctx := fabctx.New(fabctx.NoSignals)
		ctx = fabctx.WithLinting(ctx)

		eval := cmd.NewEvaluator()
		defer func() {
			eval.Cleanup(nil)
		}()

		diags := cmd.Lint(ctx, eval, sourceDir, fullLint)

		diagAsserts.AssertMatch(t, diags, eval.FileMap)
	})
}

func fullLintTest(t *testing.T, testName string, files []string, diagAsserts diagtest.Asserts) {
	t.Helper()
	lintTest(t, true, testName, files, diagAsserts)
}

func limitedLintTest(t *testing.T, testName string, files []string, diagAsserts diagtest.Asserts) {
	t.Helper()
	lintTest(t, false, testName, files, diagAsserts)
}

func TestE2ELint(t *testing.T) {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelWarn,
	})))

	t.Parallel()
	fullLintTest(
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
		diagtest.Asserts{
			{diagtest.IsError, diagtest.SummaryContains("Circular reference detected")},
		},
	)
	fullLintTest(
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
		diagtest.Asserts{
			{diagtest.IsWarning, diagtest.SummaryContains("Data conflict")},
		},
	)

	limitedLintTest(
		t, "Unknown plugins are fine in limited lint",
		[]string{
			`
			document "doc1" {
				data made_up_data_source "name1" {

				}
			}
			document "doc2" {
				content made_up_content_provider "name2" {

				}
			}
			`,
		},
		diagtest.Asserts{},
	)
	fullLintTest(
		t, "Unknown plugins generate diags in full lint",
		[]string{
			`
			document "doc1" {
				data made_up_data_source "name1" {

				}
			}
			document "doc2" {
				content made_up_content_provider "name2" {

				}
			}
			`,
		},
		diagtest.Asserts{
			{diagtest.IsError, diagtest.SummaryContains("Missing data source")},
			{diagtest.IsError, diagtest.SummaryContains("Missing content provider")},
		},
	)
	limitedLintTest(
		t, "Unknown config is fine in limited lint",
		[]string{
			`
			document "doc1" {
				data inline "name1" {
					config {}
				}
			}
			`,
		},
		diagtest.Asserts{},
	)

	fullLintTest(
		t, "Unknown config generate diags in full lint",
		[]string{
			`
			document "doc1" {
				data inline "name1" {
					config {}
				}
			}
			`,
		},
		diagtest.Asserts{
			{diagtest.IsWarning, diagtest.DetailContains("support configuration")},
		},
	)
}
