package parser

import (
	"context"
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"

	"github.com/blackstork-io/fabric/parser/definitions"
	evaltree "github.com/blackstork-io/fabric/parser/parsetree"
	"github.com/blackstork-io/fabric/parser/plugincaller"
	"github.com/blackstork-io/fabric/pkg/diagnostics"
	"github.com/blackstork-io/fabric/plugin"
)

// Evaluates a chosen document

type Evaluator struct {
	caller         plugincaller.PluginCaller
	topLevelBlocks *DefinedBlocks
	context        plugin.MapData
}

func NewEvaluator(caller plugincaller.PluginCaller, blocks *DefinedBlocks) *Evaluator {
	return &Evaluator{
		caller:         caller,
		topLevelBlocks: blocks,
		context:        plugin.MapData{},
	}
}

func (e *Evaluator) EvaluateDocument(d *definitions.Document) (results []string, diags diagnostics.Diag) {
	node, diags := e.docToEvalTree(d)
	if diags.HasErrors() {
		return
	}
	results, diag := node.EvalContent(context.TODO(), e.caller)
	diags.Extend(diag)
	return
}

func (e *Evaluator) docToEvalTree(d *definitions.Document) (node *evaltree.DocumentNode, diags diagnostics.Diag) {
	node = new(evaltree.DocumentNode)
	if title := d.Block.Body.Attributes["title"]; title != nil {
		pluginName := "text"
		node.AddContent(&definitions.ParsedPlugin{
			PluginName: pluginName,
			Config:     e.topLevelBlocks.DefaultConfig(definitions.BlockKindContent, pluginName),
			Invocation: definitions.NewTitle(title.Expr),
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
			call, diag := e.topLevelBlocks.ParsePlugin(plugin)
			if diags.Extend(diag) {
				continue
			}
			switch block.Type {
			case definitions.BlockKindContent:
				node.AddContent(call)
			case definitions.BlockKindData:
				node.AddData(call)
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
			var meta definitions.MetaBlock
			if diags.ExtendHcl(gohcl.DecodeBody(block.Body, nil, &meta)) {
				continue
			}
			node.AddMeta(&meta)
			origMeta = block.DefRange().Ptr()
		case definitions.BlockKindSection:
			section, diag := definitions.DefineSection(block, false)
			if diags.Extend(diag) {
				continue
			}
			parsedSection, diag := e.topLevelBlocks.ParseSection(section)
			if diags.Extend(diag) {
				continue
			}
			node.AddSection(parsedSection)
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
				},
			))
			continue
		}
	}

	return
}
