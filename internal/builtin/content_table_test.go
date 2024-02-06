package builtin

import (
	"context"
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/stretchr/testify/suite"
	"github.com/zclconf/go-cty/cty"

	"github.com/blackstork-io/fabric/plugin"
)

type TableGeneratorTestSuite struct {
	suite.Suite
	plugin *plugin.Schema
}

func TestTableGeneratorTestSuite(t *testing.T) {
	suite.Run(t, &TableGeneratorTestSuite{})
}

func (s *TableGeneratorTestSuite) SetupSuite() {
	s.plugin = &plugin.Schema{
		ContentProviders: plugin.ContentProviders{
			"table": makeTableContentProvider(),
		},
	}
}

func (s *TableGeneratorTestSuite) TestSchema() {
	schema := s.plugin.ContentProviders["table"]
	s.NotNil(schema)
	s.Nil(schema.Config)
	s.NotNil(schema.Args)
	s.NotNil(schema.ContentFunc)
}

func (s *TableGeneratorTestSuite) TestNilQueryResult() {
	args := cty.ObjectVal(map[string]cty.Value{
		"columns": cty.ListVal([]cty.Value{
			cty.ObjectVal(map[string]cty.Value{
				"header": cty.StringVal("{{.col_prefix}} Name"),
				"value":  cty.StringVal("{{.name}}"),
			}),
			cty.ObjectVal(map[string]cty.Value{
				"header": cty.StringVal("{{.col_prefix}} Age"),
				"value":  cty.StringVal("{{.age}}"),
			}),
		}),
	})
	ctx := context.Background()
	content, diags := s.plugin.ProvideContent(ctx, "table", &plugin.ProvideContentParams{
		Args: args,
		DataContext: plugin.MapData{
			"col_prefix":   plugin.StringData("User"),
			"query_result": nil,
		},
	})
	s.Equal(&plugin.Content{
		Markdown: "|User Name|User Age|\n|-|-|\n",
	}, content)
	s.Nil(diags)
}

func (s *TableGeneratorTestSuite) TestEmptyQueryResult() {
	args := cty.ObjectVal(map[string]cty.Value{
		"columns": cty.ListVal([]cty.Value{
			cty.ObjectVal(map[string]cty.Value{
				"header": cty.StringVal("{{.col_prefix}} Name"),
				"value":  cty.StringVal("{{.name}}"),
			}),
			cty.ObjectVal(map[string]cty.Value{
				"header": cty.StringVal("{{.col_prefix}} Age"),
				"value":  cty.StringVal("{{.age}}"),
			}),
		}),
	})
	ctx := context.Background()
	content, diags := s.plugin.ProvideContent(ctx, "table", &plugin.ProvideContentParams{
		Args: args,
		DataContext: plugin.MapData{
			"col_prefix":   plugin.StringData("User"),
			"query_result": plugin.ListData{},
		},
	})
	s.Equal(&plugin.Content{
		Markdown: "|User Name|User Age|\n|-|-|\n",
	}, content)
	s.Nil(diags)
}

func (s *TableGeneratorTestSuite) TestBasic() {
	args := cty.ObjectVal(map[string]cty.Value{
		"columns": cty.ListVal([]cty.Value{
			cty.ObjectVal(map[string]cty.Value{
				"header": cty.StringVal("{{.col_prefix}} Name"),
				"value":  cty.StringVal("{{.name}}"),
			}),
			cty.ObjectVal(map[string]cty.Value{
				"header": cty.StringVal("{{.col_prefix}} Age"),
				"value":  cty.StringVal("{{.age}}"),
			}),
		}),
	})
	ctx := context.Background()
	content, diags := s.plugin.ProvideContent(ctx, "table", &plugin.ProvideContentParams{
		Args: args,
		DataContext: plugin.MapData{
			"col_prefix": plugin.StringData("User"),
			"query_result": plugin.ListData{
				plugin.MapData{
					"name": plugin.StringData("John"),
					"age":  plugin.NumberData(42),
				},
				plugin.MapData{
					"name": plugin.StringData("Jane"),
					"age":  plugin.NumberData(43),
				},
			},
		},
	})
	s.Equal(&plugin.Content{
		Markdown: "|User Name|User Age|\n|-|-|\n|John|42|\n|Jane|43|\n",
	}, content)
	s.Nil(diags)
}

