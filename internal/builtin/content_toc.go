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
	minTOCLevel          = 1
	maxTOCLevel          = 6
	defaultTOCStartLevel = 1
	defaultTOCEndLevel   = 3
	defaultTOCOrdered    = false
	defaultTOCScope      = "document"
)

var availableTOCScopes = []string{"document", "section"}

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
		scope = cty.StringVal(defaultTOCScope)
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

func genTOC(ctx context.Context, params *plugin.ProvideContentParams) (*plugin.Content, hcl.Diagnostics) {
	args, err := parseTOCArgs(params.Args)
	if err != nil {
		return nil, hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "Failed to parse arguments",
			Detail:   err.Error(),
		}}
	}
	var scopedCtx plugin.Data
	if args.scope == "section" {
		section, ok := params.DataContext["section"]
		if !ok {
			return nil, hcl.Diagnostics{{
				Severity: hcl.DiagError,
				Summary:  "No section context",
				Detail:   "No section context found",
			}}
		}
		scopedCtx = section
	} else {
		doc, ok := params.DataContext["document"]
		if !ok {
			return nil, hcl.Diagnostics{{
				Severity: hcl.DiagError,
				Summary:  "No document context",
				Detail:   "No document context found",
			}}
		}
		scopedCtx = doc
	}
	content, ok := scopedCtx.(plugin.MapData)["content"]
	if !ok {
		return nil, hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "No content context",
			Detail:   "No content context found",
		}}
	}
	titles, err := parseContentTitles(content, args.startLevel, args.endLevel)
	if err != nil {
		return nil, hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "Failed to parse content titles",
			Detail:   err.Error(),
		}}
	}

	return &plugin.Content{
		Markdown: titles.render(0, args.ordered),
		Location: &plugin.Location{
			Index:  0,
			Effect: plugin.LocationEffectBefore,
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

func parseContentTitles(data plugin.Data, startLvl, endLvl int) (tocNodeList, error) {
	list, ok := data.(plugin.ListData)
	if !ok {
		return nil, fmt.Errorf("expected a list of content titles")
	}
	var result tocNodeList
	for _, item := range list {
		elem, ok := item.(plugin.MapData)
		if !ok {
			return nil, fmt.Errorf("expected a string")
		}
		content := plugin.ParseContentData(elem)
		line := strings.TrimSpace(content.Markdown)
		if strings.HasPrefix(line, "#") {
			level := strings.Count(line, "#")
			if level < startLvl || level > endLvl {
				continue
			}
			title := strings.TrimSpace(line[level:])
			result = result.add(tocNode{level: level, title: title})
		}
	}

	return result, nil
}
