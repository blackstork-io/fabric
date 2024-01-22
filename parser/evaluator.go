package parser

import (
	"fmt"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"

	"github.com/blackstork-io/fabric/parser/definitions"
	"github.com/blackstork-io/fabric/parser/evaluation"
	"github.com/blackstork-io/fabric/pkg/diagnostics"
)

// Evaluates a chosen document

type Evaluator struct {
	caller         PluginCaller
	contentCalls   []*evaluation.Plugin
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

func (e *Evaluator) EvaluateDocument(d *definitions.DocumentOrSection) (output string, diags diagnostics.Diag) {
	// sections are basically documents
	diags = e.EvaluateSection(d)
	if diags.HasErrors() {
		return
	}

	results := make([]string, 0, len(e.contentCalls))
	for _, call := range e.contentCalls {
		result, diag := e.caller.CallContent(call.PluginName, call.Config, call.Invocation, e.context)
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

func (e *Evaluator) EvaluateSection(d *definitions.DocumentOrSection) (diags diagnostics.Diag) {
	if title := d.Block.Body.Attributes["title"]; title != nil {
		e.contentCalls = append(e.contentCalls, &evaluation.Plugin{
			PluginName: "text",
			Config:     nil,
			Invocation: definitions.NewTitle(title),
		})
	}

	var origMeta *hcl.Range

	for _, block := range d.Block.Body.Blocks {
		switch block.Type {
		case definitions.BlockKindContent, definitions.BlockKindData:
			plugin, diag := definitions.DefinePlugin(block, false)
			if diags.Extend(diag) {
				continue
			}
			call, diag := e.topLevelBlocks.EvaluatePlugin(plugin)
			if diags.Extend(diag) {
				continue
			}
			switch block.Type {
			case definitions.BlockKindContent:
				// delaying content calls until all data calls are completed
				e.contentCalls = append(e.contentCalls, call)
			case definitions.BlockKindData:
				res, diag := e.caller.CallData(
					call.PluginName,
					call.Config,
					call.Invocation,
				)
				if diags.Extend(diag) {
					continue
				}
				// XXX: place the result in the correct path here
				e.context[call.BlockName] = res
			default:
				panic("must be exhaustive")
			}

		case definitions.BlockKindMeta:
			if origMeta != nil {
				diags.Append(&hcl.Diagnostic{
					Severity: hcl.DiagError,
					Summary:  "Meta block redefinition",
					Detail: fmt.Sprintf(
						"%s block allows at most one meta block, original meta block was defined at %s:%d",
						d.Block.Type, origMeta.Filename, origMeta.Start.Line,
					),
					Subject: block.DefRange().Ptr(),
					Context: d.Block.Body.Range().Ptr(),
				})
				continue
			}
			diags.ExtendHcl(gohcl.DecodeBody(block.Body, nil, &d.Meta))
			origMeta = block.DefRange().Ptr()
		case definitions.BlockKindSection:
			section, diag := definitions.DefineSectionOrDocument(block, false)
			if diags.Extend(diag) {
				continue
			}
			diag = e.EvaluateSection(section)
			if diags.Extend(diag) {
				continue
			}
		default:
			diags.Append(definitions.NewNestingDiag(
				d.Block.Type,
				block,
				d.Block.Body,
				[]string{
					definitions.BlockKindContent,
					definitions.BlockKindData,
					definitions.BlockKindMeta,
					definitions.BlockKindSection,
				}))
		}
	}

	return
}
