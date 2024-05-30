package engine

import (
	"log/slog"
	"os"
	"testing"

	"github.com/blackstork-io/fabric/pkg/diagnostics/diagtest"
	"github.com/blackstork-io/fabric/plugin"
)

func TestEngineFetchData(t *testing.T) {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelWarn,
	})))
	fetchDataTest(
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
		[][]diagtest.Assert{},
	)
	fetchDataTest(
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
		[][]diagtest.Assert{},
	)
}

func TestEngineLint(t *testing.T) {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelWarn,
	})))
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
		[][]diagtest.Assert{
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
		[][]diagtest.Assert{
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
		[][]diagtest.Assert{},
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
		[][]diagtest.Assert{
			{diagtest.IsError, diagtest.SummaryContains("Missing datasource")},
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
		[][]diagtest.Assert{},
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
		[][]diagtest.Assert{
			{diagtest.IsWarning, diagtest.DetailContains("support configuration")},
		},
	)
}

func TestEngineRenderContent(t *testing.T) {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelWarn,
	})))
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
		diagtest.Asserts{},
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
		diagtest.Asserts{},
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
		diagtest.Asserts{},
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
		diagtest.Asserts{},
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
		diagtest.Asserts{},
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
		diagtest.Asserts{
			{diagtest.IsError, diagtest.SummaryContains("Circular reference detected")},
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
		diagtest.Asserts{},
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
		diagtest.Asserts{},
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
		diagtest.Asserts{},
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
		diagtest.Asserts{
			{diagtest.IsWarning, diagtest.SummaryContains("Data conflict")},
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
		diagtest.Asserts{},
	)
	renderTest(
		t, "No fabric files",
		[]string{},
		"test-doc",
		[]string{},
		diagtest.Asserts{
			{diagtest.IsError, diagtest.SummaryContains("No valid fabric files found")},
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
		diagtest.Asserts{},
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
		diagtest.Asserts{},
	)
	renderTest(
		t, "Document meta",
		[]string{
			`
			document "test" {
				meta {
					authors = ["foo", "bar"]
					version = "0.1.2"
					tags = ["xxx", "yyy"]
				}
				content text {
					query = ".document.meta.authors"
					value = <<-EOT
						authors={{ .query_result | join "," }},
						version={{ .document.meta.version }},
						tag={{ index .document.meta.tags 0 }}
					EOT
				}
			}
			`,
		},
		"test",
		[]string{"authors=foo,bar,\nversion=0.1.2,\ntag=xxx"},
		diagtest.Asserts{},
	)
	renderTest(
		t, "Document and content meta",
		[]string{
			`
			document "test" {
				meta {
				  authors = ["foo"]
				}
				section {
				  meta {
					authors = ["bar"]
				  }
				  content text {
					meta {
					  authors = ["baz"]
					}
					query = "(.document.meta.authors[0] + .section.meta.authors[0] + .content.meta.authors[0])" //
					value = "author = {{ .query_result }}"
				  }
				}
			  }
			`,
		},
		"test",
		[]string{"author = foobarbaz"},
		diagtest.Asserts{},
	)
	renderTest(
		t, "Meta scoping and nesting",
		[]string{
			`
			content text get_section_author {
				query = ".section.meta.authors[0] // \"unknown\""
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
							authors = ["foo"]
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
									authors = ["bar"]
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
		diagtest.Asserts{},
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
		diagtest.Asserts{},
	)
	renderTest(
		t, "vars refs query in override",
		[]string{`
			content text "base" {
				vars {
					a = "original"
					b = "unique to base"
				}
				value = "base:\n{{toPrettyJson .vars}}"
			}
			document "test-doc" {
				content ref "base" {
					base = content.text.base
				}
				content ref "ref" {
					base = content.text.base
					vars {
						c = "unique to ref"
						q_b = query_jq(".vars.b") // works as expected
						a = query_jq(".vars.c") // doesn't have access to c, because "a" overrides the variable in the "base", and therefore is executed before c is set
				  	}
				  	value = "ref: {{toPrettyJson .vars}}"
				}
			}
		`},
		"test-doc",
		[]string{
			`base:
{
  "a": "original",
  "b": "unique to base"
}`,
			`ref: {
  "a": null,
  "b": "unique to base",
  "c": "unique to ref",
  "q_b": "unique to base"
}`,
		},
		diagtest.Asserts{},
	)
	renderTest(
		t, "vars refs query in override",
		[]string{`
			content text "base" {
				vars {
					a = "original"
					b = query_jq(".vars.a")
				}
				value = "base:\n{{toPrettyJson .vars}}"
			}
			document "test-doc" {
				content ref "base" {
					base = content.text.base
				}
				content ref "ref" {
					base = content.text.base
					vars {
						a = "redefined"
				  	}
				  	value = "ref: {{toPrettyJson .vars}}"
				}
			}
		`},
		"test-doc",
		[]string{
			`base:
{
  "a": "original",
  "b": "original"
}`,
			`ref: {
  "a": "redefined",
  "b": "redefined"
}`,
		},
		diagtest.Asserts{},
	)
	renderTest(
		t, "vars inheritance",
		[]string{`
		  document "test-doc" {
			vars {
			  docVar = "docVar"
			}
			section "sect" {
			  vars {
				sectVar = "sectVar"
			  }
			  content text {
				vars {
				  contentVar = "contentVar"
				}
				value = "1: {{toPrettyJson .vars}}"
			  }
			  content text {
				value = "2: {{toPrettyJson .vars}}"
			  }
			}
			content text {
			  value = "3: {{toPrettyJson .vars}}"
			}
		  }
		`},
		"test-doc",
		[]string{
			`3: {
  "docVar": "docVar"
}`,
			`1: {
  "contentVar": "contentVar",
  "docVar": "docVar",
  "sectVar": "sectVar"
}`,
			`2: {
  "docVar": "docVar",
  "sectVar": "sectVar"
}`,
		},
		diagtest.Asserts{},
	)
	renderTest(
		t, "vars combined inheritance and shadowing",
		[]string{`
			document "test-doc" {
				vars {
					v1 = 1
					v4 = "not evaluated"
					v2 = query_jq(".vars.v1 + 1")
			  	}
			  	section "sect" {
					vars {
						v7 = "not evaluated"
						v3 = query_jq(".vars.v2 + 1")
						v4 = query_jq(".vars.v3 + 1")
						v5 = query_jq(".vars.v4 + 1")
					}
					content text {
						vars {
							v6 = query_jq(".vars.v5 + 1")
							v7 = query_jq(".vars.v6 + 1")
							v8 = query_jq(".vars.v7 + 1")
						}
						value = "{{toPrettyJson .vars}}"
					}
				}
			}
		`},
		"test-doc",
		[]string{
			`{
  "v1": 1,
  "v2": 2,
  "v3": 3,
  "v4": 4,
  "v5": 5,
  "v6": 6,
  "v7": 7,
  "v8": 8
}`,
		},
		diagtest.Asserts{},
	)
	renderTest(
		t, "vars deep nesting and complex result type",
		[]string{`
			document "test-doc" {
				content text {
					vars {
						a = {
							b = [1, 10, 100]
						}
						c = {
							d = [
								query_jq(<<EOT
									{
										"e": [(.vars.a.b[0] + 1)]
									}
								EOT
								)
							]
						}
					}
					value = "{{toPrettyJson .vars}}"
				}
			}
		`},
		"test-doc",
		[]string{
			`{
  "a": {
    "b": [
      1,
      10,
      100
    ]
  },
  "c": {
    "d": [
      {
        "e": [
          2
        ]
      }
    ]
  }
}`,
		},
		diagtest.Asserts{},
	)
}
