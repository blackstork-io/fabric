package parser

import (
	"fmt"
	"maps"
	"slices"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"

	"github.com/blackstork-io/fabric/parser/definitions"
	"github.com/blackstork-io/fabric/parser/evaluation"
	circularRefDetector "github.com/blackstork-io/fabric/pkg/cirularRefDetector"
	"github.com/blackstork-io/fabric/pkg/diagnostics"
)

// Evaluates a defined plugin

func (db *DefinedBlocks) EvaluatePlugin(plugin *definitions.Plugin) (res *evaluation.Plugin, diags diagnostics.Diag) {
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
		res, diags = db.evaluatePlugin(plugin)
		if diags.HasErrors() {
			return
		}
		plugin.EvalResult = res
		plugin.Evaluated = true
	})
	if !plugin.Evaluated {
		if diags == nil {
			diags.Append(diagnostics.RepeatedError)
		}
		return
	}
	res = plugin.EvalResult
	return
}

func (db *DefinedBlocks) evaluatePlugin(plugin *definitions.Plugin) (eval *evaluation.Plugin, diags diagnostics.Diag) {
	var diag hcl.Diagnostics
	var res evaluation.Plugin
	res.PluginName = plugin.Name()
	res.BlockName = plugin.BlockName()

	// Parsing body

	var attrs []hcl.AttributeSchema
	if plugin.IsRef() {
		attrs = make([]hcl.AttributeSchema, 1, 2)
		attrs[0] = hcl.AttributeSchema{Name: "base", Required: true}
	} else {
		attrs = make([]hcl.AttributeSchema, 0, 1)
	}
	attrs = append(attrs, hcl.AttributeSchema{Name: definitions.BlockKindConfig, Required: false})

	content, restHcl, diag := plugin.Block.Body.PartialContent(&hcl.BodySchema{
		Attributes: attrs,
		Blocks:     []hcl.BlockHeaderSchema{{Type: definitions.BlockKindConfig, LabelNames: nil}},
	})
	if diags.ExtendHcl(diag) {
		return
	}

	rest := ToHclsyntaxBody(restHcl)
	invocation := &evaluation.BlockInvocation{
		Body: &hclsyntax.Body{
			Attributes: maps.Clone(rest.Attributes),
			Blocks:     slices.Clone(rest.Blocks),
			SrcRange:   rest.SrcRange,
			EndRange:   rest.EndRange,
		},
		DefinitionRange: plugin.DefRange(),
	}
	delete(invocation.Body.Attributes, "config")
	delete(invocation.Body.Attributes, "base")

	invocation.Blocks = slices.DeleteFunc(invocation.Blocks, func(b *hclsyntax.Block) bool {
		return b.Type == definitions.BlockKindConfig && len(b.Labels) == 0
	})

	// Parsing config

	configAttr := content.Attributes[definitions.BlockKindConfig]
	var configBlock *hcl.Block
	switch len(content.Blocks) {
	case 0:
		// configBlock is nil already
	default:
		diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagWarning,
			Summary:  "More than one embedded config block",
			Detail:   "No more than one config block is allowed. Only the first one will be evaluated.",
			Subject:  hcl.RangeOver(content.Blocks[1].DefRange, content.Blocks[len(content.Blocks)-1].DefRange).Ptr(),
			Context:  plugin.Block.Range().Ptr(),
		})
		fallthrough
	case 1:
		configBlock = content.Blocks[0]
	}

	if configAttr != nil && configBlock != nil {
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
	} else if configAttr != nil {
		cfg, diag := db.GetConfig(configAttr.Expr)
		if diags.Extend(diag) {
			return
		}

		res.Config = &definitions.ConfigPtr{
			Cfg: cfg,
			Ptr: configAttr,
		}
	} else if configBlock != nil {
		res.Config = &definitions.Config{
			Block: configBlock,
		}
	} else {
		// if no config provided, look up the default config
		res.Config = db.Config[definitions.Key{
			PluginKind: plugin.Kind(),
			PluginName: plugin.Name(),
			BlockName:  "",
		}]
	}

	// Parsing the ref

	if plugin.IsRef() {
		base := content.Attributes["base"].Expr
		parentPlugin, dgs := db.GetPlugin(base)
		if diags.Extend(dgs) {
			return
		}

		if plugin.Kind() != parentPlugin.Kind() {
			diags.Append(&hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Invalid reference",
				Detail:   fmt.Sprintf("'%s ref' block references a different kind of block (%s) in 'base' attribute", plugin.Kind(), parentPlugin.Kind()),
				Subject:  &configAttr.Range,
				Context:  plugin.Block.Body.Range().Ptr(),
			})
			return
		}

		circularRefDetector.Add(plugin, base.Range().Ptr())
		defer circularRefDetector.Remove(plugin, &diags)
		if circularRefDetector.Check(parentPlugin) {
			diags.Append(&hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Circular reference detected",
				Detail:   "Looped back to this block through reference chain:",
				Subject:  plugin.DefRange().Ptr(),
				Extra:    circularRefDetector.ExtraMarker,
			})
			return
		}
		parent, dgs := db.EvaluatePlugin(parentPlugin)
		if diags.Extend(dgs) {
			return
		}

		res.PluginName = parent.PluginName
		if res.BlockName == "" {
			res.BlockName = parent.BlockName
			// TODO: display warning for data plugins? See issue #25
		}
		if res.Config == nil {
			res.Config = parent.Config
		}
		parentInvocation := parent.AsBlockInvocation()

		for k, v := range parentInvocation.Attributes {
			switch k {
			case "base", "config":
				continue
			default:
				if _, found := invocation.Attributes[k]; found {
					continue
				}
				invocation.Attributes[k] = v
			}
		}
		key := func(b *hclsyntax.Block) string {
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

		definedBlocks := make(map[string]struct{}, len(invocation.Blocks)+1)
		definedBlocks["config"] = struct{}{}
		for _, b := range invocation.Blocks {
			definedBlocks[key(b)] = struct{}{}
		}
		for _, b := range parentInvocation.Blocks {
			if _, found := definedBlocks[key(b)]; found {
				continue
			}
			invocation.Blocks = append(invocation.Blocks, b)
		}
	}

	if diags.HasErrors() {
		return
	}

	res.Invocation = invocation

	// Future-proofing: be careful when refactoring, the rest of the program
	// (specifically the ref handeling) relies on res.invocation being *evaluation.BlockInvocation
	_, ok := res.Invocation.(*evaluation.BlockInvocation)
	if !ok {
		panic("Plugin invocation must be block invocation")
	}

	eval = &res
	return
}
