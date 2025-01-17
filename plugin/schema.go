package plugin

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/hcl/v2"

	"github.com/blackstork-io/fabric/pkg/diagnostics"
	"github.com/blackstork-io/fabric/plugin/ast/nodes"
	"github.com/blackstork-io/fabric/plugin/plugindata"
)

type Schema struct {
	Name             string
	Version          string
	Doc              string
	Tags             []string
	DataSources      DataSources
	ContentProviders ContentProviders
	Publishers       Publishers
	NodeRenderers    NodeRenderers
}

func (p *Schema) Validate() diagnostics.Diag {
	var diags diagnostics.Diag
	if p.Name == "" {
		diags = append(diags, &hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Incomplete PluginSchema",
			Detail:   "Name not defined",
		})
	}
	if p.DataSources != nil {
		diags = append(diags, p.DataSources.Validate()...)
	}
	if p.ContentProviders != nil {
		diags = append(diags, p.ContentProviders.Validate()...)
	}
	if p.Publishers != nil {
		diags = append(diags, p.Publishers.Validate()...)
	}
	if p.DataSources == nil && p.ContentProviders == nil && p.Publishers == nil {
		diags = append(diags, &hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Incomplete PluginSchema",
			Detail:   "No data sources, content providers or publishers defined",
		})
	}
	return diags
}

func (p *Schema) RetrieveData(ctx context.Context, name string, params *RetrieveDataParams) (_ plugindata.Data, diags diagnostics.Diag) {
	if p == nil {
		return nil, diagnostics.Diag{{
			Severity: hcl.DiagError,
			Summary:  "No schema",
			Detail:   "No schema defined",
		}}
	}
	if p.DataSources == nil {
		return nil, diagnostics.Diag{{
			Severity: hcl.DiagError,
			Summary:  "No data sources",
			Detail:   "No data sources defined in schema",
		}}
	}
	source, ok := p.DataSources[name]
	if !ok || source == nil {
		return nil, diagnostics.Diag{{
			Severity: hcl.DiagError,
			Summary:  "Data source not found",
			Detail:   fmt.Sprintf("Data source '%s' not found in schema", name),
		}}
	}
	return source.Execute(ctx, params)
}

func (p *Schema) ProvideContent(ctx context.Context, name string, params *ProvideContentParams) (_ *ContentElement, diags diagnostics.Diag) {
	if p == nil {
		return nil, diagnostics.Diag{{
			Severity: hcl.DiagError,
			Summary:  "No schema",
			Detail:   "No schema defined",
		}}
	}
	if p.ContentProviders == nil {
		return nil, diagnostics.Diag{{
			Severity: hcl.DiagError,
			Summary:  "No content providers",
			Detail:   "No content providers defined in schema",
		}}
	}
	provider, ok := p.ContentProviders[name]
	if !ok || provider == nil {
		return nil, diagnostics.Diag{{
			Severity: hcl.DiagError,
			Summary:  "Content provider not found",
			Detail:   fmt.Sprintf("Content provider '%s' not found in schema", name),
		}}
	}
	result, diags := provider.Execute(ctx, params)
	if diags.HasErrors() {
		return nil, diags
	}
	result.SetPluginMeta(&nodes.FabricContentMetadata{
		Provider: name,
		Plugin:   p.Name,
		Version:  p.Version,
	})
	return result, diags
}

func (p *Schema) getPublisher(name string) (publisher *Publisher, diags diagnostics.Diag) {
	if p == nil {
		diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "No schema",
			Detail:   "No schema defined",
		})
		return
	}
	if p.Publishers == nil {
		diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "No publishers",
			Detail:   "No publishers defined in schema",
		})
		return
	}
	publisher, ok := p.Publishers[name]
	if !ok || publisher == nil {
		diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Publisher not found",
			Detail:   fmt.Sprintf("Publisher '%s' not found in schema", name),
		})
		return
	}
	return
}

func (p *Schema) PublisherInfo(ctx context.Context, name string, params *PublisherInfoParams) (info PublisherInfo, diags diagnostics.Diag) {
	publisher, diag := p.getPublisher(name)
	if diags.Extend(diag) {
		return
	}
	info, diag = publisher.Info(ctx, params)
	diags.Extend(diag)
	return
}

func (p *Schema) Publish(ctx context.Context, name string, params *PublishParams) (diags diagnostics.Diag) {
	publisher, diag := p.getPublisher(name)
	if diags.Extend(diag) {
		return
	}
	diags.Extend(publisher.Execute(ctx, params))
	return
}

func (p *Schema) RenderNode(ctx context.Context, params *RenderNodeParams) (repl *nodes.Node, diags diagnostics.Diag) {
	if p == nil {
		diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "No schema",
			Detail:   "No schema defined",
		})
		return
	}
	if p.NodeRenderers == nil {
		diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "No node renderers",
			Detail:   "No node renderers defined in schema",
		})
		return
	}
	if params == nil {
		diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "No parameters",
			Detail:   "No parameters provided",
		})
		return
	}
	custom_node := params.Subtree.TraversePath(params.NodePath)
	if custom_node == nil {
		diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Node not found",
			Detail:   "Node not found in the subtree",
		})
		return
	}
	custom, ok := custom_node.Content.(*nodes.Custom)
	if !ok {
		diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Invalid node type",
			Detail:   fmt.Sprintf("Expected custom node, got %T", custom_node.Content),
		})
		return
	}
	// types.blackstork.io/fabric/v1/custom_nodes/<plugin_name>/<node_type>
	split := strings.SplitAfterN(custom.Data.GetTypeUrl(), "/", 6)
	nodeType := split[len(split)-1]
	renderer, found := p.NodeRenderers[custom.GetStrippedNodeType()]
	if !found {
		diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Node renderer not found",
			Detail:   fmt.Sprintf("Node renderer for type '%s' (%s) not found in schema", nodeType, custom.Data.GetTypeUrl()),
		})
		return
	}
	return renderer(ctx, params)
}