func (s *TableGeneratorTestSuite) TestMissingHeader() {
	args := cty.ObjectVal(map[string]cty.Value{
		"columns": cty.ListVal([]cty.Value{
			cty.ObjectVal(map[string]cty.Value{
				"value": cty.StringVal("{{.name}}"),
			}),
			cty.ObjectVal(map[string]cty.Value{
				"value": cty.StringVal("{{.age}}"),
			}),
		}),
	})
	ctx := context.Background()
	content, diags := s.plugin.ProvideContent(ctx, "table", &plugin.ProvideContentParams{
		Args: args,
		DataContext: plugin.MapData{
			"col_prefix": plugin.StringData("User"),
			"query_result": plugin.ListData{
				plugin.MapData{
					"name": plugin.StringData("John"),
					"age":  plugin.NumberData(42),
				},
				plugin.MapData{
					"name": plugin.StringData("Jane"),
					"age":  plugin.NumberData(43),
				},
			},
		},
	})
	s.Nil(content)
	s.Equal(hcl.Diagnostics{{
		Severity: hcl.DiagError,
		Summary:  "Failed to parse arguments",
		Detail:   "missing header in table cell",
	}}, diags)
}

func (s *TableGeneratorTestSuite) TestNilHeader() {
	args := cty.ObjectVal(map[string]cty.Value{
		"columns": cty.ListVal([]cty.Value{
			cty.ObjectVal(map[string]cty.Value{
				"header": cty.NullVal(cty.String),
				"value":  cty.StringVal("{{.name}}"),
			}),
			cty.ObjectVal(map[string]cty.Value{
				"header": cty.NullVal(cty.String),
				"value":  cty.StringVal("{{.age}}"),
			}),
		}),
	})
	ctx := context.Background()
	content, diags := s.plugin.ProvideContent(ctx, "table", &plugin.ProvideContentParams{
		Args: args,
		DataContext: plugin.MapData{
			"col_prefix": plugin.StringData("User"),
			"query_result": plugin.ListData{
				plugin.MapData{
					"name": plugin.StringData("John"),
					"age":  plugin.NumberData(42),
				},
				plugin.MapData{
					"name": plugin.StringData("Jane"),
					"age":  plugin.NumberData(43),
				},
			},
		},
	})
	s.Nil(content)
	s.Equal(hcl.Diagnostics{{
		Severity: hcl.DiagError,
		Summary:  "Failed to parse arguments",
		Detail:   "missing header in table cell",
	}}, diags)
}

func (s *TableGeneratorTestSuite) TestNilValue() {
	args := cty.ObjectVal(map[string]cty.Value{
		"columns": cty.ListVal([]cty.Value{
			cty.ObjectVal(map[string]cty.Value{
				"header": cty.StringVal("{{.col_prefix}} Name"),
				"value":  cty.NullVal(cty.String),
			}),
			cty.ObjectVal(map[string]cty.Value{
				"header": cty.StringVal("{{.col_prefix}} Age"),
				"value":  cty.NullVal(cty.String),
			}),
		}),
	})
	ctx := context.Background()
	content, diags := s.plugin.ProvideContent(ctx, "table", &plugin.ProvideContentParams{
		Args: args,
		DataContext: plugin.MapData{
			"col_prefix": plugin.StringData("User"),
			"query_result": plugin.ListData{
				plugin.MapData{
					"name": plugin.StringData("John"),
					"age":  plugin.NumberData(42),
				},
				plugin.MapData{
					"name": plugin.StringData("Jane"),
					"age":  plugin.NumberData(43),
				},
			},
		},
	})
	s.Nil(content)
	s.Equal(hcl.Diagnostics{{
		Severity: hcl.DiagError,
		Summary:  "Failed to parse arguments",
		Detail:   "missing value in table cell",
	}}, diags)
}

