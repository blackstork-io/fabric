package parser

import (
	"fmt"
	"maps"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/itchyny/gojq"
	"github.com/zclconf/go-cty/cty"

	"github.com/blackstork-io/fabric/parser/definitions"
	"github.com/blackstork-io/fabric/pkg/diagnostics"
	"github.com/blackstork-io/fabric/pkg/jsontools"
	"github.com/blackstork-io/fabric/pkg/utils"
)

// Evaluates a chosen document

type contentCall struct {
	*definitions.ParsedPlugin
	query string
}

type Evaluator struct {
	caller         PluginCaller
	contentCalls   []*definitions.ParsedPlugin
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

func (e *Evaluator) evaluateQuery(call *definitions.ParsedPlugin) (context map[string]any, diags diagnostics.Diag) {
	context = e.context
	body := call.Invocation.GetBody()
	queryAttr, found := body.Attributes["query"]
	if !found {
		return
	}
	val, newBody, dgs := hcldec.PartialDecode(body, &hcldec.ObjectSpec{
		"query": &hcldec.AttrSpec{
			Name:     "query",
			Type:     cty.String,
			Required: true,
		},
	}, nil)
	call.Invocation.SetBody(utils.ToHclsyntaxBody(newBody))
	if diags.ExtendHcl(dgs) {
		return
	}
	query := val.GetAttr("query").AsString()
	q, err := gojq.Parse(query)
	if err != nil {
		diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Failed to parse the query",
			Detail:   err.Error(),
			Subject:  &queryAttr.SrcRange,
		})
		return
	}

	code, err := gojq.Compile(q)
	if err != nil {
		diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Failed to compile the query",
			Detail:   err.Error(),
			Subject:  &queryAttr.SrcRange,
		})
		return
	}
	queryResultIter := code.Run(context)
	queryResult, ok := queryResultIter.Next()
	if ok {
		context = maps.Clone(context)
		context["query_result"] = queryResult
	}
	return
}

func (e *Evaluator) EvaluateDocument(d *definitions.DocumentOrSection) (output string, diags diagnostics.Diag) {
	// sections are basically documents
	diags = e.evaluateSectionOrDocument(d)
	if diags.HasErrors() {
		return
	}

	results := make([]string, 0, len(e.contentCalls))
	for _, call := range e.contentCalls {
		context, diag := e.evaluateQuery(call)
		if diags.Extend(diag) {
			// query failed, but context is always valid
			// TODO: #28 #29
		}
		result, diag := e.caller.CallContent(call.PluginName, call.Config, call.Invocation, context)
		if diags.Extend(diag) {
			// XXX: What to do if we have errors while executing content blocks?
			// just skipping the value for now...
			continue
		}
		results = append(results, result)
		// TODO: Here's the place to implement local context #17
		// However I think we need to rework it a bit before done
	}
	output = strings.Join(results, "\n")
	return
}

func (e *Evaluator) evaluateSectionOrDocument(d *definitions.DocumentOrSection) (diags diagnostics.Diag) {
	if title := d.Block.Body.Attributes["title"]; title != nil {
		e.contentCalls = append(e.contentCalls, &definitions.ParsedPlugin{
			PluginName: "text",
			Config:     nil,
			Invocation: definitions.NewTitle(title),
		})
	}

	var origMeta *hcl.Range

	var validChildren []string

	if d.IsDocument() {
		validChildren = []string{
			definitions.BlockKindContent,
			definitions.BlockKindData,
			definitions.BlockKindMeta,
			definitions.BlockKindSection,
		}
	} else {
		validChildren = []string{
			definitions.BlockKindContent,
			definitions.BlockKindMeta,
			definitions.BlockKindSection,
		}
	}

	for _, block := range d.Block.Body.Blocks {
		switch block.Type {
		case definitions.BlockKindData:
			if !d.IsDocument() {
				// Deny data blocks in sections
				diags.Append(definitions.NewNestingDiag(
					d.Block.Type,
					block,
					d.Block.Body,
					validChildren,
				))
				continue
			}
			fallthrough
		case definitions.BlockKindContent:
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
				var err error
				e.context, err = jsontools.MapSet(e.context, []string{
					definitions.BlockKindData,
					call.PluginName,
					call.BlockName,
				}, res)
				diags.AppendErr(err, "Failed to save data plugin result")
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
			d.Meta = &meta
			origMeta = block.DefRange().Ptr()
		case definitions.BlockKindSection:
			section, diag := definitions.DefineSectionOrDocument(block, false)
			if diags.Extend(diag) {
				continue
			}
			diag = e.evaluateSectionOrDocument(section)
			if diags.Extend(diag) {
				continue
			}
		default:
			diags.Append(definitions.NewNestingDiag(
				d.Block.Type,
				block,
				d.Block.Body,
				validChildren,
			))
		}
	}

	return
}
