package builtin

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zclconf/go-cty/cty"

	"github.com/blackstork-io/fabric/internal/testtools"
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
	val := cty.ObjectVal(map[string]cty.Value{
		"start_level": cty.NullVal(cty.Number),
		"end_level":   cty.NullVal(cty.Number),
		"ordered":     cty.NullVal(cty.Bool),
		"scope":       cty.NullVal(cty.String),
	})
	args := testtools.ReencodeCTY(s.T(), s.schema.Args, val, nil)
	ctx := context.Background()
	titleMeta := plugin.MapData{
		"provider": plugin.StringData("title"),
		"plugin":   plugin.StringData("blackstork/builtin"),
	}
	res, diags := s.schema.ContentFunc(ctx, &plugin.ProvideContentParams{
		Args: args,
		DataContext: plugin.MapData{
			"document": plugin.MapData{
				"content": plugin.MapData{
					"type": plugin.StringData("section"),
					"children": plugin.ListData{
						plugin.MapData{
							"type":     plugin.StringData("element"),
							"markdown": plugin.StringData("# Header 1"),
							"meta":     titleMeta,
						},
						plugin.MapData{
							"type":     plugin.StringData("element"),
							"markdown": plugin.StringData("Lorem ipsum dolor sit amet, consectetur adipiscing elit."),
						},
						plugin.MapData{
							"type":     plugin.StringData("element"),
							"markdown": plugin.StringData("## Header 2"),
							"meta":     titleMeta,
						},
						plugin.MapData{
							"type":     plugin.StringData("element"),
							"markdown": plugin.StringData("Vestibulum nec odio."),
						},
						plugin.MapData{
							"type":     plugin.StringData("element"),
							"markdown": plugin.StringData("### Header 3"),
							"meta":     titleMeta,
						},
						plugin.MapData{
							"type":     plugin.StringData("element"),
							"markdown": plugin.StringData("Integer sit amet."),
						},
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
	}, "\n")+"\n", res.Content.Print())
}

func (s *TOCContentTestSuite) TestAdvanced() {
	val := cty.ObjectVal(map[string]cty.Value{
		"start_level": cty.NumberIntVal(1),
		"end_level":   cty.NumberIntVal(2),
		"ordered":     cty.True,
		"scope":       cty.StringVal("document"),
	})
	args := testtools.ReencodeCTY(s.T(), s.schema.Args, val, nil)
	ctx := context.Background()
	titleMeta := plugin.MapData{
		"provider": plugin.StringData("title"),
		"plugin":   plugin.StringData("blackstork/builtin"),
	}
	res, diags := s.schema.ContentFunc(ctx, &plugin.ProvideContentParams{
		Args: args,

		DataContext: plugin.MapData{
			"document": plugin.MapData{
				"content": plugin.MapData{
					"type": plugin.StringData("section"),
					"children": plugin.ListData{
						plugin.MapData{
							"type":     plugin.StringData("element"),
							"markdown": plugin.StringData("# Header 1"),
							"meta":     titleMeta,
						},
						plugin.MapData{
							"type":     plugin.StringData("element"),
							"markdown": plugin.StringData("Lorem ipsum dolor sit amet, consectetur adipiscing elit."),
							"meta":     titleMeta,
						},
						plugin.MapData{
							"type":     plugin.StringData("element"),
							"markdown": plugin.StringData("## Header 2"),
							"meta":     titleMeta,
						},
						plugin.MapData{
							"type":     plugin.StringData("element"),
							"markdown": plugin.StringData("Vestibulum nec odio."),
						},
						plugin.MapData{
							"type":     plugin.StringData("element"),
							"markdown": plugin.StringData("### Header 3"),
							"meta":     titleMeta,
						},
						plugin.MapData{
							"type":     plugin.StringData("element"),
							"markdown": plugin.StringData("Integer sit amet."),
						},
						plugin.MapData{
							"type":     plugin.StringData("element"),
							"markdown": plugin.StringData("## Header 4"),
							"meta":     titleMeta,
						},
						plugin.MapData{
							"type":     plugin.StringData("element"),
							"markdown": plugin.StringData("Vestibulum nec odio."),
						},
						plugin.MapData{
							"type":     plugin.StringData("element"),
							"markdown": plugin.StringData("## Header 5"),
							"meta":     titleMeta,
						},
						plugin.MapData{
							"type":     plugin.StringData("element"),
							"markdown": plugin.StringData("Vestibulum nec odio."),
						},
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
	}, "\n")+"\n", res.Content.Print())
}
