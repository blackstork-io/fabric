package builtin

import (
	"context"
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/stretchr/testify/suite"
	"github.com/zclconf/go-cty/cty"

	"github.com/blackstork-io/fabric/plugin"
)

type TextGeneratorTestSuite struct {
	suite.Suite
	schema *plugin.ContentProvider
}

func TestTextGeneratorSuite(t *testing.T) {
	suite.Run(t, &TextGeneratorTestSuite{})
}

func (s *TextGeneratorTestSuite) SetupSuite() {
	s.schema = makeTextContentProvider()
}

func (s *TextGeneratorTestSuite) TestSchema() {
	s.Nil(s.schema.Config)
	s.NotNil(s.schema.Args)
	s.NotNil(s.schema.ContentFunc)
}

func (s *TextGeneratorTestSuite) TestMissingText() {
	ctx := context.Background()
	content, diags := s.schema.ContentFunc(ctx, &plugin.ProvideContentParams{
		Args: cty.ObjectVal(map[string]cty.Value{
			"text":                cty.NullVal(cty.String),
			"format_as":           cty.NullVal(cty.String),
			"absolute_title_size": cty.NullVal(cty.Number),
			"code_language":       cty.NullVal(cty.String),
		}),
		DataContext: plugin.MapData{
			"name": plugin.StringData("World"),
		},
	})
	s.Nil(content)
	s.Equal(hcl.Diagnostics{{
		Severity: hcl.DiagError,
		Summary:  "Failed to parse arguments",
		Detail:   "text is required",
	}}, diags)
}

func (s *TextGeneratorTestSuite) TestBasic() {
	ctx := context.Background()
	content, diags := s.schema.ContentFunc(ctx, &plugin.ProvideContentParams{
		Args: cty.ObjectVal(map[string]cty.Value{
			"text":                cty.StringVal("Hello {{.name}}!"),
			"format_as":           cty.NullVal(cty.String),
			"absolute_title_size": cty.NullVal(cty.Number),
			"code_language":       cty.NullVal(cty.String),
		}),
		DataContext: plugin.MapData{
			"name": plugin.StringData("World"),
		},
	})
	s.Empty(diags)
	s.Equal(&plugin.Content{
		Markdown: "Hello World!",
	}, content)
}

func (s *TextGeneratorTestSuite) TestNoTemplate() {
	ctx := context.Background()
	content, diags := s.schema.ContentFunc(ctx, &plugin.ProvideContentParams{
		Args: cty.ObjectVal(map[string]cty.Value{
			"text":                cty.StringVal("Hello World!"),
			"format_as":           cty.NullVal(cty.String),
			"absolute_title_size": cty.NullVal(cty.Number),
			"code_language":       cty.NullVal(cty.String),
		}),
		DataContext: nil,
	})
	s.Empty(diags)
	s.Equal(&plugin.Content{
		Markdown: "Hello World!",
	}, content)
}

func (s *TextGeneratorTestSuite) TestTitleDefault() {
	ctx := context.Background()
	content, diags := s.schema.ContentFunc(ctx, &plugin.ProvideContentParams{
		Args: cty.ObjectVal(map[string]cty.Value{
			"text":                cty.StringVal("Hello {{.name}}!"),
			"format_as":           cty.StringVal("title"),
			"absolute_title_size": cty.NullVal(cty.Number),
			"code_language":       cty.NullVal(cty.String),
		}),
		DataContext: plugin.MapData{
			"name": plugin.StringData("World"),
		},
	})
	s.Empty(diags)
	s.Equal(&plugin.Content{
		Markdown: "# Hello World!",
	}, content)
}

func (s *TextGeneratorTestSuite) TestTitleWithTextMultiline() {
	ctx := context.Background()
	content, diags := s.schema.ContentFunc(ctx, &plugin.ProvideContentParams{
		Args: cty.ObjectVal(map[string]cty.Value{
			"text":                cty.StringVal("Hello\n{{.name}}\nfor you!"),
			"format_as":           cty.StringVal("title"),
			"absolute_title_size": cty.NullVal(cty.Number),
			"code_language":       cty.NullVal(cty.String),
		}),
		DataContext: plugin.MapData{
			"name": plugin.StringData("World"),
		},
	})
	s.Empty(diags)
	s.Equal(&plugin.Content{
		Markdown: "# Hello World for you!",
	}, content)
}

