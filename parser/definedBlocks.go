package parser

import (
	"fmt"
	"maps"
	"slices"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/sanity-io/litter"

	circularRefDetector "github.com/blackstork-io/fabric/pkg/cirularRefDetector"
	"github.com/blackstork-io/fabric/pkg/diagnostics"
)

type DefinedBlocks struct {
	Config    map[Key]*Config
	Documents map[string]*DocumentOrSection
	Sections  map[string]*DocumentOrSection
	Plugins   map[Key]*Plugin
}

func traversalFromExpr(expr hcl.Expression) (path []string, diags diagnostics.Diag) {
	// ignore diags, just checking if the val is null
	val, _ := expr.Value(nil)
	if val.IsNull() {
		// empty ref
		return
	}
	traversal, diag := hcl.AbsTraversalForExpr(expr)
	if diags.ExtendHcl(diag) {
		return
	}
	path = make([]string, len(traversal))
	for i, trav := range traversal {
		switch traverser := trav.(type) {
		case hcl.TraverseRoot:
			path[i] = traverser.Name
		case hcl.TraverseAttr:
			path[i] = traverser.Name
		default:
			diags.Append(&hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Invalid path",
				Detail:   "The path in the attribute can not contain this operation",
				Subject:  traverser.SourceRange().Ptr(),
			})
		}
	}
	if diag.HasErrors() {
		path = nil
	}
	return
}

func (db *DefinedBlocks) Traverse(expr hcl.Expression) (res any, diags diagnostics.Diag) {
	var found bool
	path, diags := traversalFromExpr(expr)
	if diags.HasErrors() {
		return
	}
	if len(path) == 0 {
		diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Invalid path",
			Detail:   "The path is empty",
			Subject:  expr.Range().Ptr(),
		})
		return
	}
	switch path[0] {
	case BlockKindConfig:
		if len(path) != 4 {
			diags.Append(&hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Invalid path",
				Detail:   "The config path should have format config.<plugin_kind>.<plugin_name>.<config_name>",
				Subject:  expr.Range().Ptr(),
			})
			return
		}
		res, found = db.Config[Key{
			PluginKind: path[1],
			PluginName: path[2],
			BlockName:  path[3],
		}]
	case BlockKindContent, BlockKindData:
		if len(path) != 3 {
			diags.Append(&hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Invalid path",
				Detail: fmt.Sprintf(
					"The %s path should have format %s.<plugin_kind>.<plugin_name>",
					path[0], path[0],
				),
				Subject: expr.Range().Ptr(),
			})
			return
		}
		res, found = db.Plugins[Key{
			PluginKind: path[0],
			PluginName: path[1],
			BlockName:  path[2],
		}]
	default:
		diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Invalid path",
			Detail:   fmt.Sprintf("Unknown path root '%s'", path[0]),
			Subject:  expr.Range().Ptr(),
		})
		return
	}
	if !found {
		diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Invalid path",
			Detail:   "Referenced item not found",
			Subject:  expr.Range().Ptr(),
		})
	}
	return
}

func (db *DefinedBlocks) GetPlugin(expr hcl.Expression) (plugin *Plugin, diags diagnostics.Diag) {
	res, diags := db.Traverse(expr)
	if diags.HasErrors() {
		return
	}
	plugin, ok := res.(*Plugin)
	if !ok {
		diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Invalid path",
			Detail:   "This path is not referencing a plugin block",
			Subject:  expr.Range().Ptr(),
		})
	}
	return
}

func (db *DefinedBlocks) GetConfig(expr hcl.Expression) (cfg *Config, diags diagnostics.Diag) {
	res, diags := db.Traverse(expr)
	if diags.HasErrors() {
		return
	}
	cfg, ok := res.(*Config)
	if !ok {
		diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Invalid path",
			Detail:   "This path is not referencing a config block",
			Subject:  expr.Range().Ptr(),
		})
	}
	return
}

func ToHclsyntaxBody(body hcl.Body) *hclsyntax.Body {
	hclsyntaxBody, ok := body.(*hclsyntax.Body)
	if !ok {
		// Should never happen: hcl.Body for hcl documents is always *hclsyntax.Body
		panic("hcl.Body to *hclsyntax.Body failed")
	}
	return hclsyntaxBody
}

func (db *DefinedBlocks) EvaluatePlugin(plugin *Plugin) (res PluginEvaluation, diags diagnostics.Diag) {
	if circularRefDetector.Check(plugin) {
		// This produces a bit of an incorrect error and shouldn't trigger in normal operation
		// but I re-check for the circular refs here out of abundance of caution
		diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Circular reference detected",
			Detail:   "Looped back to this block through reference chain:",
			Subject:  plugin.DefRange().Ptr(),
			Extra:    circularRefDetector.ExtraMarker,
		})
		return
	}
	plugin.once.Do(func() {
		res, diags = db.evaluatePlugin(plugin)
		if diags.HasErrors() {
			return
		}
		var ok bool
		plugin.invoke, ok = res.invocation.(*blockInvocation)
		if !ok {
			// Should never happen
			litter.Dump("Should never happen happened", plugin, res, diags)
			litter.Dump("plugin", plugin.Kind(), plugin.PluginName(), plugin.BlockName())

			diags.Append(&hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Plugin evaluation failed",
				Detail:   "Incorrect invocation type (please contact fabric developers about this error)",
				Subject:  plugin.block.Range().Ptr(),
			})
			return
		}
		plugin.config = res.config
		// May have changed after the ref traversal
		plugin.pluginName = res.PluginName
		plugin.isValid = true
	})
	if !plugin.isValid {
		if diags == nil {
			diags.Append(diagnostics.RepeatedError)
		}
		return
	}
	res = PluginEvaluation{
		PluginName: plugin.PluginName(),
		BlockName:  plugin.BlockName(),
		config:     plugin.config,
		invocation: plugin.invoke,
	}
	return
}

