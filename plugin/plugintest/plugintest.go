package plugintest

import (
	"context"
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"

	"github.com/blackstork-io/fabric/parser/evaluation"
	"github.com/blackstork-io/fabric/pkg/diagnostics"
	"github.com/blackstork-io/fabric/pkg/diagnostics/diagtest"
	"github.com/blackstork-io/fabric/plugin"
	"github.com/blackstork-io/fabric/plugin/dataspec"
)

// We have a massive amount of tests that break as soon as we add
// schemas with default values. This function is a workaround.
// It reencodes provided cty.Value to hcl text and then re-parses that text
// in accordance to spec. Ugly hack, but there's over a 100 tests in need of
// a rewrite that can't be automated with regex or similar.
// New tests should use Decode and provide a string of hcl.
func ReencodeCTY(t *testing.T, spec dataspec.RootSpec, val cty.Value, asserts [][]diagtest.Assert) cty.Value {
	t.Helper()
	ty := val.Type()
	if !(ty.IsMapType() || ty.IsObjectType()) {
		panic("Can't handle type " + ty.FriendlyName())
	}
	f := hclwrite.NewEmptyFile()
	b := f.Body()
	for it := val.ElementIterator(); it.Next(); {
		k, v := it.Element()
		b.SetAttributeValue(k.AsString(), v)
	}
	return DecodeAndAssert(t, spec, string(hclwrite.Format(f.Bytes())), nil, asserts)
}

const filename = "<inline-data>"

// We have a massive amount of tests that break as soon as we add
// schemas with default values. These functions are a workaround
// Use the Decode in the newer tests

// Decodes a string (representing content of a config/data/content block)
// into cty.Value according to given spec (i.e. respecting default values)
func DecodeAndAssert(t *testing.T, spec dataspec.RootSpec, body string, dataCtx plugin.MapData, asserts diagtest.Asserts) (v cty.Value) {
	t.Helper()
	var diags diagnostics.Diag
	var fm map[string]*hcl.File
	defer func() {
		t.Helper()
		asserts.AssertMatch(t, diags, fm)
	}()
	f, diag := hclsyntax.ParseConfig([]byte(body), filename, hcl.InitialPos)
	if diags.Extend(diag) {
		return
	}
	fm = map[string]*hcl.File{
		filename: f,
	}
	val, dgs := dataspec.Decode(f.Body, spec, evaluation.EvalContext())
	if diags.Extend(dgs) {
		return
	}
	v, dgs = plugin.CustomEvalTransform(context.Background(), dataCtx, val)
	if diags.Extend(dgs) {
		return
	}
	return
}

func Decode(t *testing.T, spec dataspec.RootSpec, body string) (v cty.Value, diags diagnostics.Diag) {
	t.Helper()
	f, stdDiag := hclsyntax.ParseConfig([]byte(body), filename, hcl.InitialPos)
	if diags.Extend(stdDiag) {
		return
	}

	v, diag := dataspec.Decode(f.Body, spec, evaluation.EvalContext())
	diags.Extend(diag)
	return
}
