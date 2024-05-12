package stixview

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zclconf/go-cty/cty"

	"github.com/blackstork-io/fabric/plugin"
	"github.com/blackstork-io/fabric/print/mdprint"
)

type StixViewTestSuite struct {
	suite.Suite
	schema *plugin.ContentProvider
}

func TestStixViewTestSuite(t *testing.T) {
	suite.Run(t, new(StixViewTestSuite))
}

func (s *StixViewTestSuite) SetupTest() {
	s.schema = makeStixViewContentProvider()
}

func (s *StixViewTestSuite) TestSchema() {
	s.Require().NotNil(s.schema)
	s.Nil(s.schema.Config)
	s.NotNil(s.schema.Args)
	s.NotNil(s.schema.ContentFunc)
}

func (s *StixViewTestSuite) TestGistID() {
	res, diags := s.schema.ContentFunc(context.Background(), &plugin.ProvideContentParams{
		Args: cty.ObjectVal(map[string]cty.Value{
			"gist_id":            cty.StringVal("123"),
			"stix_url":           cty.NullVal(cty.String),
			"caption":            cty.NullVal(cty.String),
			"show_footer":        cty.NullVal(cty.Bool),
			"show_sidebar":       cty.NullVal(cty.Bool),
			"show_tlp_as_tags":   cty.NullVal(cty.Bool),
			"show_marking_nodes": cty.NullVal(cty.Bool),
			"show_labels":        cty.NullVal(cty.Bool),
			"show_idrefs":        cty.NullVal(cty.Bool),
			"width":              cty.NullVal(cty.Number),
			"height":             cty.NullVal(cty.Number),
		}),
		DataContext: plugin.MapData{},
	})
	s.Len(diags, 0)
	s.Equal(strings.Join([]string{
		`<script src="https://unpkg.com/stixview/dist/stixview.bundle.js" type="text/javascript"></script>`,
		`<div data-stix-gist-id="123">`,
		`</div>`,
	}, "\n"), mdprint.PrintString(res.Content))
}

func (s *StixViewTestSuite) TestStixURL() {
	res, diags := s.schema.ContentFunc(context.Background(), &plugin.ProvideContentParams{
		Args: cty.ObjectVal(map[string]cty.Value{
			"gist_id":            cty.NullVal(cty.String),
			"stix_url":           cty.StringVal("https://example.com/stix.json"),
			"caption":            cty.NullVal(cty.String),
			"show_footer":        cty.NullVal(cty.Bool),
			"show_sidebar":       cty.NullVal(cty.Bool),
			"show_tlp_as_tags":   cty.NullVal(cty.Bool),
			"show_marking_nodes": cty.NullVal(cty.Bool),
			"show_labels":        cty.NullVal(cty.Bool),
			"show_idrefs":        cty.NullVal(cty.Bool),
			"width":              cty.NullVal(cty.Number),
			"height":             cty.NullVal(cty.Number),
		}),
		DataContext: plugin.MapData{},
	})
	s.Len(diags, 0)
	s.Equal(strings.Join([]string{
		`<script src="https://unpkg.com/stixview/dist/stixview.bundle.js" type="text/javascript"></script>`,
		`<div data-stix-url="https://example.com/stix.json">`,
		`</div>`,
	}, "\n"), mdprint.PrintString(res.Content))
}

func (s *StixViewTestSuite) TestAllArgs() {
	res, diags := s.schema.ContentFunc(context.Background(), &plugin.ProvideContentParams{
		Args: cty.ObjectVal(map[string]cty.Value{
			"gist_id":            cty.StringVal("123"),
			"stix_url":           cty.NullVal(cty.String),
			"caption":            cty.StringVal("test caption"),
			"show_footer":        cty.BoolVal(true),
			"show_sidebar":       cty.BoolVal(true),
			"show_tlp_as_tags":   cty.BoolVal(true),
			"show_marking_nodes": cty.BoolVal(true),
			"show_labels":        cty.BoolVal(true),
			"show_idrefs":        cty.BoolVal(true),
			"width":              cty.NumberIntVal(400),
			"height":             cty.NumberIntVal(300),
		}),
		DataContext: plugin.MapData{},
	})
	s.Len(diags, 0)
	s.Equal(strings.Join([]string{
		`<script src="https://unpkg.com/stixview/dist/stixview.bundle.js" type="text/javascript"></script>`,
		`<div data-stix-gist-id="123" data-show-sidebar=true data-show-footer=true data-show-tlp-as-tags=true data-caption="test caption" data-show-marking-nodes=true data-show-labels=true data-show-idrefs=true data-graph-width=400 data-graph-height=300>`,
		`</div>`,
	}, "\n"), mdprint.PrintString(res.Content))
}

func (s *StixViewTestSuite) TestQueryResult() {
	res, diags := s.schema.ContentFunc(context.Background(), &plugin.ProvideContentParams{
		Args: cty.ObjectVal(map[string]cty.Value{
			"gist_id":            cty.NullVal(cty.String),
			"stix_url":           cty.NullVal(cty.String),
			"caption":            cty.NullVal(cty.String),
			"show_footer":        cty.NullVal(cty.Bool),
			"show_sidebar":       cty.NullVal(cty.Bool),
			"show_tlp_as_tags":   cty.NullVal(cty.Bool),
			"show_marking_nodes": cty.NullVal(cty.Bool),
			"show_labels":        cty.NullVal(cty.Bool),
			"show_idrefs":        cty.NullVal(cty.Bool),
			"width":              cty.NullVal(cty.Number),
			"height":             cty.NullVal(cty.Number),
		}),
		DataContext: plugin.MapData{
			"query_result": plugin.ListData{
				plugin.MapData{
					"key": plugin.StringData("value"),
				},
			},
		},
	})
	s.Len(diags, 0)
	s.Contains(mdprint.PrintString(res.Content), `<script src="https://unpkg.com/stixview/dist/stixview.bundle.js" type="text/javascript"></script>`)
	s.Contains(mdprint.PrintString(res.Content), `<div id="graph-`)
	s.Contains(mdprint.PrintString(res.Content), `window.stixview.init(`)
	s.Contains(mdprint.PrintString(res.Content), `"objects":  [{"key":"value"}]}`)
}
