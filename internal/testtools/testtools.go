package testtools

import (
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"

	"github.com/blackstork-io/fabric/pkg/diagnostics"
	"github.com/blackstork-io/fabric/plugin/dataspec"
)

// We have a massive amount of tests that break as soon as we add
// schemas with default values. This function is a workaround.
// It reencodes provided cty.Value to hcl text and then re-parses that text
// in accordance to spec. Ugly hack, but there's over a 100 tests in need of
// a rewrite that can't be automated with regex or similar.
// New tests should use Decode and provide a string of hcl.
func ReencodeCTY(t *testing.T, spec dataspec.RootSpec, val cty.Value, asserts [][]Assert) cty.Value {
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
	return Decode(t, spec, string(hclwrite.Format(f.Bytes())), asserts)
}

const filename = "<inline-data>"

// We have a massive amount of thests that break as soon as we add
// schemas with default values. These functions are a workaround
// Use the Decode in the newer tests

// Decodes a string (representing content of a config/data/content block)
// into cty.Value according to given spec (i.e. respecting default values)
func Decode(t *testing.T, spec dataspec.RootSpec, body string, asserts [][]Assert) (v cty.Value) {
	t.Helper()
	var f *hcl.File
	var diags diagnostics.Diag
	defer func() {
		t.Helper()
		var fm map[string]*hcl.File
		if f != nil {
			fm = map[string]*hcl.File{
				filename: f,
			}
		}
		CompareDiags(t, fm, diags, asserts)
	}()
	var diag hcl.Diagnostics
	f, diag = hclsyntax.ParseConfig([]byte(body), filename, hcl.InitialPos)
	if diags.ExtendHcl(diag) {
		return
	}

	v, diag = hcldec.Decode(f.Body, spec.HcldecSpec(), nil)
	diags = append(diags, diag...)
	return
}
