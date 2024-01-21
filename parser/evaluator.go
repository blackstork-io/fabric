package parser

import (
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/sanity-io/litter"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/convert"

	"github.com/blackstork-io/fabric/pkg/diagnostics"
)

// Evaluates a chosen document

type MetaBlock struct {
	// XXX: is empty sting enougth or use a proper ptr-nil-if-missing?
	Author string   `hcl:"author,optional"`
	Tags   []string `hcl:"tags,optional"`
}

var _ Configuration = (*Config)(nil)

type ConfigPtr struct {
	cfg *Config
	ptr *hcl.Attribute
}

// Parse implements ConfigurationObject.
func (c *ConfigPtr) Parse(spec hcldec.Spec) (val cty.Value, diags diagnostics.Diag) {
	return c.cfg.Parse(spec)
}

// Range implements ConfigurationObject.
func (c *ConfigPtr) Range() hcl.Range {
	// Use the location of "config = *traversal*" for error reporting, not original config's Range
	return c.ptr.Range
}

var _ Configuration = (*ConfigPtr)(nil)

type InlineConfigBlock struct {
	*Config
}

var _ Configuration = (*InlineConfigBlock)(nil)

type titleInvocation hclsyntax.Attribute

func newTitleInvocation(title *hclsyntax.Attribute) *titleInvocation {
	return (*titleInvocation)(title)
}

var _ Invocation = (*titleInvocation)(nil)

func (t *titleInvocation) DefRange() hcl.Range {
	return t.SrcRange
}

func (t *titleInvocation) MissingItemRange() hcl.Range {
	return t.SrcRange
}

// Range implements InvocationObject.
func (t *titleInvocation) Range() hcl.Range {
	return t.SrcRange
}

func (t *titleInvocation) Parse(spec hcldec.Spec) (val cty.Value, diags diagnostics.Diag) {
	titleVal, diag := t.Expr.Value(nil)
	if diags.ExtendHcl(diag) {
		return
	}

	titleStrVal, err := convert.Convert(titleVal, cty.String)
	if err != nil {
		diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Failed to turn title into a string",
			Detail:   err.Error(),
			Subject:  t.Expr.Range().Ptr(),
		})
		return
	}
	// cty.MapVal()?
	val = cty.ObjectVal(map[string]cty.Value{
		"text":      titleStrVal,
		"format_as": cty.StringVal("title"),
	})
	return
}

type Evaluator struct {
	caller         PluginCaller
	contentCalls   []PluginEvaluation
	topLevelBlocks *DefinedBlocks
	context        map[string]any
}

func NewEvaluator(caller PluginCaller, blocks *DefinedBlocks) *Evaluator {
	return &Evaluator{
		caller:         caller,
		topLevelBlocks: blocks,
		context:        map[string]any{},
	}
}

type PluginEvaluation struct {
	PluginName string
	BlockName  string
	config     Configuration
	invocation Invocation
}

func (e *Evaluator) EvaluateDocument(d *DocumentOrSection) (output string, diags diagnostics.Diag) {
	// sections are basically documents
	diags = e.EvaluateSection(d)
	if diags.HasErrors() {
		return
	}

	results := make([]string, 0, len(e.contentCalls))
	for _, call := range e.contentCalls {
		result, diag := e.caller.CallContent(call.PluginName, call.config, call.invocation, e.context)
		if diags.Extend(diag) {
			// XXX: What to do if we have errors while executing content blocks?
			// just skipping the value for now...
			continue
		}
		results = append(results, result)
	}
	output = strings.Join(results, "\n")
	return
}

func (e *Evaluator) EvaluateSection(d *DocumentOrSection) (diags diagnostics.Diag) {
	if title := d.block.Body.Attributes["title"]; title != nil {
		e.contentCalls = append(e.contentCalls, PluginEvaluation{
			PluginName: "text",
			config:     nil,
			invocation: newTitleInvocation(title),
		})
	}

	for _, block := range d.block.Body.Blocks {
		switch block.Type {
		case BlockKindContent, BlockKindData:
			plugin, diag := DefinePlugin(block, false)
			if diags.Extend(diag) {
				continue
			}
			call, diag := e.topLevelBlocks.EvaluatePlugin(plugin)
			if diags.Extend(diag) {
				continue
			}
			switch block.Type {
			case BlockKindContent:
				// delaying content calls until all data calls are completed
				e.contentCalls = append(e.contentCalls, call)
			case BlockKindData:
				litter.Dump("call", call)
				res, diag := e.caller.CallData(
					call.PluginName,
					call.config,
					call.invocation,
				)
				if diags.Extend(diag) {
					continue
				}
				// XXX: place the result in the correct path here
				e.context[call.BlockName] = res
			default:
				panic("must be exhaustive")
			}

		case BlockKindMeta:
			diags.ExtendHcl(gohcl.DecodeBody(block.Body, nil, &d.meta))
		case BlockKindSection:
			section, diag := DefineSectionOrDocument(block, false)
			if diags.Extend(diag) {
				continue
			}
			diag = e.EvaluateSection(section)
			if diags.Extend(diag) {
				continue
			}
		default:
			diags.Append(newNestingDiag(
				d.block.Type,
				block,
				d.block.Body,
				[]string{
					BlockKindContent,
					BlockKindData,
					BlockKindMeta,
					BlockKindSection,
				}))
		}
	}

	return
}

type blockInvocation struct {
	*hclsyntax.Body
	defRange hcl.Range
}

// DefRange implements InvocationObject.
func (b *blockInvocation) DefRange() hcl.Range {
	return b.defRange
}

// Parse implements InvocationObject.
func (b *blockInvocation) Parse(spec hcldec.Spec) (cty.Value, diagnostics.Diag) {
	res, diag := hcldec.Decode(b.Body, spec, nil)
	return res, diagnostics.Diag(diag)
}

// Range implements InvocationObject.
func (b *blockInvocation) Range() hcl.Range {
	return b.Body.Range()
}

var _ Invocation = (*blockInvocation)(nil)