func (s *TableGeneratorTestSuite) TestNilColumns() {
	args := cty.ObjectVal(map[string]cty.Value{
		"columns": cty.NullVal(cty.List(cty.Object(map[string]cty.Type{}))),
	})
	ctx := context.Background()
	content, diags := s.plugin.ProvideContent(ctx, "table", &plugin.ProvideContentParams{
		Args: args,
		DataContext: plugin.MapData{
			"col_prefix": plugin.StringData("User"),
			"query_result": plugin.ListData{
				plugin.MapData{
					"name": plugin.StringData("John"),
					"age":  plugin.NumberData(42),
				},
				plugin.MapData{
					"name": plugin.StringData("Jane"),
					"age":  plugin.NumberData(43),
				},
			},
		},
	})
	s.Nil(content)
	s.Equal(hcl.Diagnostics{{
		Severity: hcl.DiagError,
		Summary:  "Failed to parse arguments",
		Detail:   "columns is required",
	}}, diags)
}

func (s *TableGeneratorTestSuite) TestEmptyColumns() {
	args := cty.ObjectVal(map[string]cty.Value{
		"columns": cty.ListValEmpty(cty.Object(map[string]cty.Type{})),
	})
	ctx := context.Background()
	content, diags := s.plugin.ProvideContent(ctx, "table", &plugin.ProvideContentParams{
		Args: args,
		DataContext: plugin.MapData{
			"col_prefix": plugin.StringData("User"),
			"query_result": plugin.ListData{
				plugin.MapData{
					"name": plugin.StringData("John"),
					"age":  plugin.NumberData(42),
				},
				plugin.MapData{
					"name": plugin.StringData("Jane"),
					"age":  plugin.NumberData(43),
				},
			},
		},
	})
	s.Nil(content)
	s.Equal(hcl.Diagnostics{{
		Severity: hcl.DiagError,
		Summary:  "Failed to parse arguments",
		Detail:   "columns must not be empty",
	}}, diags)
}

func (s *TableGeneratorTestSuite) TestInvalidHeaderTemplate() {
	args := cty.ObjectVal(map[string]cty.Value{
		"columns": cty.ListVal([]cty.Value{
			cty.ObjectVal(map[string]cty.Value{
				"header": cty.StringVal("{{.col_prefix} Name"),
				"value":  cty.StringVal("{{.name}}"),
			}),
			cty.ObjectVal(map[string]cty.Value{
				"header": cty.StringVal("{{.col_prefix}} Age"),
				"value":  cty.StringVal("{{.age}}"),
			}),
		}),
	})
	ctx := context.Background()
	content, diags := s.plugin.ProvideContent(ctx, "table", &plugin.ProvideContentParams{
		Args: args,
		DataContext: plugin.MapData{
			"col_prefix": plugin.StringData("User"),
			"query_result": plugin.ListData{
				plugin.MapData{
					"name": plugin.StringData("John"),
					"age":  plugin.NumberData(42),
				},
				plugin.MapData{
					"name": plugin.StringData("Jane"),
					"age":  plugin.NumberData(43),
				},
			},
		},
	})
	s.Nil(content)
	s.Equal(hcl.Diagnostics{{
		Severity: hcl.DiagError,
		Summary:  "Failed to parse arguments",
		Detail:   "failed to parse header template: template: header:1: bad character U+007D '}'",
	}}, diags)
}

func (s *TableGeneratorTestSuite) TestInvalidValueTemplate() {
	args := cty.ObjectVal(map[string]cty.Value{
		"columns": cty.ListVal([]cty.Value{
			cty.ObjectVal(map[string]cty.Value{
				"header": cty.StringVal("{{.col_prefix}} Name"),
				"value":  cty.StringVal("{{.name}"),
			}),
			cty.ObjectVal(map[string]cty.Value{
				"header": cty.StringVal("{{.col_prefix}} Age"),
				"value":  cty.StringVal("{{.age}}"),
			}),
		}),
	})
	ctx := context.Background()
	content, diags := s.plugin.ProvideContent(ctx, "table", &plugin.ProvideContentParams{
		Args: args,
		DataContext: plugin.MapData{
			"col_prefix": plugin.StringData("User"),
			"query_result": plugin.ListData{
				plugin.MapData{
					"name": plugin.StringData("John"),
					"age":  plugin.NumberData(42),
				},
				plugin.MapData{
					"name": plugin.StringData("Jane"),
					"age":  plugin.NumberData(43),
				},
			},
		},
	})
	s.Nil(content)
	s.Equal(hcl.Diagnostics{{
		Severity: hcl.DiagError,
		Summary:  "Failed to parse arguments",
		Detail:   "failed to parse value template: template: value:1: bad character U+007D '}'",
	}}, diags)
}
