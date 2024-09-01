package builtin

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"

	"github.com/blackstork-io/fabric/pkg/diagnostics"
	"github.com/blackstork-io/fabric/plugin"
	"github.com/blackstork-io/fabric/plugin/dataspec"
	"github.com/blackstork-io/fabric/plugin/dataspec/constraint"
	"github.com/blackstork-io/fabric/plugin/plugindata"
)

func makeTOCContentProvider() *plugin.ContentProvider {
	return &plugin.ContentProvider{
		Args: &dataspec.RootSpec{
			Attrs: []*dataspec.AttrSpec{
				{
					Name:         "start_level",
					Type:         cty.Number,
					DefaultVal:   cty.NumberIntVal(0),
					Doc:          `Largest header size which produces entries in the table of contents`,
					MinInclusive: cty.NumberIntVal(0),
					MaxInclusive: cty.NumberIntVal(5),
					Constraints:  constraint.Integer,
				},
				{
					Name:         "end_level",
					Type:         cty.Number,
					DefaultVal:   cty.NumberIntVal(2),
					Doc:          `Smallest header size which produces entries in the table of contents`,
					MinInclusive: cty.NumberIntVal(0),
					MaxInclusive: cty.NumberIntVal(5),
					Constraints:  constraint.Integer,
				},
				{
					Name:       "ordered",
					Type:       cty.Bool,
					DefaultVal: cty.False,
					Doc:        `Whether to use ordered list for the contents`,
				},
				{
					Name: "scope",
					Type: cty.String,
					Doc: `
				Scope of the headers to evaluate.
				  "document" – look for headers in the whole document
				  "section" – look for headers only in the current section
				  "auto" – behaves as "section" if the "toc" block is inside of a section; else – behaves as "document"
				`,
					OneOf: []cty.Value{
						cty.StringVal("document"),
						cty.StringVal("section"),
						cty.StringVal("auto"),
					},
					DefaultVal: cty.StringVal("auto"),
				},
			},
		},
		InvocationOrder: plugin.InvocationOrderEnd,
		ContentFunc:     genTOC,
		Doc: `
			Produces table of contents.

			Inspects the rendered document for headers of a certain size and creates a linked
			table of contents
		`,
	}
}

type tocArgs struct {
	startLevel int
	endLevel   int
	ordered    bool
	scope      string
}

func parseTOCArgs(args *dataspec.Block) (*tocArgs, error) {
	startLevel, _ := args.GetAttrVal("start_level").AsBigFloat().Int64()
	endLevel, _ := args.GetAttrVal("end_level").AsBigFloat().Int64()
	ordered := args.GetAttrVal("ordered").True()
	scope := args.GetAttrVal("scope").AsString()

	return &tocArgs{
		startLevel: int(startLevel),
		endLevel:   int(endLevel),
		ordered:    ordered,
		scope:      scope,
	}, nil
}

func genTOC(ctx context.Context, params *plugin.ProvideContentParams) (*plugin.ContentResult, diagnostics.Diag) {
	args, err := parseTOCArgs(params.Args)
	if err != nil {
		return nil, diagnostics.Diag{{
			Severity: hcl.DiagError,
			Summary:  "Failed to parse arguments",
			Detail:   err.Error(),
		}}
	}
	titles, err := parseContentTitles(params.DataContext, args.startLevel, args.endLevel, args.scope)
	if err != nil {
		return nil, diagnostics.Diag{{
			Severity: hcl.DiagError,
			Summary:  "Failed to parse content titles",
			Detail:   err.Error(),
		}}
	}

	return &plugin.ContentResult{
		Content: &plugin.ContentElement{
			Markdown: titles.render(0, args.ordered),
		},
	}, nil
}

type tocNode struct {
	level    int
	title    string
	children tocNodeList
}

func (n tocNode) render(pos, depth int, ordered bool) string {
	format := "%s- [%s](#%s)\n"
	if ordered {
		format = "%s" + strconv.Itoa(pos+1) + ". [%s](#%s)\n"
	}
	const indentStep = "  ";
	dst := []string{
		fmt.Sprintf(format, strings.Repeat(indentStep, depth), n.title, anchorize(n.title)),
		n.children.render(depth+1, ordered),
	}
	return strings.Join(dst, "")
}

type tocNodeList []tocNode

func (l tocNodeList) render(depth int, ordered bool) string {
	dst := []string{}
	for i, node := range l {
		dst = append(dst, node.render(i, depth, ordered))
	}
	return strings.Join(dst, "")
}

func (l tocNodeList) add(node tocNode) tocNodeList {
	if len(l) == 0 {
		return append(l, node)
	}
	last := l[len(l)-1]
	if last.level < node.level {
		last.children = last.children.add(node)
		l[len(l)-1] = last
	} else {
		l = append(l, node)
	}
	return l
}

func anchorize(s string) string {
	return strings.ToLower(strings.ReplaceAll(s, " ", "-"))
}

func extractTitles(section *plugin.ContentSection) []string {
	var titles []string
	for _, content := range section.Children {
		switch content := content.(type) {
		case *plugin.ContentSection:
			titles = append(titles, extractTitles(content)...)
		case *plugin.ContentElement:
			meta := content.Meta()
			if meta == nil || meta.Plugin != Name || meta.Provider != "title" {
				continue
			}
			titles = append(titles, content.Markdown)
		}
	}
	return titles
}

func parseContentTitles(data plugindata.Map, startLvl, endLvl int, scope string) (tocNodeList, error) {
	document, section := parseScope(data)
	var list []string
	if scope == "auto" {
		if section != nil {
			scope = "section"
		} else {
			scope = "document"
		}
	}
	if scope == "document" {
		list = extractTitles(document)
	} else if scope == "section" && section != nil {
		list = extractTitles(section)
	} else {
		return nil, fmt.Errorf("no content to parse")
	}
	var result tocNodeList
	for _, item := range list {
		line := strings.TrimSpace(item)
		if strings.HasPrefix(line, "#") {
			level := -1
			for i := 0; i < len(line); i++ {
				if line[i] != '#' {
					break
				}
				level++
			}
			if level < startLvl || level > endLvl {
				continue
			}
			title := strings.TrimSpace(line[level+1:])
			result = result.add(tocNode{level: level, title: title})
		}
	}

	return result, nil
}
