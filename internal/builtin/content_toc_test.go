package builtin

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zclconf/go-cty/cty"

	"github.com/blackstork-io/fabric/plugin"
)

type TOCContentTestSuite struct {
	suite.Suite
	schema *plugin.ContentProvider
}

func TestTOCContentTestSuite(t *testing.T) {
	suite.Run(t, new(TOCContentTestSuite))
}

func (s *TOCContentTestSuite) SetupSuite() {
	s.schema = makeTOCContentProvider()
}

func (s *TOCContentTestSuite) TestSchema() {
	s.Require().NotNil(s.schema)
	s.Nil(s.schema.Config)
	s.NotNil(s.schema.Args)
	s.NotNil(s.schema.ContentFunc)
}

func (s *TOCContentTestSuite) TestSimple() {
	schema := makeTOCContentProvider()
	ctx := context.Background()
	res, diags := schema.ContentFunc(ctx, &plugin.ProvideContentParams{
		Args: cty.ObjectVal(map[string]cty.Value{
			"start_level": cty.NullVal(cty.Number),
			"end_level":   cty.NullVal(cty.Number),
			"ordered":     cty.NullVal(cty.Bool),
			"scope":       cty.NullVal(cty.String),
		}),
		DataContext: plugin.MapData{
			"document": plugin.MapData{
				"content": plugin.ListData{
					plugin.MapData{
						"markdown": plugin.StringData("# Header 1"),
					},
					plugin.MapData{
						"markdown": plugin.StringData("Lorem ipsum dolor sit amet, consectetur adipiscing elit."),
					},
					plugin.MapData{
						"markdown": plugin.StringData("## Header 2"),
					},
					plugin.MapData{
						"markdown": plugin.StringData("Vestibulum nec odio."),
					},
					plugin.MapData{
						"markdown": plugin.StringData("### Header 3"),
					},
					plugin.MapData{
						"markdown": plugin.StringData("Integer sit amet."),
					},
				},
			},
		},
	})
	s.Len(diags, 0, "no errors")
	s.Equal(strings.Join([]string{
		"- [Header 1](#header-1)",
		"   - [Header 2](#header-2)",
		"      - [Header 3](#header-3)",
	}, "\n")+"\n", res.Markdown)
}

func (s *TOCContentTestSuite) TestAdvanced() {
	schema := makeTOCContentProvider()
	ctx := context.Background()
	res, diags := schema.ContentFunc(ctx, &plugin.ProvideContentParams{
		Args: cty.ObjectVal(map[string]cty.Value{
			"start_level": cty.NumberIntVal(2),
			"end_level":   cty.NumberIntVal(3),
			"ordered":     cty.True,
			"scope":       cty.StringVal("document"),
		}),
		DataContext: plugin.MapData{
			"document": plugin.MapData{
				"content": plugin.ListData{
					plugin.MapData{
						"markdown": plugin.StringData("# Header 1"),
					},
					plugin.MapData{
						"markdown": plugin.StringData("Lorem ipsum dolor sit amet, consectetur adipiscing elit."),
					},
					plugin.MapData{
						"markdown": plugin.StringData("## Header 2"),
					},
					plugin.MapData{
						"markdown": plugin.StringData("Vestibulum nec odio."),
					},
					plugin.MapData{
						"markdown": plugin.StringData("### Header 3"),
					},
					plugin.MapData{
						"markdown": plugin.StringData("Integer sit amet."),
					},
					plugin.MapData{
						"markdown": plugin.StringData("## Header 4"),
					},
					plugin.MapData{
						"markdown": plugin.StringData("Vestibulum nec odio."),
					},
					plugin.MapData{
						"markdown": plugin.StringData("## Header 5"),
					},
					plugin.MapData{
						"markdown": plugin.StringData("Vestibulum nec odio."),
					},
				},
			},
		},
	})
	s.Len(diags, 0, "no errors")
	s.Equal(strings.Join([]string{
		"1. [Header 2](#header-2)",
		"   1. [Header 3](#header-3)",
		"2. [Header 4](#header-4)",
		"3. [Header 5](#header-5)",
	}, "\n")+"\n", res.Markdown)
}
