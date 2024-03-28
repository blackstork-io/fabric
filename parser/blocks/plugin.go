package blocks

import (
	"fmt"
	"maps"
	"slices"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/hashicorp/hcl/v2/hclsyntax"

	"github.com/blackstork-io/fabric/parser/blocks/internal/tree"
	"github.com/blackstork-io/fabric/parser/definitions"
	"github.com/blackstork-io/fabric/parser/evaluation"
	"github.com/blackstork-io/fabric/pkg/diagnostics"
	"github.com/blackstork-io/fabric/pkg/utils"
	plg "github.com/blackstork-io/fabric/plugin"
	"github.com/blackstork-io/fabric/plugin/runner"
)

type plugin struct {
	Block *hclsyntax.Block

	// Parse fills these
	Meta       *definitions.MetaBlock
	Config     evaluation.Configuration
	ParsedBody hclsyntax.Body
	Base       Plugin

	// Filled by invoke
	Result string
}

func (p *plugin) HCLBlock() *hclsyntax.Block {
	return p.Block
}

func (p *plugin) commonPlugin() *plugin {
	return p
}

func (p *plugin) DefRange() hcl.Range {
	return p.Block.DefRange()
}

func (p *plugin) Kind() string {
	return p.Block.Type
}

// Returns originally defined plugin name (will be `ref` on ref plugins)
func (p *plugin) DefName() string {
	return p.Block.Labels[0]
}

func (p *plugin) Name() string {
	if !p.IsRef() || p.Base == nil {
		return p.DefName()
	} else {
		return p.Base.Name()
	}
}

// Whether or not the original block is a ref.
func (p *plugin) IsRef() bool {
	return p.Block.Labels[0] == PluginTypeRef
}

func (p *plugin) BlockName() string {
	if len(p.Block.Labels) < 2 {
		return ""
	}
	return p.Block.Labels[1]
}

func definePlugin(block *hclsyntax.Block, atTopLevel bool) (res plugin, diags diagnostics.Diag) {
	nameRequired := atTopLevel

	// diags.Append(validatePluginKind(block, block.Type, block.TypeRange))
	diags.Append(validatePluginName(block, 0))
	diags.Append(validateBlockName(block, 1, nameRequired))
	var usage string
	if nameRequired {
		usage = "plugin_name block_name"
	} else {
		usage = "plugin_name <block_name>"
	}
	diags.Append(validateLabelsLength(block, 2, usage))
	if diags.HasErrors() {
		return
	}
	res = plugin{
		Block: block,
	}
	return
}

func (p *plugin) fillParsedBody(ctx interface {
	Traverser(expr hclsyntax.Expression) (tree.Node, diagnostics.Diag)
},
) (diags diagnostics.Diag) {
	p.ParsedBody = *p.Block.Body
	p.ParsedBody.Attributes = maps.Clone(p.ParsedBody.Attributes)
	p.ParsedBody.Blocks = slices.Clone(p.ParsedBody.Blocks)

	if !p.IsRef() {
		if baseAttr, hasBase := p.ParsedBody.Attributes[AttrBase]; hasBase {
			diags.Append(&hcl.Diagnostic{
				Severity: hcl.DiagWarning,
				Summary:  "Non-ref block contains '" + AttrBase + "' attribute",
				Detail:   "Did you mean to make it a 'ref'?",
				Subject:  baseAttr.Range().Ptr(),
				Context:  &p.ParsedBody.SrcRange,
			})
		}
		return
	}
	baseAttr, _ := utils.Pop(p.ParsedBody.Attributes, AttrBase)
	if baseAttr == nil {
		diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Ref block missing '" + AttrBase + "' attribute",
			Detail:   "Ref blocks must contain the '" + AttrBase + "' attribute",
			Subject:  p.ParsedBody.MissingItemRange().Ptr(),
			Context:  &p.ParsedBody.SrcRange,
		})
		return
	}

	baseNode, diag := ctx.Traverser(baseAttr.Expr)
	if diags.Extend(diag) {
		return
	}
	basePlugin, ok := baseNode.(Plugin)

	// TODO: ensure that the base node is parsed
	if !ok || basePlugin.Kind() != p.Kind() {
		diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Invalid " + AttrBase,
			Detail:   fmt.Sprintf("'%s' should refer to %s block, not %s", AttrBase, p.Kind(), baseNode.FriendlyName()),
			Subject:  &baseAttr.SrcRange,
			Context:  &p.ParsedBody.SrcRange,
		})
		return
	}

	p.Base = basePlugin
	commonPlugin := p.Base.commonPlugin()

	// Extend attributes and blocks from base
	for k, v := range commonPlugin.ParsedBody.Attributes {
		if _, found := p.ParsedBody.Attributes[k]; !found {
			p.ParsedBody.Attributes[k] = v
		}
	}
	definedBlocks := map[string]struct{}{}
	for _, b := range p.ParsedBody.Blocks {
		definedBlocks[b.Type] = struct{}{}
	}
	for _, b := range commonPlugin.ParsedBody.Blocks {
		if _, found := definedBlocks[b.Type]; !found {
			p.ParsedBody.Blocks = append(p.ParsedBody.Blocks, b)
		}
	}

	p.Config = commonPlugin.Config
	p.Meta = commonPlugin.Meta

	return
}

