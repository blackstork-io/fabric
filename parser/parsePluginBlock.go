package parser

import (
	"fmt"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclsyntax"

	"github.com/blackstork-io/fabric/parser/definitions"
	"github.com/blackstork-io/fabric/parser/evaluation"
	circularRefDetector "github.com/blackstork-io/fabric/pkg/cirularRefDetector"
	"github.com/blackstork-io/fabric/pkg/diagnostics"
	"github.com/blackstork-io/fabric/pkg/utils"
)

// Evaluates a defined plugin.
func (db *DefinedBlocks) ParsePlugin(plugin *definitions.Plugin) (res *definitions.ParsedPlugin, diags diagnostics.Diag) {
	if circularRefDetector.Check(plugin) {
		// This produces a bit of an incorrect error and shouldn't trigger in normal operation
		// but I re-check for the circular refs here out of abundance of caution:
		// deadlocks or infinite loops may occur, and are hard to debug
		diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Circular reference detected",
			Detail:   "Looped back to this block through reference chain:",
			Subject:  plugin.DefRange().Ptr(),
			Extra:    circularRefDetector.ExtraMarker,
		})
		return
	}
	plugin.Once.Do(func() {
		res, diags = db.parsePlugin(plugin)
		if diags.HasErrors() {
			return
		}
		plugin.ParseResult = res
		plugin.Parsed = true
	})
	if !plugin.Parsed {
		if diags == nil {
			diags.Append(diagnostics.RepeatedError)
		}
		return
	}
	res = plugin.ParseResult
	return
}

func (db *DefinedBlocks) parsePlugin(plugin *definitions.Plugin) (parsed *definitions.ParsedPlugin, diags diagnostics.Diag) {
	var diag hcl.Diagnostics
	// var res evaluation.Plugin
	res := definitions.ParsedPlugin{
		PluginName: plugin.Name(),
		BlockName:  plugin.BlockName(),
		// Config and Invocation are to-be filled
	}

	// Parsing body

	attrs := []hcl.AttributeSchema{
		{Name: definitions.BlockKindConfig, Required: false},
	}
	if plugin.IsRef() {
		attrs = append(attrs, hcl.AttributeSchema{Name: definitions.AttrRefBase, Required: true})
	}

	partialBody, restHcl, diag := plugin.Block.Body.PartialContent(&hcl.BodySchema{
		Attributes: attrs,
		Blocks: []hcl.BlockHeaderSchema{
			{Type: definitions.BlockKindConfig, LabelNames: nil},
			{Type: definitions.BlockKindMeta, LabelNames: nil},
		},
	})
	if diags.ExtendHcl(diag) {
		return
	}

	configAttr := partialBody.Attributes[definitions.BlockKindConfig]
	var configBlock *hcl.Block

	for _, blk := range partialBody.Blocks {
		switch blk.Type {
		case definitions.BlockKindConfig:
			if configBlock != nil {
				diags.Append(&hcl.Diagnostic{
					Severity: hcl.DiagWarning,
					Summary:  "More than one embedded config block",
					Detail:   "No more than one config block is allowed. Only the first one will be evaluated.",
					Subject:  &blk.DefRange,
					Context:  plugin.Block.Range().Ptr(),
				})
				continue
			}
			configBlock = blk
		case definitions.BlockKindMeta:
			if res.Meta != nil {
				diags.Append(&hcl.Diagnostic{
					Severity: hcl.DiagWarning,
					Summary:  "More than one meta block",
					Detail:   "No more than one meta block is allowed. Only the first one will be used.",
					Subject:  &blk.DefRange,
					Context:  plugin.Block.Range().Ptr(),
				})
				continue
			}

			var meta definitions.MetaBlock
			if diags.ExtendHcl(gohcl.DecodeBody(blk.Body, nil, &meta)) {
				continue
			}
			res.Meta = &meta
		}
	}

	invocation := &evaluation.BlockInvocation{
		Body:            utils.ToHclsyntaxBody(restHcl),
		DefinitionRange: plugin.DefRange(),
	}

	// Parsing the ref
	var refBaseConfig evaluation.Configuration
	if plugin.IsRef() {
		baseEval, diag := db.parseRefBase(plugin, partialBody.Attributes[definitions.AttrRefBase].Expr)
		if diags.Extend(diag) {
			return
		}

		// replaces "ref" with actual name
		res.PluginName = baseEval.PluginName
		// inherit config from parent. Can be overridden later
		refBaseConfig = baseEval.Config
		if res.BlockName == "" {
			// TODO: display warning for data plugins? See issue #25
			res.BlockName = baseEval.BlockName
		}

		updateRefBody(invocation.Body, baseEval.GetBlockInvocation().Body)
	}
	var dgs diagnostics.Diag
	res.Config, dgs = db.parsePluginConfig(plugin, configAttr, configBlock, refBaseConfig)
	if diags.Extend(dgs) {
		return
	}

	res.Invocation = invocation

	// Future-proofing: be careful when refactoring, the rest of the program
	// (specifically the ref handeling) relies on res.invocation being *evaluation.BlockInvocation
	_, ok := res.Invocation.(*evaluation.BlockInvocation)
	if !ok {
		panic("Plugin invocation must be block invocation")
	}

	parsed = &res
	return
}