func (db *DefinedBlocks) evaluatePlugin(plugin *Plugin) (res PluginEvaluation, diags diagnostics.Diag) {
	var diag hcl.Diagnostics
	res.PluginName = plugin.PluginName()
	res.BlockName = plugin.BlockName()

	// Parsing body

	var attrs []hcl.AttributeSchema
	if plugin.IsRef() {
		attrs = make([]hcl.AttributeSchema, 1, 2)
		attrs[0] = hcl.AttributeSchema{Name: "base", Required: true}
	} else {
		attrs = make([]hcl.AttributeSchema, 0, 1)
	}
	attrs = append(attrs, hcl.AttributeSchema{Name: BlockKindConfig, Required: false})

	content, restHcl, diag := plugin.block.Body.PartialContent(&hcl.BodySchema{
		Attributes: attrs,
		Blocks:     []hcl.BlockHeaderSchema{{Type: BlockKindConfig, LabelNames: nil}},
	})
	if diags.ExtendHcl(diag) {
		return
	}

	rest := ToHclsyntaxBody(restHcl)
	invocation := &blockInvocation{
		Body: &hclsyntax.Body{
			Attributes: maps.Clone(rest.Attributes),
			Blocks:     slices.Clone(rest.Blocks),
			SrcRange:   rest.SrcRange,
			EndRange:   rest.EndRange,
		},
		defRange: plugin.DefRange(),
	}
	delete(invocation.Body.Attributes, "config")
	delete(invocation.Body.Attributes, "base")

	invocation.Blocks = slices.DeleteFunc(invocation.Blocks, func(b *hclsyntax.Block) bool {
		return b.Type == BlockKindConfig && len(b.Labels) == 0
	})

	// Parsing config

	configAttr := content.Attributes[BlockKindConfig]
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
			Context:  plugin.block.Range().Ptr(),
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
			Context:  plugin.block.Body.Range().Ptr(),
		})
		diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Both config attribute and block are specified",
			Detail:   "Remove one of them",
			Subject:  &configAttr.Range,
			Context:  plugin.block.Body.Range().Ptr(),
		})
	} else if configAttr != nil {
		cfg, diag := db.GetConfig(configAttr.Expr)
		if diags.Extend(diag) {
			return
		}

		res.config = &ConfigPtr{
			cfg: cfg,
			ptr: configAttr,
		}
	} else if configBlock != nil {
		res.config = &Config{
			Block: configBlock,
		}
	} else {
		// if no config provided, look up of the default config
		// will happen in the plugin call
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
				Context:  plugin.block.Body.Range().Ptr(),
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

		litter.Dump("***", res.PluginName, parent.PluginName)
		res.PluginName = parent.PluginName
		if res.BlockName == "" {
			res.BlockName = parent.BlockName
			// TODO: display warning for data plugins? See issue #25
		}
		if res.config == nil {
			res.config = parent.config
		}
		parentInvocation, ok := parent.invocation.(*blockInvocation)
		if !ok {
			// Shouldn't ever happen, just a safeguard
			diags.Append(&hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Invalid reference",
				Detail:   "Parent is of the wrong type",
				Subject:  &configAttr.Range,
				Context:  plugin.block.Body.Range().Ptr(),
			})
			return
		}
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

	res.invocation = invocation

	return
}

func NewDefinedBlocks() *DefinedBlocks {
	return &DefinedBlocks{
		Config:    map[Key]*Config{},
		Documents: map[string]*DocumentOrSection{},
		Sections:  map[string]*DocumentOrSection{},
		Plugins:   map[Key]*Plugin{},
	}
}

func (db *DefinedBlocks) Merge(other *DefinedBlocks) (diags diagnostics.Diag) {
	for k, v := range other.Config {
		diags.Append(AddIfMissing(db.Config, k, v))
	}
	for k, v := range other.Documents {
		diags.Append(AddIfMissing(db.Documents, k, v))
	}
	for k, v := range other.Sections {
		diags.Append(AddIfMissing(db.Sections, k, v))
	}
	for k, v := range other.Plugins {
		diags.Append(AddIfMissing(db.Plugins, k, v))
	}
	return
}

func AddIfMissing[M ~map[K]V, K comparable, V FabricBlock](m M, key K, newBlock V) *hcl.Diagnostic {
	if origBlock, found := m[key]; found {
		kind := origBlock.GetHCLBlock().Type
		origDefRange := origBlock.GetHCLBlock().DefRange
		return &hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  fmt.Sprintf("Duplicate '%s' declaration", kind),
			Detail:   fmt.Sprintf("'%s' with the same name originally defined at %s:%d", kind, origDefRange.Filename, origDefRange.Start.Line),
			Subject:  newBlock.GetHCLBlock().DefRange.Ptr(),
		}
	}
	m[key] = newBlock
	return nil
}
