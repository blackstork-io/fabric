package engine

import (
	"testing"
)

func TestDynamicAndMeta(t *testing.T) {
	renderTest(
		t, "dynamic content",
		[]string{`
			document "test-doc" {
				dynamic {
					items = ["a", "b", "c"]
					content text {
						value = "{{.vars.dynamic_item_index}}: {{.vars.dynamic_item}}"
					}
				}
			}
		`},
		[]string{
			"0: a",
			"1: b",
			"2: c",
		},
	)
	renderTest(
		t, "dynamic conditional content",
		[]string{`
			document "test-doc" {
				vars {
				   show1 = true
				   show2 = false
				}
				dynamic {
					condition = query_jq(".vars.show1")
					content text {
						value = "show1 block"
					}
				}
				dynamic {
					condition = query_jq(".vars.show2")
					content text {
						value = "show2 block"
					}
				}
			}
		`},
		[]string{
			"show1 block",
		},
	)

	renderTest(
		t, "dynamic conditional + items",
		[]string{`
			document "test-doc" {
				vars {
				   show1 = true
				   show2 = false
				}
				dynamic {
					condition = query_jq(".vars.show1")
					items = ["dyn_item"]
					content text {
						value = "show1 block {{.vars.dynamic_item}}"
					}
				}
				dynamic {
					condition = query_jq(".vars.show2")
					items = ["dyn_item"]
					content text {
						value = "show2 block {{.vars.dynamic_item}}"
					}
				}
			}
		`},
		[]string{
			"show1 block dyn_item",
		},
	)

	renderTest(
		t, "nested dynamics",
		[]string{`
			document "test-doc" {
				dynamic {
					items = ["abc", "def"]
					content text {
						value = "{{.vars.dynamic_item}} by letters:"
					}
					dynamic {
						items = query_jq(".vars.dynamic_item | split(\"\")")
						content text {
							vars {
								idx_human = query_jq(".vars.dynamic_item_index + 1")
							}
							value = "{{.vars.idx_human}}: {{.vars.dynamic_item}}"
						}
					}
				}
			}
		`},
		[]string{
			"abc by letters:",
			"1: a",
			"2: b",
			"3: c",
			"def by letters:",
			"1: d",
			"2: e",
			"3: f",
		},
	)
	renderTest(
		t, "nested dynamics",
		[]string{`
			document "test-doc" {
				dynamic {
					items = ["abc", "def"]
					content text {
						value = "{{.vars.dynamic_item}} by letters:"
					}
					dynamic {
						items = query_jq(".vars.dynamic_item | split(\"\")")
						content text {
							vars {
								idx_human = query_jq(".vars.dynamic_item_index + 1")
							}
							value = "{{.vars.idx_human}}: {{.vars.dynamic_item}}"
						}
					}
				}
			}
		`},
		[]string{
			"abc by letters:",
			"1: a",
			"2: b",
			"3: c",
			"def by letters:",
			"1: d",
			"2: e",
			"3: f",
		},
	)
	renderTest(
		t, "dynamics sections and nested dynamic",
		[]string{`
			document "test-doc" {
				dynamic {
					items = ["A", "B"]
					section {
						content text {
							value = query_jq("\"Section \" + .vars.dynamic_item")
						}
						dynamic {
							items = ["x", "y"]
							content text {
								value = "Content {{.vars.dynamic_item}}"
							}
						}
					}
				}
			}
		`},
		[]string{
			"Section A",
			"Content x",
			"Content y",
			"Section B",
			"Content x",
			"Content y",
		},
	)

	renderTest(
		t, "dynamics sections with titles",
		[]string{`
			document "test-doc" {
				dynamic {
					items = ["A", "B"]
					section {
						title = query_jq("\"Section \" + .vars.dynamic_item")
						dynamic {
							items = ["x", "y"]
							content text {
								value = "Content {{.vars.dynamic_item}}"
							}
						}
					}
				}
			}
		`},
		[]string{
			"## Section A",
			"Content x",
			"Content y",
			"## Section B",
			"Content x",
			"Content y",
		},
	)
}
