package plugintest

import (
	"context"
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"

	"github.com/blackstork-io/fabric/pkg/diagnostics"
	"github.com/blackstork-io/fabric/pkg/diagnostics/diagtest"
	"github.com/blackstork-io/fabric/plugin/dataspec"
	"github.com/blackstork-io/fabric/plugin/plugindata"
)

// We have a massive amount of tests that break as soon as we add
// schemas with default values. This function is a workaround.
// It reencodes provided cty.Value to hcl text and then re-parses that text
// in accordance to spec. Ugly hack, but there's over a 100 tests in need of
// a rewrite that can't be automated with regex or similar.
// New tests should use Decode and provide a string of hcl.
func ReencodeCTY(t *testing.T, spec *dataspec.RootSpec, val cty.Value, asserts [][]diagtest.Assert) *dataspec.Block {
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

const filename = "<test-data>"

// Decodes a string (representing content of a config/data/content block)
// into cty.Value according to given spec (i.e. respecting default values)
func DecodeAndAssert(t *testing.T, spec *dataspec.RootSpec, body string, dataCtx plugindata.Map, asserts diagtest.Asserts) (val *dataspec.Block) {
	t.Helper()
	var diags diagnostics.Diag
	var fm map[string]*hcl.File
	defer func() {
		t.Helper()
		asserts.AssertMatch(t, diags, fm)
	}()
	src := []byte("block {\n")
	src = append(src, body...)
	src = append(src, "\n}\n"...)
	f, diag := hclsyntax.ParseConfig(src, filename, hcl.InitialPos)
	if diags.Extend(diag) {
		return
	}
	fm = map[string]*hcl.File{
		filename: f,
	}
	val, dgs := dataspec.DecodeAndEvalBlock(context.Background(), f.Body.(*hclsyntax.Body).Blocks[0], spec, dataCtx)
	if diags.Extend(dgs) {
		return
	}
	return
}

func Decode(t *testing.T, spec *dataspec.RootSpec, body string) (v *dataspec.Block, diags diagnostics.Diag) {
	t.Helper()
	src := []byte("block {\n")
	src = append(src, body...)
	src = append(src, "\n}\n"...)
	f, stdDiag := hclsyntax.ParseConfig(src, filename, hcl.InitialPos)
	if diags.Extend(stdDiag) {
		return
	}

	v, diag := dataspec.DecodeAndEvalBlock(context.Background(), f.Body.(*hclsyntax.Body).Blocks[0], spec, nil)
	diags.Extend(diag)
	return
}
