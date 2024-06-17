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
				data json "test" {
					path = "testdata/a.json"
				}

				content text {
					text = "hello"
				}
			}
			`,
		},
		"document.hello.data.json.test",
		plugin.MapData{
			"property_for": plugin.StringData("a.json"),
		},
		[][]diagtest.Assert{},
	)
	fetchDataTest(
		t, "Basic",
		[]string{
			`
			data json "test" {
				path = "testdata/a.json"
			}
			document "hello" {
				content text {
					text = "hello"
				}
			}
			`,
		},
		"data.json.test",
		plugin.MapData{
			"property_for": plugin.StringData("a.json"),
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
			data json "name" {
				path = "testdata/a.json"
			}
			document "test-doc" {
				data ref {
					base = data.json.name
				}
				data ref {
					base = data.json.name
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
				data json "name1" {
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
				data json "name1" {
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

func TestEnvPrefix(t *testing.T) {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelWarn,
	})))
	t.Setenv("OTHER_VAR", "OTHER_VAR")
	t.Setenv("FABRIC_VAR", "FABRIC_VAR")
	t.Setenv("FABRIC_TEST_VAR", "FABRIC_TEST_VAR")

	renderTest(
		t, "Default",
		[]string{
			`
			document "test-doc" {
				content text {
					value = "{{.env.OTHER_VAR}}\n{{.env.FABRIC_VAR}}\n{{.env.FABRIC_TEST_VAR}}"
				}
			}
			`,
		},
		"test-doc",
		[]string{
			"<no value>\nFABRIC_VAR\nFABRIC_TEST_VAR",
		},
		diagtest.Asserts{},
	)
	renderTest(
		t, "Custom",
		[]string{
			`
			fabric {
				expose_env_vars_with_pattern = "FABRIC_TEST_*"
			}
			document "test-doc" {
				content text {
					value = "{{.env.OTHER_VAR}}\n{{.env.FABRIC_VAR}}\n{{.env.FABRIC_TEST_VAR}}"
				}
			}
			`,
		},
		"test-doc",
		[]string{
			"<no value>\n<no value>\nFABRIC_TEST_VAR",
		},
		diagtest.Asserts{},
	)
	renderTest(
		t, "Empty",
		[]string{
			`
			fabric {
				expose_env_vars_with_pattern = ""
			}
			document "test-doc" {
				content text {
					value = "{{.env.OTHER_VAR}}\n{{.env.FABRIC_VAR}}\n{{.env.FABRIC_TEST_VAR}}"
				}
			}
			`,
		},
		"test-doc",
		[]string{
			"<no value>\n<no value>\n<no value>",
		},
		diagtest.Asserts{},
	)
	renderTest(
		t, "Empty",
		[]string{
			`
			fabric {
				expose_env_vars_with_pattern =  "\t FABRIC_TEST_*   "
			}
			document "test-doc" {
				content text {
					value = "{{.env.OTHER_VAR}}\n{{.env.FABRIC_VAR}}\n{{.env.FABRIC_TEST_VAR}}"
				}
			}
			`,
		},
		"test-doc",
		[]string{
			"<no value>\n<no value>\nFABRIC_TEST_VAR",
		},
		diagtest.Asserts{{
			diagtest.IsWarning,
			diagtest.SummaryContains("contains whitespace"),
		}},
	)
	renderTest(
		t, "ErrorInPattern",
		[]string{
			`
			fabric {
				expose_env_vars_with_pattern =  "FABRIC_TEST_["
			}
			document "test-doc" {
				content text {
					value = "{{.env.OTHER_VAR}}\n{{.env.FABRIC_VAR}}\n{{.env.FABRIC_TEST_VAR}}"
				}
			}
			`,
		},
		"test-doc",
		[]string{
			"<no value>\n<no value>\n<no value>",
		},
		diagtest.Asserts{{
			diagtest.IsWarning,
			diagtest.SummaryContains("Failed to parse", "expose_env_vars_with_pattern"),
		}},
	)
	// TODO: uncomment this test when switched away from gocty parsing for global config
	// renderTest(
	// 	t, "Null",
	// 	[]string{
	// 		`
	// 		fabric {
	// 			expose_env_vars_with_pattern = null
	// 		}
	// 		document "test-doc" {
	// 			content text {
	// 				value = "{{.env.OTHER_VAR}}\n{{.env.FABRIC_VAR}}\n{{.env.FABRIC_TEST_VAR}}"
	// 			}
	// 		}
	// 		`,
	// 	},
	// 	"test-doc",
	// 	[]string{
	// 		"<no value>\n<no value>\n<no value>",
	// 	},
	// 	diagtest.Asserts{},
	// )
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
			data json "name" {
				path = "testdata/a.json"
			}
			document "test-doc" {
				data ref {
					base = data.json.name
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
			data json "name" {
				path = "testdata/a.json"
			}
			document "test-doc" {
				data ref {
					base = data.json.name
				}
				data ref {
					base = data.json.name
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
				data json "name" {
					path = "testdata/a.json"
				}
				content text {
					value = "From data block: {{.data.json.name.property_for}}"
				}
			}
			`,
		},
		"test-doc",
		[]string{"From data block: a.json"},
		diagtest.Asserts{},
	)
	renderTest(
		t, "Data with jq interaction",
		[]string{
			`
			document "test" {
				vars {
					items = ["a", "b", "c"]
					x = 1
					y = 2
				}
				content text {
					local_var = query_jq(".vars.items | length")
					value = "There are {{ .vars.local }} items"
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
					local_var = query_jq(".document.meta.authors")
					value = <<-EOT
						authors={{ .vars.local | join "," }},
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
						local_var = query_jq(".document.meta.authors[0] + .section.meta.authors[0] + .content.meta.authors[0]")
						value = "author = {{ .vars.local }}"
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
				local_var = query_jq(".section.meta.authors[0] // \"unknown\"")
				value = "author = {{ .vars.local }}"
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
					local_var = query_jq(".document.content.children[0].markdown")
					value = "content[0] = {{ .vars.local }}"
				}
				content text {
					local_var = query_jq(".document.content.children[1].markdown")
					value = "content[1] = {{ .vars.local }}"
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
}