func (p *plugin) parseSpecial(ctx interface {
	GetEvalContext() *hcl.EvalContext
	GetDefinedBlocks() *DefinedBlocks
	Traverser(expr hclsyntax.Expression) (tree.Node, diagnostics.Diag)
},
) (diags diagnostics.Diag) {
	evalCtx := ctx.GetEvalContext()

	// parse special blocks and attributes
	configAttr, _ := utils.Pop(p.ParsedBody.Attributes, AttrConfig)

	var configBlock *Config

	diag := parseBlocksOnce(&p.ParsedBody, map[string]func(*hclsyntax.Block) (consumed bool, diags diagnostics.Diag){
		BlockKindConfig: func(b *hclsyntax.Block) (consumed bool, diags diagnostics.Diag) {
			cb, diags := DefineConfig(b)
			if diags.HasErrors() {
				return
			}
			if diags.Append(cb.ApplicabilityTest(
				p.Kind(), p.Name(), b.DefRange().Ptr(),
			)) {
				return
			}
			configBlock = cb
			consumed = true
			return
		},
		BlockKindMeta: func(b *hclsyntax.Block) (consumed bool, diags diagnostics.Diag) {
			var meta definitions.MetaBlock

			if diags.ExtendHcl(
				gohcl.DecodeBody(b.Body, evalCtx, &meta),
			) {
				return
			}
			p.Meta = &meta
			consumed = true
			return
		},
	})

	diags.Extend(diag)

	if configAttr != nil && configBlock != nil {
		diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Both config attribute and block are specified",
			Detail:   "Remove one of them",
			Subject:  configBlock.Block.DefRange().Ptr(),
			Context:  &p.Block.Body.SrcRange,
		})
		diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Both config attribute and block are specified",
			Detail:   "Remove one of them",
			Subject:  configAttr.Range().Ptr(),
			Context:  &p.Block.Body.SrcRange,
		})
	}

	if configAttr != nil {
		diags.Extend(p.parseConfigAttr(ctx, configAttr))
	} else if configBlock != nil {
		p.Config = configBlock
	} else if p.Config == nil { // no config specified and none inherited
		defCfg := ctx.GetDefinedBlocks().DefaultConfig(p.Kind(), p.Name())
		if defCfg != nil {
			p.Config = defCfg
		} else {
			p.Config = NewEmptyConfig(p.Block.Body.MissingItemRange())
		}
	}

	return
}

func (p *plugin) parseConfigAttr(ctx interface {
	Traverser(expr hclsyntax.Expression) (tree.Node, diagnostics.Diag)
}, configAttr *hclsyntax.Attribute,
) (diags diagnostics.Diag) {
	node, diag := ctx.Traverser(configAttr.Expr)
	if diags.Extend(diag) {
		return
	}
	config, ok := node.(*Config)
	if !ok {
		diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Invalid " + AttrConfig,
			Detail:   fmt.Sprintf("'%s' should refer to %s, not %s", AttrConfig, (*Config)(nil).FriendlyName(), node.FriendlyName()),
			Subject:  configAttr.Expr.Range().Ptr(),
		})
		return
	}

	// TODO: ensure that the config node is parsed
	if diags.Append(
		config.ApplicabilityTest(
			p.Kind(),
			p.BlockName(),
			configAttr.Expr.Range().Ptr(),
		),
	) {
		return
	}

	p.Config = &ConfigPtr{
		Config: config,
		Rng:    configAttr.SrcRange,
	}
	return
}

func parseBlocksOnce(body *hclsyntax.Body, parsers map[string]func(*hclsyntax.Block) (consumed bool, diags diagnostics.Diag)) (diags diagnostics.Diag) {
	parsedBlocks := make(map[string]struct{})
	body.Blocks = slices.DeleteFunc(body.Blocks, func(block *hclsyntax.Block) bool {
		_, found := parsedBlocks[block.Type]
		if found {
			diags.Append(&hcl.Diagnostic{
				Severity: hcl.DiagWarning,
				Summary:  fmt.Sprintf("Repeated '%s' block", block.Type),
				Detail:   fmt.Sprintf("No more than one '%s' block is allowed. Only the first one will be used.", block.Type),
				Subject:  block.DefRange().Ptr(),
				Context:  &body.SrcRange,
			})
			return true
		}
		parser, found := parsers[block.Type]
		if !found {
			return false
		}
		consumed, diag := parser(block)
		diags.Extend(diag)
		if !consumed {
			return false
		}
		parsedBlocks[block.Type] = struct{}{}
		return true
	})
	return
}

func (p *plugin) Invoke(ctx interface {
	evaluation.PluginCaller
}, r *runner.Runner,
) (diags diagnostics.Diag) {
	switch p.Kind() {
	case BlockKindContent:
		d, _ := r.ContentProvider(p.Name())
		cfg, _ := hcldec.Decode(p.Config.GetBody(), d.Config, nil)
		arg, _ := hcldec.Decode(&p.ParsedBody, d.Args, nil)
		content, _ := d.ContentFunc(nil, &plg.ProvideContentParams{
			Config: cfg,
			Args:   arg,
		})
		p.Result = content.Markdown
	}
	return
}
