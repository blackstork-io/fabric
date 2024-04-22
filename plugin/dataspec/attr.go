package dataspec

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
)

// AttrSpec represents the attribute value (hcldec.AttrSpec).
type AttrSpec struct {
	Name       string
	Type       cty.Type
	DefaultVal cty.Value
	ExampleVal cty.Value
	Doc        string
	Required   bool
}

func (a *AttrSpec) KeyForObjectSpec() string {
	return a.Name
}

func (a *AttrSpec) getSpec() Spec {
	return a
}

func (a *AttrSpec) DocComment() hclwrite.Tokens {
	tokens := comment(nil, a.Doc)
	if len(tokens) != 0 {
		tokens = appendCommentNewLine(tokens)
	}

	if a.Required {
		tokens = comment(
			tokens,
			fmt.Sprintf("Required %s. For example:", a.Type.FriendlyNameForConstraint()),
		)
	} else {
		if a.ExampleVal != cty.NilVal {
			f := hclwrite.NewEmptyFile()
			f.Body().SetAttributeValue(a.Name, a.ExampleVal)
			tokens = comment(tokens, "For example:\n"+string(hclwrite.Format(f.Bytes())))
			tokens = appendCommentNewLine(tokens)
		}
		tokens = comment(
			tokens,
			fmt.Sprintf("Optional %s. Default value:", a.Type.FriendlyNameForConstraint()),
		)
	}
	return tokens
}

func (a *AttrSpec) WriteDoc(w *hclwrite.Body) {
	// write out documnetation
	w.AppendUnstructuredTokens(a.DocComment())

	// write the attribute
	var val cty.Value
	if a.Required {
		val = a.ExampleVal
		if val.IsNull() {
			val = exampleValueForType(a.Type)
		}
	} else {
		val = a.DefaultVal
	}

	w.SetAttributeValue(a.Name, val)
}

func (a *AttrSpec) HcldecSpec() hcldec.Spec {
	res := &hcldec.AttrSpec{
		Name:     a.Name,
		Type:     a.Type,
		Required: a.Required,
	}
	if a.Required {
		return &hcldec.ValidateSpec{
			Wrapped: res,
			Func: func(value cty.Value) hcl.Diagnostics {
				if value.IsNull() {
					return []*hcl.Diagnostic{{
						Severity: hcl.DiagError,
						Summary:  "Argument must be non-null",
						Detail:   fmt.Sprintf("The argument %q was either not defined or is null.", res.Name),
					}}
				}
				return nil
			},
		}
	} else if !a.DefaultVal.IsNull() {
		return &hcldec.DefaultSpec{
			Primary: res,
			Default: &hcldec.LiteralSpec{
				Value: a.DefaultVal,
			},
		}
	}

	return res
}

func (a *AttrSpec) Validate() (errs []string) {
	if a.Required {
		if a.ExampleVal == cty.NilVal {
			errs = append(errs, fmt.Sprintf("Missing example value on required attibute %q", a.Name))
		}
		if a.DefaultVal != cty.NilVal {
			errs = append(errs, fmt.Sprintf("Default value is specified for the required attribute %q = %s", a.Name, a.DefaultVal.GoString()))
		}
	}
	return
}
