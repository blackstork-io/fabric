package e2e_test

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/assert"

	"github.com/blackstork-io/fabric/cmd"
	"github.com/blackstork-io/fabric/pkg/diagnostics"
	"github.com/blackstork-io/fabric/pkg/fabctx"
	"github.com/blackstork-io/fabric/test/e2e/diag_test"
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
		eval := cmd.NewEvaluator()
		defer func() {
			eval.Cleanup(nil)
		}()

		var res string
		diags := eval.ParseFabricFiles(sourceDir)
		ctx := fabctx.New(fabctx.NoSignals)
		if !diags.HasErrors() {
			if !diags.Extend(eval.LoadPluginResolver(false)) && !diags.Extend(eval.LoadPluginRunner(ctx)) {
				var diag diagnostics.Diag
				res, diag = cmd.Render(ctx, eval.Blocks, eval.PluginCaller(), docName)
				diags.Extend(diag)
			}
		}

		if len(expectedResult) == 0 {
			// so nil == []string{}
			assert.Empty(t, res)
		} else {
			assert.EqualValues(
				t,
				strings.Join(expectedResult, "\n\n"),
				res,
			)
		}
		diag_test.CompareDiags(t, diags, diagAsserts)
	})
}

func TestE2ERender(t *testing.T) {
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
					value = "Hello from fabric"
				}
			}

			document "goodbye" {
				title = "Goodbye"
				content text {
					value = "Goodbye from fabric"
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
				value = "Hello from ref"
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
				value = "Hello from ref"
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

			content blockquote "actual_block" {
				value = "Hello from ref chain"
			}
			`,
			`
			content ref "add_format_as" {
				base = content.blockquote.actual_block
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
				value = "Near refloop"
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
					value = "final section"
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
					value = "s1"
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
						value = "s2"
					}
					section {
						title = "sect3"
						content text {
							value = "s3"
						}
						content ref {
							base = content.text.some_text
						}
						content text {
							value = "s3 extra"
						}
						section ref {
							base = section.sect4
						}
					}
				}
			}

			content text "some_text" {
				value = "some_text"
			}

			`,
		},
		"test-doc",
		[]string{
			"## sect1",
			"s1",
			"some_text",
			"### sect2",
			"some_text",
			"s2",
			"#### sect3",
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
					value = "${2+2}"
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
		t, "Data ref name warning missing",
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
				data ref {
					base = data.inline.name
				}
			}
			`,
		},
		"test-doc",
		[]string{},
		[][]diag_test.Assert{
			{diag_test.IsWarning, diag_test.SummaryContains("Data conflict")},
		},
	)
	renderTest(
		t, "Content ref name no-error",
		[]string{
			`
			content text "name" {
				value = "txt"
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
	renderTest(
		t, "Data block result access",
		[]string{
			`
			document "test-doc" {
				data inline "name" {
					attr = "val"
				}
				content text {
					value = "From data block: {{.data.inline.name.attr}}"
				}
			}
			`,
		},
		"test-doc",
		[]string{"From data block: val"},
		[][]diag_test.Assert{},
	)
	renderTest(
		t, "Data with jq interaction",
		[]string{
			`
			document "test" {
				data inline "foo" {
				  items = ["a", "b", "c"]
				  x = 1
				  y = 2
				}
				content text {
				  query = ".data.inline.foo.items | length"
				  value = "There are {{ .query_result }} items"
				}
			}
			`,
		},
		"test",
		[]string{"There are 3 items"},
		[][]diag_test.Assert{},
	)
	renderTest(
		t, "Document meta",
		[]string{
			`
			document "test" {
				meta {
					author = "foo"
				}
				content text {
				  query = ".document.meta.author"
				  value = "author = {{ .query_result }}"
				}
			}
			`,
		},
		"test",
		[]string{"author = foo"},
		[][]diag_test.Assert{},
	)
	renderTest(
		t, "Document and content meta",
		[]string{
			`
			document "test" {
				meta {
				  author = "foo"
				}
				section {
				  meta {
					author = "bar"
				  }
				  content text {
					meta {
					  author = "baz"
					}
					query = "(.document.meta.author + .section.meta.author + .content.meta.author)" //
					value = "author = {{ .query_result }}"
				  }
				}
			  }
			`,
		},
		"test",
		[]string{"author = foobarbaz"},
		[][]diag_test.Assert{},
	)
	renderTest(
		t, "Meta scoping and nesting",
		[]string{
			`
			content text get_section_author {
				query = ".section.meta.author // \"unknown\""
				value = "author = {{ .query_result }}"
			}
			document "test" {
				content ref {
					base = content.text.get_section_author
				}
				section {
					content ref {
						base = content.text.get_section_author
					}
					section {
						meta {
							author = "foo"
						}
						content ref {
							base = content.text.get_section_author
						}
						section {
							content ref {
								base = content.text.get_section_author
							}
							section {
								meta {
									author = "bar"
								}
								content ref {
									base = content.text.get_section_author
								}
							}
						}
					}
				}
			}
			`,
		},
		"test",
		[]string{
			"author = unknown",
			"author = unknown",
			"author = foo",
			"author = unknown",
			"author = bar",
		},
		[][]diag_test.Assert{},
	)
	renderTest(
		t, "Reference rendered blocks",
		[]string{
			`
			document "test" {
				content text {
					value = "first result"
				  }
				  content text {
					query = ".document.content.children[0].markdown"
					value = "content[0] = {{ .query_result }}"
				  }
				  content text {
					query = ".document.content.children[1].markdown"
					value = "content[1] = {{ .query_result }}"
				  }
			  }
			`,
		},
		"test",
		[]string{
			"first result",
			"content[0] = first result",
			"content[1] = content[0] = first result",
		},
		[][]diag_test.Assert{},
	)
}
