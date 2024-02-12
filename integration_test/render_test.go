//go:build integration

package integration_test

import (
	"fmt"
	"log/slog"
	"os"
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/assert"

	"github.com/blackstork-io/fabric/cmd"
	"github.com/blackstork-io/fabric/integration_test/diag_test"
)

func renderTest(t *testing.T, testName string, files []string, docName string, expectedResult []string, diagAsserts [][]diag_test.Assert) {
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

		res, _, diags := cmd.Render(pluginsDir, sourceDir, docName)
		if len(expectedResult) == 0 {
			// so nil == []string{}
			assert.Empty(t, res)
		} else {
			assert.EqualValues(
				t,
				expectedResult,
				res,
			)
		}
		if !diag_test.MatchBiject(diags, diagAsserts) {
			assert.Fail(t, "Diagnostics do not match", diags)
		}
	})
}

func TestRender(t *testing.T) {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelWarn,
	})))

	t.Parallel()
	renderTest(
		t, "Basic",
		[]string{
			`
			document "hello" {
				title = "Welcome"
				content text {
					text = "Hello from fabric"
				}
			}

			document "goodbye" {
				title = "Goodbye"
				content text {
					text = "Goodbye from fabric"
				}
			}
			`,
		},
		"hello",
		[]string{
			"# Welcome",
			"Hello from fabric",
		},
		[][]diag_test.Assert{},
	)
	renderTest(
		t, "Ref",
		[]string{
			`
			document "test-doc" {
				title = "Welcome"
				content ref {
					base = content.text.external_block
				}
			}

			content text "external_block" {
				text = "Hello from ref"
			}
			`,
		},
		"test-doc",
		[]string{
			"# Welcome",
			"Hello from ref",
		},
		[][]diag_test.Assert{},
	)
	renderTest(
		t, "Ref across files",
		[]string{
			`
			document "test-doc" {
				title = "Welcome"
				content ref {
					base = content.text.external_block
				}
			}
			`,
			`
			content text "external_block" {
				text = "Hello from ref"
			}
			`,
		},
		"test-doc",
		[]string{
			"# Welcome",
			"Hello from ref",
		},
		[][]diag_test.Assert{},
	)
	renderTest(
		t, "Ref chain",
		[]string{
			`
			document "test-doc" {
				content ref {
					base = content.ref.add_format_as
				}
			}

			content text "actual_block" {
				text = "Hello from ref chain"
			}
			`,
			`
			content ref "add_format_as" {
				base = content.text.actual_block
				format_as = "blockquote"
			}
			`,
		},
		"test-doc",
		[]string{
			"> Hello from ref chain",
		},
		[][]diag_test.Assert{},
	)
	renderTest(
		t, "Ref loop untouched",
		[]string{
			`
			document "test-doc" {
				content ref {
					base = content.text.actual_block
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

			content text "actual_block" {
				text = "Near refloop"
			}
			`,
		},
		"test-doc",
		[]string{
			"Near refloop",
		},
		[][]diag_test.Assert{},
	)

	renderTest(
		t, "Ref loop touched",
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
		"test-doc",
		[]string{},
		[][]diag_test.Assert{
			{diag_test.IsError, diag_test.SummaryContains("Circular reference detected")},
		},
	)
	renderTest(
		t, "Sections",
		[]string{
			`
			section "sect4" {
				content text {
					text = "final section"
				}
			}
			document "test-doc" {
				section ref {
					base = section.sect1
				}
			}
			`,
			`
			section "sect1" {
				title = "sect1"
				content text {
					text = "s1"
				}
				content ref {
					base = content.text.some_text
				}
				section {
					title = "sect2"
					content ref {
						base = content.text.some_text
					}
					content text {
						text = "s2"
					}
					section {
						title = "sect3"
						content text {
							text = "s3"
						}
						content ref {
							base = content.text.some_text
						}
						content text {
							text = "s3 extra"
						}
						section ref {
							base = section.sect4
						}
					}
				}
			}

			content text "some_text" {
				text = "some_text"
			}

			`,
		},
		"test-doc",
		[]string{
			"# sect1",
			"s1",
			"some_text",
			"# sect2",
			"some_text",
			"s2",
			"# sect3",
			"s3",
			"some_text",
			"s3 extra",
			"final section",
		},
		[][]diag_test.Assert{},
	)
	renderTest(
		t, "Templating support",
		[]string{
			`
			document "test-doc" {
				content text {
					text = "${2+2}"
				}
			}
			`,
		},
		"test-doc",
		[]string{
			"4",
		},
		[][]diag_test.Assert{},
	)
	renderTest(
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
			}
			`,
		},
		"test-doc",
		[]string{},
		[][]diag_test.Assert{
			{diag_test.IsWarning, diag_test.SummaryContains("Potential data conflict")},
		},
	)
	renderTest(
		t, "Content ref name no-error",
		[]string{
			`
			content text "name" {
				text = "txt"
			}
			document "test-doc" {
				content ref {
					base = content.text.name
				}
			}
			`,
		},
		"test-doc",
		[]string{"txt"},
		[][]diag_test.Assert{},
	)
	renderTest(
		t, "No fabric files",
		[]string{},
		"test-doc",
		[]string{},
		[][]diag_test.Assert{
			{diag_test.IsError, diag_test.SummaryContains("No fabric files found")},
		},
	)
}
