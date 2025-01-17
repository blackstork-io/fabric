package plugin

import (
	"context"

	"github.com/hashicorp/hcl/v2"

	"github.com/blackstork-io/fabric/pkg/diagnostics"
	"github.com/blackstork-io/fabric/plugin/ast/nodes"
	"github.com/blackstork-io/fabric/plugin/dataspec"
	"github.com/blackstork-io/fabric/plugin/plugindata"
)

type (
	PublisherInfoFunc func(ctx context.Context, params *PublisherInfoParams) (PublisherInfo, diagnostics.Diag)
	PublishFunc       func(ctx context.Context, params *PublishParams) diagnostics.Diag
)

type PublishParams struct {
	DocumentName string
	Config       *dataspec.Block
	Args         *dataspec.Block
	DataContext  plugindata.Map
	Document     *nodes.Node
}

type PublisherInfoParams struct {
	Config *dataspec.Block
	Args   *dataspec.Block
}

type PublisherInfo struct {
	SupportedCustomNodes []string
	UnsupportedNodes     []string
}

type Publisher struct {
	Doc         string
	Tags        []string
	PublishFunc PublishFunc
	// optional function to provide additional information about the publisher
	PublisherInfoFunc PublisherInfoFunc
	Args              *dataspec.RootSpec
	Config            *dataspec.RootSpec
}

func (pub *Publisher) Validate() diagnostics.Diag {
	var diags diagnostics.Diag
	if pub.PublishFunc == nil {
		diags = append(diags, &hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Incomplete Publisher schema",
			Detail:   "Publisher function not loaded",
		})
	}
	if pub.Args == nil {
		diags = append(diags, &hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Incomplete Publisher schema",
			Detail:   "Missing args schema",
		})
	}
	return diags
}

func (pub *Publisher) Info(ctx context.Context, params *PublisherInfoParams) (info PublisherInfo, diags diagnostics.Diag) {
	if pub == nil {
		diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Missing Publisher schema",
		})
		return
	}
	if pub.PublisherInfoFunc == nil {
		return
	}
	return pub.PublisherInfoFunc(ctx, params)
}

func (pub *Publisher) Execute(ctx context.Context, params *PublishParams) (diags diagnostics.Diag) {
	if pub == nil {
		return diagnostics.Diag{{
			Severity: hcl.DiagError,
			Summary:  "Missing Publisher schema",
		}}
	}
	if pub.PublishFunc == nil {
		return diagnostics.Diag{{
			Severity: hcl.DiagError,
			Summary:  "Incomplete Publisher schema",
			Detail:   "Publish function not loaded",
		}}
	}
	return pub.PublishFunc(ctx, params)
}

type Publishers map[string]*Publisher

func (pubs Publishers) Validate() diagnostics.Diag {
	var diags diagnostics.Diag
	for name, provider := range pubs {
		if provider == nil {
			diags = append(diags, &hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Incomplete Publisher schema",
				Detail:   "Publisher '" + name + "' not loaded",
			})
		} else {
			diags = append(diags, provider.Validate()...)
		}
	}
	return diags
}
