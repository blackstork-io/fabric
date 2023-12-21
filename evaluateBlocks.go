package main

import (
	"fmt"
	"slices"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/itchyny/gojq"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/json"
	"golang.org/x/exp/maps"

	"github.com/blackstork-io/fabric/pkg/diagnostics"
	"github.com/blackstork-io/fabric/pkg/jsontools"
	"github.com/blackstork-io/fabric/pkg/parexec"
	"github.com/blackstork-io/fabric/plugins/content"
	"github.com/blackstork-io/fabric/plugins/data"
)

// TODO: define magic number
var limiter = parexec.NewLimiter(5) //nolint: gomnd

// data block evaluation

type dataBlocksEvaluator struct {
	plugins map[string]any
}

type dataEvalResult struct {
	diagnostics.Diag
	Type string
	Name string
	Res  any
}

func (eb *dataBlocksEvaluator) evalBlock(db *DataBlock) (res dataEvalResult) {
	if !db.Decoded {
		res.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Undecoded block",
			Detail:   fmt.Sprintf(`%s block '%s %s "%s"' wasn't decoded`, DataBlockName, DataBlockName, db.Type, db.Name),
		})
		return
	}
	rawPlugin, found := eb.plugins[db.Type]
	if !found {
		res.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Plugin not found",
			Detail:   fmt.Sprintf("plugin %s.%s not found", DataBlockName, db.Type),
		})
		return
	}

	attrs, diags := AttrsToJSON(db.Attrs)
	if res.ExtendHcl(diags) {
		return
	}

	var err error
	res.Res, err = rawPlugin.(data.Plugin).Execute(attrs)
	if res.AppendErr(err, "Data plugin error") {
		return
	}

	res.Type = db.Type
	res.Name = db.Name
	return
}

func EvaluateDataBlocks(plugins map[string]any, dbs []DataBlock) (dict map[string]any, diags diagnostics.Diag) {
	ev := dataBlocksEvaluator{
		plugins: plugins,
	}
	// access through pe lock
	dataDict := map[string]any{}
	pe := parexec.New(
		limiter,
		func(res dataEvalResult, _ int) (cmd parexec.Command) {
			if diags.Extend(res.Diag) {
				return parexec.CmdStop
			}
			var err error
			dataDict, err = jsontools.MapSet(dataDict, []string{res.Type, res.Name}, res.Res)
			diags.AppendErr(err, "Data dict set key error")
			return
		},
	)
	parexec.MapRef(pe, dbs, ev.evalBlock)
	pe.WaitDoneAndLock()
	if diags.HasErrors() {
		return
	}
	dict = map[string]any{
		DataBlockName: dataDict,
	}
	return
}

// content block queries

type queryEvaluator struct {
	pe              parexec.Executor[diagnostics.Diag]
	dict            map[string]any
	goEvaluateQuery func(*ContentBlock)
}

func EvaluateQueries(dict map[string]any, cbs []ContentBlock) (diags diagnostics.Diag) {
	ev := queryEvaluator{
		pe: *parexec.New(
			limiter,
			func(res diagnostics.Diag, idx int) (cmd parexec.Command) {
				if diags.Extend(res) {
					return parexec.CmdStop
				}
				return
			},
		),
		dict: dict,
	}
	ev.goEvaluateQuery = parexec.GoWithArg(&ev.pe, ev.evaluateQuery)
	ev.evaluateQueries(cbs)
	ev.pe.WaitDoneAndLock()
	return
}

func (ev *queryEvaluator) evaluateQueries(cbs []ContentBlock) {
	for i := range cbs {
		cb := &cbs[i]
		if cb.Query != nil {
			ev.goEvaluateQuery(cb)
		} else {
			// no query -> no modifications -> no need to clone the dict
			cb.localDict = ev.dict
		}
		ev.evaluateQueries(cb.NestedContentBlocks)
	}
}

func (ev *queryEvaluator) evaluateQuery(cb *ContentBlock) (diags diagnostics.Diag) {
	query, err := gojq.Parse(*cb.Query)
	if err != nil {
		diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Can't parse jq query",
			Detail:   fmt.Sprintf("Error: %s Query: %s", err, *cb.Query),
		})
	}

	iter := query.Run(ev.dict)

	qRes, ok := iter.Next()
	if !ok {
		diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagWarning,
			Summary:  "Jq query returned nothing",
			Detail:   fmt.Sprintf("Query: %s", *cb.Query),
		})
		return
	}
	// decouple from the ev.dict
	cb.localDict = maps.Clone(ev.dict)
	cb.localDict["query_result"] = qRes
	return
}

