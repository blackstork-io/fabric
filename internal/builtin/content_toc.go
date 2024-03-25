package builtin

import (
	"context"
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/zclconf/go-cty/cty"

	"github.com/blackstork-io/fabric/plugin"
)

const (
	minTOCLevel          = 0
	maxTOCLevel          = 5
	defaultTOCStartLevel = 0
	defaultTOCEndLevel   = 2
	defaultTOCOrdered    = false
)

var availableTOCScopes = []string{"document", "section", "auto"}

func makeTOCContentProvider() *plugin.ContentProvider {
	return &plugin.ContentProvider{
		Args: hcldec.ObjectSpec{
			"start_level": &hcldec.AttrSpec{
				Name:     "start_level",
				Type:     cty.Number,
				Required: false,
			},
			"end_level": &hcldec.AttrSpec{
				Name:     "end_level",
				Type:     cty.Number,
				Required: false,
			},
			"ordered": &hcldec.AttrSpec{
				Name:     "ordered",
				Type:     cty.Bool,
				Required: false,
			},
			"scope": &hcldec.AttrSpec{
				Name:     "scope",
				Type:     cty.String,
				Required: false,
			},
		},
		InvocationOrder: plugin.InvocationOrderEnd,
		ContentFunc:     genTOC,
	}
}

type tocArgs struct {
	startLevel int
	endLevel   int
	ordered    bool
	scope      string
}

func parseTOCArgs(args cty.Value) (*tocArgs, error) {
	if args.IsNull() {
		return nil, fmt.Errorf("arguments are null")
	}
	startLevel := args.GetAttr("start_level")
	if startLevel.IsNull() {
		startLevel = cty.NumberIntVal(defaultTOCStartLevel)
	} else {
		n, _ := startLevel.AsBigFloat().Int64()
		if n < minTOCLevel || n > maxTOCLevel {
			return nil, fmt.Errorf("start_level should be between %d and %d", minTOCLevel, maxTOCLevel)
		}
	}
	endLevel := args.GetAttr("end_level")
	if endLevel.IsNull() {
		endLevel = cty.NumberIntVal(defaultTOCEndLevel)
	} else {
		n, _ := endLevel.AsBigFloat().Int64()
		if n < minTOCLevel || n > maxTOCLevel {
			return nil, fmt.Errorf("end_level should be between %d and %d", minTOCLevel, maxTOCLevel)
		}
	}
	ordered := args.GetAttr("ordered")
	if ordered.IsNull() {
		ordered = cty.BoolVal(defaultTOCOrdered)
	}
	scope := args.GetAttr("scope")
	if scope.IsNull() {
		scope = cty.StringVal("auto")
	} else if !slices.Contains(availableTOCScopes, scope.AsString()) {
		return nil, fmt.Errorf("scope should be one of %s", strings.Join(availableTOCScopes, ", "))
	}
	startLevelI64, _ := startLevel.AsBigFloat().Int64()
	endLevelI64, _ := endLevel.AsBigFloat().Int64()
	return &tocArgs{
		startLevel: int(startLevelI64),
		endLevel:   int(endLevelI64),
		ordered:    ordered.True(),
		scope:      scope.AsString(),
	}, nil
}

func genTOC(ctx context.Context, params *plugin.ProvideContentParams) (*plugin.ContentResult, hcl.Diagnostics) {
	args, err := parseTOCArgs(params.Args)
	if err != nil {
		return nil, hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "Failed to parse arguments",
			Detail:   err.Error(),
		}}
	}
	titles, err := parseContentTitles(params.DataContext, args.startLevel, args.endLevel, args.scope)
	if err != nil {
		return nil, hcl.Diagnostics{{
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
	dst := []string{
		fmt.Sprintf(format, strings.Repeat("   ", depth), n.title, anchorize(n.title)),
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

func parseContentTitles(data plugin.MapData, startLvl, endLvl int, scope string) (tocNodeList, error) {
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
			level := strings.Count(line, "#") - 1
			if level < startLvl || level > endLvl {
				continue
			}
			title := strings.TrimSpace(line[level+1:])
			result = result.add(tocNode{level: level, title: title})
		}
	}

	return result, nil
}