func (db *DefinedBlocks) parsePluginConfig(plugin *definitions.Plugin, configAttr *hcl.Attribute, configBlock *hcl.Block, refBaseConfig evaluation.Configuration) (config evaluation.Configuration, diags diagnostics.Diag) {
	switch {
	case configAttr != nil && configBlock != nil:
		diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Both config attribute and block are specified",
			Detail:   "Remove one of them",
			Subject:  configBlock.DefRange.Ptr(),
			Context:  plugin.Block.Body.Range().Ptr(),
		})
		diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Both config attribute and block are specified",
			Detail:   "Remove one of them",
			Subject:  &configAttr.Range,
			Context:  plugin.Block.Body.Range().Ptr(),
		})
		return
	case configAttr != nil:
		// config attr referensing top-level config block
		cfg, diag := Resolve[*definitions.Config](db, configAttr.Expr)
		if diags.Extend(diag) {
			return
		}
		if !cfg.ApplicableTo(plugin) {
			diags.Append(&hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Inapplicable configuration",
				Detail:   "This configuration is for another plugin",
				Subject:  &configAttr.Range,
				Context:  plugin.Block.Body.Range().Ptr(),
			})
			return
		}

		config = &definitions.ConfigPtr{
			Cfg: cfg,
			Ptr: configAttr,
		}
	case configBlock != nil:
		// anonymous config block
		config = &definitions.Config{
			Block: configBlock,
		}
	case plugin.IsRef():
		// Config wasn't provided: inherit config from the base block
		config = refBaseConfig
	default:
		if defaultCfg := db.DefaultConfigFor(plugin); defaultCfg != nil {
			// Apply default configs to non-refs only
			config = defaultCfg
		} else {
			config = &definitions.ConfigEmpty{
				MissingItemRange: plugin.Block.Body.MissingItemRange(),
			}
		}
	}
	return
}

func (db *DefinedBlocks) parseRefBase(plugin *definitions.Plugin, base hcl.Expression) (baseEval *definitions.ParsedPlugin, diags diagnostics.Diag) {
	basePlugin, diags := Resolve[*definitions.Plugin](db, base)
	if diags.HasErrors() {
		return
	}

	if plugin.Kind() != basePlugin.Kind() {
		diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Invalid reference",
			Detail:   fmt.Sprintf("'%s ref' block references a different kind of block (%s) in 'base' attribute", plugin.Kind(), basePlugin.Kind()),
			Subject:  base.Range().Ptr(),
			Context:  plugin.Block.Body.Range().Ptr(),
		})
		return
	}

	circularRefDetector.Add(plugin, base.Range().Ptr())
	defer circularRefDetector.Remove(plugin, &diags)
	if circularRefDetector.Check(basePlugin) {
		diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Circular reference detected",
			Detail:   "Looped back to this block through reference chain:",
			Subject:  plugin.DefRange().Ptr(),
			Extra:    circularRefDetector.ExtraMarker,
		})
		return
	}
	baseEval, diag := db.ParsePlugin(basePlugin)
	diags.Extend(diag)
	return
}

// unique key (concatenation of type and labels).
func hclBlockKey(b *hclsyntax.Block) string {
	var sb strings.Builder
	length := len(b.Type) + len(b.Labels)
	for _, l := range b.Labels {
		length += len(l)
	}
	sb.Grow(length)
	sb.WriteString(b.Type)
	for _, l := range b.Labels {
		sb.WriteByte(0)
		sb.WriteString(l)
	}
	return sb.String()
}

func updateRefBody(ref, base *hclsyntax.Body) {
	for k, v := range base.Attributes {
		switch k {
		case definitions.AttrRefBase, definitions.BlockKindConfig:
			continue
		default:
			if _, found := ref.Attributes[k]; found {
				continue
			}
			ref.Attributes[k] = v
		}
	}
	refBlocks := make(map[string]struct{}, len(ref.Blocks)+1)
	refBlocks[definitions.BlockKindConfig] = struct{}{} // to prevent us from copying the anonymous config block
	for _, b := range ref.Blocks {
		refBlocks[hclBlockKey(b)] = struct{}{}
	}
	for _, b := range base.Blocks {
		if _, found := refBlocks[hclBlockKey(b)]; found {
			continue
		}
		ref.Blocks = append(ref.Blocks, b)
	}
}
