package stixview

import (
	"bytes"
	"context"
	"crypto/rand"
	_ "embed"
	"encoding/base32"
	"encoding/json"
	"fmt"
	"text/template"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/zclconf/go-cty/cty"

	"github.com/blackstork-io/fabric/plugin"
)

//go:embed stixview.gohtml
var stixViewTmplStr string

var stixViewTmpl *template.Template

func init() {
	stixViewTmpl = template.Must(template.New("stixview").Funcs(template.FuncMap{
		"json": func(v interface{}) string {
			data, err := json.Marshal(v)
			if err != nil {
				return fmt.Sprintf("error: %s", err)
			}
			return string(data)
		},
	}).Parse(stixViewTmplStr))
}

func makeStixViewContentProvider() *plugin.ContentProvider {
	return &plugin.ContentProvider{
		Args: hcldec.ObjectSpec{
			"gist_id": &hcldec.AttrSpec{
				Name:     "gist_id",
				Type:     cty.String,
				Required: false,
			},
			"stix_url": &hcldec.AttrSpec{
				Name:     "stix_url",
				Type:     cty.String,
				Required: false,
			},
			"caption": &hcldec.AttrSpec{
				Name:     "caption",
				Type:     cty.String,
				Required: false,
			},
			"show_footer": &hcldec.AttrSpec{
				Name:     "show_footer",
				Type:     cty.Bool,
				Required: false,
			},
			"show_sidebar": &hcldec.AttrSpec{
				Name:     "show_sidebar",
				Type:     cty.Bool,
				Required: false,
			},
			"show_tlp_as_tags": &hcldec.AttrSpec{
				Name:     "show_tlp_as_tags",
				Type:     cty.Bool,
				Required: false,
			},
			"show_marking_nodes": &hcldec.AttrSpec{
				Name:     "show_marking_nodes",
				Type:     cty.Bool,
				Required: false,
			},
			"show_labels": &hcldec.AttrSpec{
				Name:     "show_labels",
				Type:     cty.Bool,
				Required: false,
			},
			"show_idrefs": &hcldec.AttrSpec{
				Name:     "show_idrefs",
				Type:     cty.Bool,
				Required: false,
			},
			"width": &hcldec.AttrSpec{
				Name:     "width",
				Type:     cty.Number,
				Required: false,
			},
			"height": &hcldec.AttrSpec{
				Name:     "height",
				Type:     cty.Number,
				Required: false,
			},
		},
		ContentFunc: renderStixView,
	}
}

func renderStixView(ctx context.Context, params *plugin.ProvideContentParams) (*plugin.Content, hcl.Diagnostics) {
	args, err := parseStixViewArgs(params.Args)
	if err != nil {
		return nil, hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "Failed to parse arguments",
			Detail:   err.Error(),
		}}
	}
	var uid [16]byte
	_, err = rand.Read(uid[:])
	if err != nil {
		return nil, hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "Failed to generate UID",
			Detail:   err.Error(),
		}}
	}
	rctx := &renderContext{
		Args: args,
		UID:  base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(uid[:]),
	}
	if queryResult, ok := params.DataContext["query_result"]; ok {
		rctx.QueryResult, ok = queryResult.(plugin.ListData)
		if !ok {
			return nil, hcl.Diagnostics{{
				Severity: hcl.DiagError,
				Summary:  "Invalid query result",
				Detail:   "Query result is not a list",
			}}
		}
	} else if args.StixURL == nil && args.GistID == nil {
		return nil, hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "Missing arugments",
			Detail:   "Must provide either stix_url or gist_id or query_result",
		}}
	}
	buf := bytes.NewBufferString("")
	err = stixViewTmpl.Execute(buf, rctx)
	if err != nil {
		return nil, hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "Failed to render template",
			Detail:   err.Error(),
		}}
	}

	return &plugin.Content{
		Markdown: buf.String(),
	}, nil
}

type renderContext struct {
	Args        *stixViewArgs
	UID         string
	QueryResult plugin.ListData
}

type stixViewArgs struct {
	GistID           *string
	StixURL          *string
	Caption          *string
	ShowFooter       *bool
	ShowSidebar      *bool
	ShowTLPAsTags    *bool
	ShowMarkingNodes *bool
	ShowLabels       *bool
	ShowIDRefs       *bool
	Width            *int
	Height           *int
}

func stringPtr(s string) *string {
	return &s
}

func boolPtr(b bool) *bool {
	return &b
}

func intPtr(i int) *int {
	return &i
}

func parseStixViewArgs(args cty.Value) (*stixViewArgs, error) {
	if args.IsNull() {
		return nil, fmt.Errorf("arguments are null")
	}
	var dst stixViewArgs
	gistID := args.GetAttr("gist_id")
	if !gistID.IsNull() && gistID.AsString() != "" {
		dst.GistID = stringPtr(gistID.AsString())
	}
	stixURL := args.GetAttr("stix_url")
	if !stixURL.IsNull() && stixURL.AsString() != "" {
		dst.StixURL = stringPtr(stixURL.AsString())
	}
	caption := args.GetAttr("caption")
	if !caption.IsNull() && caption.AsString() != "" {
		dst.Caption = stringPtr(caption.AsString())
	}
	showFooter := args.GetAttr("show_footer")
	if !showFooter.IsNull() {
		dst.ShowFooter = boolPtr(showFooter.True())
	}
	showSidebar := args.GetAttr("show_sidebar")
	if !showSidebar.IsNull() {
		dst.ShowSidebar = boolPtr(showSidebar.True())
	}
	showTLPAsTags := args.GetAttr("show_tlp_as_tags")
	if !showTLPAsTags.IsNull() {
		dst.ShowTLPAsTags = boolPtr(showTLPAsTags.True())
	}
	showMarkingNodes := args.GetAttr("show_marking_nodes")
	if !showMarkingNodes.IsNull() {
		dst.ShowMarkingNodes = boolPtr(showMarkingNodes.True())
	}
	showLabels := args.GetAttr("show_labels")
	if !showLabels.IsNull() {
		dst.ShowLabels = boolPtr(showLabels.True())
	}
	showIDRefs := args.GetAttr("show_idrefs")
	if !showIDRefs.IsNull() {
		dst.ShowIDRefs = boolPtr(showIDRefs.True())
	}
	width := args.GetAttr("width")
	if !width.IsNull() {
		n, _ := width.AsBigFloat().Int64()
		dst.Width = intPtr(int(n))
	}
	height := args.GetAttr("height")
	if !height.IsNull() {
		n, _ := height.AsBigFloat().Int64()
		dst.Height = intPtr(int(n))
	}
	return &dst, nil
}