// content block queries

type contentBlocksEvaluator struct {
	pe                     parexec.Executor[contentEvalResult]
	contentPlugins         map[string]any
	goEvaluateContentBlock func(*ContentBlock)
}

type contentEvalResult struct {
	diagnostics.Diag
	res string
}

func EvaluateContentBlocks(contentPlugins map[string]any, cbs []ContentBlock) (output string, diags diagnostics.Diag) {
	var orderedResult []string
	ev := contentBlocksEvaluator{
		pe: *parexec.New(
			limiter,
			func(res contentEvalResult, idx int) (cmd parexec.Command) {
				if diags.Extend(res.Diag) {
					return parexec.CmdStop
				}
				orderedResult = parexec.SetAt(orderedResult, idx, res.res)
				return
			},
		),
		contentPlugins: contentPlugins,
	}
	ev.goEvaluateContentBlock = parexec.GoWithArg(&ev.pe, ev.evaluateContentBlock)
	ev.evaluateContentBlocks(cbs)
	ev.pe.WaitDoneAndLock()
	if diags.HasErrors() {
		return
	}
	output = strings.Join(orderedResult, "\n")
	return
}

func (ev *contentBlocksEvaluator) evaluateContentBlock(cb *ContentBlock) (res contentEvalResult) {
	if !cb.Decoded {
		res.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Undecoded block",
			Detail:   fmt.Sprintf(`%s block '%s %s "%s"' wasn't decoded`, DataBlockName, DataBlockName, cb.Type, cb.Name),
		})
		return
	}
	rawPlugin, found := ev.contentPlugins[cb.Type]
	if !found {
		res.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Plugin not found",
			Detail:   fmt.Sprintf("plugin %s.%s not found", DataBlockName, cb.Type),
		})
		return
	}
	attrs, diags := AttrsToJSON(cb.Attrs)
	if res.ExtendHcl(diags) {
		return
	}

	pluginRes, err := rawPlugin.(content.Plugin).Execute(attrs, cb.localDict)
	if res.AppendErr(err, "Content plugin error") {
		return
	}
	res.res = pluginRes
	return
}

func (ev *contentBlocksEvaluator) evaluateContentBlocks(cbs []ContentBlock) {
	for i := range cbs {
		cb := &cbs[i]
		if cb.Type != "generic" {
			ev.goEvaluateContentBlock(cb)
		}
		ev.evaluateContentBlocks(cb.NestedContentBlocks)
	}
}

func (d *Decoder) FindDoc(name string) (doc *Document, diags diagnostics.Diag) {
	n := slices.IndexFunc(d.root.Documents, func(d Document) bool {
		return d.Name == name
	})
	if n == -1 {
		diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Documnent not found",
			Detail:   fmt.Sprintf("Can't find a document named '%s'", name),
		})
		return
	}
	return &d.root.Documents[n], nil
}

func (d *Decoder) Evaluate(name string) (output string, diags diagnostics.Diag) {
	doc, diag := d.FindDoc(name)
	if diags.Extend(diag) {
		return
	}

	dict, diag := EvaluateDataBlocks(d.plugins.data.plugins, doc.DataBlocks)
	if diags.Extend(diag) {
		return
	}
	diag = EvaluateQueries(dict, doc.ContentBlocks)
	if diags.Extend(diag) {
		return
	}

	output, diag = EvaluateContentBlocks(d.plugins.content.plugins, doc.ContentBlocks)
	diags.Extend(diag)
	return
}

func AttrsToJSON(attrs hcl.Attributes) (res json.SimpleJSONValue, diag hcl.Diagnostics) {
	attrsMap := make(map[string]cty.Value, len(attrs))

	for key, attr := range attrs {
		val, dgs := attr.Expr.Value(nil)
		if len(dgs) > 0 {
			for _, di := range dgs {
				di.Severity = hcl.DiagWarning
				di.Detail = fmt.Sprintf("Evaluation failed for value at key '%s': %s", key, di.Detail)
			}
			diag = diag.Extend(dgs)

			continue
		}
		attrsMap[key] = val
	}
	res = json.SimpleJSONValue{Value: cty.ObjectVal(attrsMap)}
	return
}