func (s *TextGeneratorTestSuite) TestTitleWithSize() {
	ctx := context.Background()
	content, diags := s.schema.ContentFunc(ctx, &plugin.ProvideContentParams{
		Args: cty.ObjectVal(map[string]cty.Value{
			"text":                cty.StringVal("Hello {{.name}}!"),
			"format_as":           cty.StringVal("title"),
			"absolute_title_size": cty.NumberIntVal(3),
			"code_language":       cty.NullVal(cty.String),
		}),
		DataContext: plugin.MapData{
			"name": plugin.StringData("World"),
		},
	})
	s.Empty(diags)
	s.Equal(&plugin.Content{
		Markdown: "### Hello World!",
	}, content)
}

func (s *TextGeneratorTestSuite) TestTitleWithSizeTooSmall() {
	ctx := context.Background()
	content, diags := s.schema.ContentFunc(ctx, &plugin.ProvideContentParams{
		Args: cty.ObjectVal(map[string]cty.Value{
			"text":                cty.StringVal("Hello {{.name}}!"),
			"format_as":           cty.StringVal("title"),
			"absolute_title_size": cty.NumberIntVal(0),
			"code_language":       cty.NullVal(cty.String),
		}),
		DataContext: plugin.MapData{
			"name": plugin.StringData("World"),
		},
	})
	s.Equal(hcl.Diagnostics{{
		Severity: hcl.DiagError,
		Summary:  "Failed to parse arguments",
		Detail:   "absolute_title_size must be between 1 and 6",
	}}, diags)
	s.Nil(content)
}

func (s *TextGeneratorTestSuite) TestTitleWithSizeTooBig() {
	ctx := context.Background()
	content, diags := s.schema.ContentFunc(ctx, &plugin.ProvideContentParams{
		Args: cty.ObjectVal(map[string]cty.Value{
			"text":                cty.StringVal("Hello {{.name}}!"),
			"format_as":           cty.StringVal("title"),
			"absolute_title_size": cty.NumberIntVal(7),
			"code_language":       cty.NullVal(cty.String),
		}),
		DataContext: plugin.MapData{
			"name": plugin.StringData("World"),
		},
	})
	s.Equal(hcl.Diagnostics{{
		Severity: hcl.DiagError,
		Summary:  "Failed to parse arguments",
		Detail:   "absolute_title_size must be between 1 and 6",
	}}, diags)
	s.Nil(content)
}

func (s *TextGeneratorTestSuite) TestCallInvalidFormat() {
	ctx := context.Background()
	content, diags := s.schema.ContentFunc(ctx, &plugin.ProvideContentParams{
		Args: cty.ObjectVal(map[string]cty.Value{
			"text":                cty.StringVal("Hello World!"),
			"format_as":           cty.StringVal("unknown"),
			"absolute_title_size": cty.NullVal(cty.Number),
			"code_language":       cty.NullVal(cty.String),
		}),
		DataContext: nil,
	})
	s.Nil(content)
	s.Equal(hcl.Diagnostics{{
		Severity: hcl.DiagError,
		Summary:  "Failed to parse arguments",
		Detail:   "format_as must be one of text, title, code, blockquote",
	}}, diags)
}

func (s *TextGeneratorTestSuite) TestCallInvalidTemplate() {
	ctx := context.Background()
	content, diags := s.schema.ContentFunc(ctx, &plugin.ProvideContentParams{
		Args: cty.ObjectVal(map[string]cty.Value{
			"text":                cty.StringVal("Hello {{.name}!"),
			"format_as":           cty.NullVal(cty.String),
			"absolute_title_size": cty.NullVal(cty.Number),
			"code_language":       cty.NullVal(cty.String),
		}),
		DataContext: plugin.MapData{
			"name": plugin.StringData("World"),
		},
	})
	s.Nil(content)
	s.Equal(hcl.Diagnostics{{
		Severity: hcl.DiagError,
		Summary:  "Failed to render text",
		Detail:   "failed to parse text template: template: text:1: bad character U+007D '}'",
	}}, diags)
}

