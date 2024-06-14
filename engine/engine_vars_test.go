package engine

import (
	"testing"

	"github.com/blackstork-io/fabric/pkg/diagnostics/diagtest"
)

func TestEngineVarsHandling(t *testing.T) {
	renderTest(
		t, "refs query in override",
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
		t, "refs query in override",
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
		t, "inheritance",
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
			`1: {
  "contentVar": "contentVar",
  "docVar": "docVar",
  "sectVar": "sectVar"
}`,
			`2: {
  "docVar": "docVar",
  "sectVar": "sectVar"
}`,
			`3: {
  "docVar": "docVar"
}`,
		},
		diagtest.Asserts{},
	)
	renderTest(
		t, "combined inheritance and shadowing",
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
		t, "combined inheritance and shadowing section ref",
		[]string{`
			section "sect" {
				vars {
					v3 = query_jq(".vars.v2 + 1")
					v4 = query_jq(".vars.v3 + 1")
					v5 = query_jq(".vars.v4 + 1")
					v6 = "not evaluated"
				}
				content text {
					vars {
						v7 = query_jq(".vars.v6 + 1")
						v8 = query_jq(".vars.v7 + 1")
					}
					value = "{{toPrettyJson .vars}}"
				}
			}

			document "test-doc" {
				vars {
					v1 = 1
					v2 = query_jq(".vars.v1 + 1")
					v4 = "not evaluated"
					v5 = "not evaluated"
				}
				section ref {
					base = section.sect
					vars {
						v6 = query_jq(".vars.v5 + 1")
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
		t, "deep nesting and complex result type",
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
	renderTest(
		t, "local vars",
		[]string{`
			document "test-doc" {
				content text {
					local_var = "local"
					value = "{{toPrettyJson .vars}}"
				}
			}
		`},
		"test-doc",
		[]string{
			`{
  "local": "local"
}`,
		},
		diagtest.Asserts{},
	)
	renderTest(
		t, "local vars with var local redefinition",
		[]string{`
			document "test-doc" {
				content text {
					vars {
						local = 123
					}
					local_var = "local"
					value = "{{toPrettyJson .vars}}"
				}
			}
		`},
		"test-doc",
		nil,
		diagtest.Asserts{
			{
				diagtest.IsError,
				diagtest.SummaryEquals("Local var redefinition"),
			},
		},
	)
	renderTest(
		t, "local vars with var local redefinition",
		[]string{`
			document "test-doc" {
				content text {
					vars {
						other = 123
					}
					local_var = "local"
					value = "{{toPrettyJson .vars}}"
				}
			}
		`},
		"test-doc",
		[]string{
			`{
  "local": "local",
  "other": 123
}`,
		},
		diagtest.Asserts{
			{
				diagtest.IsWarning,
				diagtest.SummaryEquals("Local var specified together with vars block"),
			},
		},
	)
	renderTest(
		t, "local vars document",
		[]string{`
			document "test-doc" {
				local_var = "local"
				content text {
					value = "{{toPrettyJson .vars}}"
				}
			}
		`},
		"test-doc",
		[]string{
			`{
  "local": "local"
}`,
		},
		diagtest.Asserts{},
	)
	renderTest(
		t, "local vars sections",
		[]string{`
			document "test-doc" {
				section {
					local_var = "local"
					content text {
						value = "{{toPrettyJson .vars}}"
					}
				}
			}
		`},
		"test-doc",
		[]string{
			`{
  "local": "local"
}`,
		},
		diagtest.Asserts{},
	)
}