func (s *TextGeneratorTestSuite) TestCallCodeDefault() {
	ctx := context.Background()
	content, diags := s.schema.ContentFunc(ctx, &plugin.ProvideContentParams{
		Args: cty.ObjectVal(map[string]cty.Value{
			"text":                cty.StringVal(`Hello {{.name}}!`),
			"format_as":           cty.StringVal("code"),
			"absolute_title_size": cty.NullVal(cty.Number),
			"code_language":       cty.NullVal(cty.String),
		}),
		DataContext: plugin.MapData{
			"name": plugin.StringData("World"),
		},
	})
	s.Empty(diags)
	s.Equal(&plugin.Content{
		Markdown: "```\nHello World!\n```",
	}, content)
}

func (s *TextGeneratorTestSuite) TestCallCodeWithLanguage() {
	ctx := context.Background()
	content, diags := s.schema.ContentFunc(ctx, &plugin.ProvideContentParams{
		Args: cty.ObjectVal(map[string]cty.Value{
			"text":                cty.StringVal(`{"hello": "{{.name}}"}`),
			"format_as":           cty.StringVal("code"),
			"absolute_title_size": cty.NullVal(cty.Number),
			"code_language":       cty.StringVal("json"),
		}),
		DataContext: plugin.MapData{
			"name": plugin.StringData("world"),
		},
	})
	s.Empty(diags)
	s.Equal(&plugin.Content{
		Markdown: "```json\n{\"hello\": \"world\"}\n```",
	}, content)
}

func (s *TextGeneratorTestSuite) TestCallBlockquote() {
	ctx := context.Background()
	content, diags := s.schema.ContentFunc(ctx, &plugin.ProvideContentParams{
		Args: cty.ObjectVal(map[string]cty.Value{
			"text":                cty.StringVal(`Hello {{.name}}!`),
			"format_as":           cty.StringVal("blockquote"),
			"absolute_title_size": cty.NullVal(cty.Number),
			"code_language":       cty.NullVal(cty.String),
		}),
		DataContext: plugin.MapData{
			"name": plugin.StringData("World"),
		},
	})
	s.Empty(diags)
	s.Equal(&plugin.Content{
		Markdown: "> Hello World!",
	}, content)
}

func (s *TextGeneratorTestSuite) TestCallBlockquoteMultiline() {
	ctx := context.Background()
	content, diags := s.schema.ContentFunc(ctx, &plugin.ProvideContentParams{
		Args: cty.ObjectVal(map[string]cty.Value{
			"text":                cty.StringVal("Hello\n{{.name}}\nfor you!"),
			"format_as":           cty.StringVal("blockquote"),
			"absolute_title_size": cty.NullVal(cty.Number),
			"code_language":       cty.NullVal(cty.String),
		}),
		DataContext: plugin.MapData{
			"name": plugin.StringData("World"),
		},
	})
	s.Empty(diags)
	s.Equal(&plugin.Content{
		Markdown: "> Hello\n> World\n> for you!",
	}, content)
}

func (s *TextGeneratorTestSuite) TestCallBlockquoteMultilineDoubleNewline() {
	ctx := context.Background()
	content, diags := s.schema.ContentFunc(ctx, &plugin.ProvideContentParams{
		Args: cty.ObjectVal(map[string]cty.Value{
			"text":                cty.StringVal("Hello\n{{.name}}\n\nfor you!"),
			"format_as":           cty.StringVal("blockquote"),
			"absolute_title_size": cty.NullVal(cty.Number),
			"code_language":       cty.NullVal(cty.String),
		}),
		DataContext: plugin.MapData{
			"name": plugin.StringData("World"),
		},
	})
	s.Empty(diags)
	s.Equal(&plugin.Content{
		Markdown: "> Hello\n> World\n> \n> for you!",
	}, content)
}
