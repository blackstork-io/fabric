package dataspec

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"

	"github.com/blackstork-io/fabric/pkg/diagnostics"
	"github.com/blackstork-io/fabric/plugin/dataspec/constraint"
)

// AttrSpec represents the attribute value (hcldec.AttrSpec).
type AttrSpec struct {
	Name       string
	Type       cty.Type
	DefaultVal cty.Value
	ExampleVal cty.Value
	Doc        string
	// TODO: replace by constraints
	Required    bool
	Constraints constraint.Constraints
	// If specified - value must be on of specified values
	OneOf constraint.OneOf
	// For numbers - min value; for collections - min number of elements; for strings - min length
	MinInclusive cty.Value
	// For numbers - max value; for collections - max number of elements; for strings - max length
	MaxInclusive cty.Value
	// If specified – a deprication warning would appear if an attribute is specified is non-null
	Depricated string
}

func (a *AttrSpec) computeMinInclusive() cty.Value {
	if !a.Constraints.Is(constraint.NonEmpty) || (a.Type.IsPrimitiveType() && a.Type == cty.Number) {
		return a.MinInclusive
	}
	// we have constraint.NonEmpty constraint on a collection type
	if a.MinInclusive.IsNull() || a.MinInclusive.LessThan(cty.NumberIntVal(1)).True() {
		return cty.NumberIntVal(1)
	}
	return a.MinInclusive
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

	var buf strings.Builder
	if a.Constraints.Is(constraint.Required) {
		buf.WriteString("Required ")
	} else {
		buf.WriteString("Optional ")
	}
	if a.Constraints.Is(constraint.Integer) {
		buf.WriteString("integer")
	} else {
		buf.WriteString(a.Type.FriendlyNameForConstraint())
	}
	buf.WriteString(".\n")

	if !a.OneOf.IsEmpty() {
		buf.WriteString("Must be one of: ")
		buf.WriteString(a.OneOf.String())
		buf.WriteString("\n")

	}
	min := a.computeMinInclusive()
	max := a.MaxInclusive

	if !min.IsNull() && !max.IsNull() {
		if a.Type.IsPrimitiveType() && a.Type == cty.Number {
			fmt.Fprintf(&buf, "Must be between %s and %s (inclusive)\n", min.AsBigFloat().String(), max.AsBigFloat().String())
		} else {
			min, _ := min.AsBigFloat().Uint64()
			max, _ := max.AsBigFloat().Uint64()
			if min == max {
				fmt.Fprintf(&buf, "Must have a length of %d\n", min)
			} else {
				fmt.Fprintf(&buf, "Must have a length between %d and %d (inclusive)\n", min, max)
			}
		}
	} else if !min.IsNull() {
		if a.Type.IsPrimitiveType() && a.Type == cty.Number {
			fmt.Fprintf(&buf, "Must be >= %s", min.AsBigFloat().String())
		} else {
			min, _ := min.AsBigFloat().Uint64()
			fmt.Fprintf(&buf, "Must have a length of at least %d", min)
		}
	} else if !max.IsNull() {
		if a.Type.IsPrimitiveType() && a.Type == cty.Number {
			fmt.Fprintf(&buf, "Must be <= %s", max.AsBigFloat().String())
		} else {
			max, _ := max.AsBigFloat().Uint64()
			fmt.Fprintf(&buf, "Must have a length of at most %d", max)
		}
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

var trimSpace = function.New(&function.Spec{
	Description: "Trim string, noop for other types",
	Params: []function.Parameter{{
		Name:        "string",
		Description: "string to be trimmed. Does noting if passed a non-string",
		Type:        cty.DynamicPseudoType,
		AllowNull:   true,
	}},
	Type: func(args []cty.Value) (cty.Type, error) {
		return args[0].Type(), nil
	},
	Impl: func(args []cty.Value, retType cty.Type) (cty.Value, error) {
		if !args[0].IsNull() && args[0].Type().IsPrimitiveType() && args[0].Type() == cty.String {
			return cty.StringVal(strings.TrimSpace(args[0].AsString())), nil
		}
		return args[0], nil
	},
})

func (a *AttrSpec) HcldecSpec() (res hcldec.Spec) {
	res = &hcldec.AttrSpec{
		Name:     a.Name,
		Type:     a.Type,
		Required: a.Constraints.Is(constraint.Required),
	}
	if !a.DefaultVal.IsNull() {
		res = &hcldec.DefaultSpec{
			Primary: res,
			Default: &hcldec.LiteralSpec{
				Value: a.DefaultVal,
			},
		}
	}

	if a.Constraints.Is(constraint.TrimSpace) {
		res = &hcldec.TransformFuncSpec{
			Wrapped: res,
			Func:    trimSpace,
		}
	}
	return &hcldec.ValidateSpec{
		Wrapped: res,
		Func:    a.ValidateValue,
	}
}

func ctyToInt(val cty.Value) int64 {
	i, acc := val.AsBigFloat().Int64()
	if acc != big.Exact {
		panic(fmt.Sprintf("%s is not an exact integer", val.GoString()))
	}
	return i
}

func (a *AttrSpec) ValidateValue(val cty.Value) (diags hcl.Diagnostics) {
	if val.IsNull() {
		if a.Constraints.Is(constraint.NonNull) {
			diags = append(diags, &hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Attribute must be non-null",
				Detail:   fmt.Sprintf("The attribute %q was either not defined or is null.", a.Name),
			})
		}
	} else {
		// non null
		if a.Depricated != "" {
			diags = append(diags, &hcl.Diagnostic{
				Severity: hcl.DiagWarning,
				Summary:  "Deprecated attribute",
				Detail:   fmt.Sprintf("The attribute %q is deprecated: %s", a.Name, a.Depricated),
			})
		}

		// Length checks:
		length := -1
		if val.Type().IsCollectionType() || val.Type().IsTupleType() {
			length = val.LengthInt()
		} else if val.Type().IsPrimitiveType() && val.Type() == cty.String {
			length = len(val.AsString())
		}
		min := a.computeMinInclusive()
		max := a.MaxInclusive
		if length != -1 {
			length := int64(length)

			// length-validating constraints

			if !min.IsNull() && !max.IsNull() {
				min := ctyToInt(min)
				max := ctyToInt(max)
				if !(min <= length && length <= max) {
					if min == max {
						diags = append(diags, &hcl.Diagnostic{
							Severity: hcl.DiagError,
							Summary:  "Attribute length is not in range",
							Detail:   fmt.Sprintf("The length of attribute %q must be exactly %d", a.Name, min),
						})
					} else {
						diags = append(diags, &hcl.Diagnostic{
							Severity: hcl.DiagError,
							Summary:  "Attribute length is not in range",
							Detail:   fmt.Sprintf("The length of attribute %q must be in range [%d; %d] (inclusive)", a.Name, min, max),
						})
					}
				}
			} else if !min.IsNull() {
				min := ctyToInt(min)
				if length < min {
					diags = append(diags, &hcl.Diagnostic{
						Severity: hcl.DiagError,
						Summary:  "Attribute length is not in range",
						Detail:   fmt.Sprintf("The length of attribute %q must be >= %d", a.Name, min),
					})
				}
			} else if !max.IsNull() {
				max := ctyToInt(max)
				if length > max {
					diags = append(diags, &hcl.Diagnostic{
						Severity: hcl.DiagError,
						Summary:  "Attribute length is not in range",
						Detail:   fmt.Sprintf("The length of attribute %q must be <= %d", a.Name, max),
					})
				}
			}
		} else if val.Type().IsPrimitiveType() && val.Type() == cty.Number {
			// Numeric checks:
			if a.Constraints.Is(constraint.Integer) {
				_, acc := val.AsBigFloat().Int64()
				if acc != big.Exact {
					diags = append(diags, &hcl.Diagnostic{
						Severity: hcl.DiagError,
						Summary:  "Attribute must be an integer",
						Detail:   fmt.Sprintf("The attribute %q must be an integer", a.Name),
					})
				}
			}
			// Range checks:
			if !min.IsNull() && !max.IsNull() {
				if val.GreaterThanOrEqualTo(min).And(val.LessThanOrEqualTo(max)).False() {
					diags = append(diags, &hcl.Diagnostic{
						Severity: hcl.DiagError,
						Summary:  "Attribute is not in range",
						Detail:   fmt.Sprintf("The attribute %q must be in range [%s; %s] (inclusive)", a.Name, min.AsBigFloat().String(), max.AsBigFloat().String()),
					})
				}
			} else if !min.IsNull() {
				if val.GreaterThanOrEqualTo(min).False() {
					diags = append(diags, &hcl.Diagnostic{
						Severity: hcl.DiagError,
						Summary:  "Attribute is not in range",
						Detail:   fmt.Sprintf("The attribute %q must be >= %s", a.Name, min.AsBigFloat().String()),
					})
				}
			} else if !max.IsNull() {
				if val.LessThanOrEqualTo(max).False() {
					diags = append(diags, &hcl.Diagnostic{
						Severity: hcl.DiagError,
						Summary:  "Attribute is not in range",
						Detail:   fmt.Sprintf("The attribute %q must be <= %s", a.Name, max.AsBigFloat().String()),
					})
				}
			}
		}
	}
	if !a.OneOf.Validate(val) {
		diags = append(diags, &hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Attribute is not one of the allowed values",
			Detail:   fmt.Sprintf("The attribute %q must be one of: %s", a.Name, a.OneOf),
		})
	}
	return
}

func (a *AttrSpec) ValidateSpec() (errs diagnostics.Diag) {
	if a.Constraints.Is(constraint.Required) {
		if a.ExampleVal == cty.NilVal {
			errs.AddWarn(fmt.Sprintf("Missing example value on required attibute %q", a.Name), "")
		}
		if a.DefaultVal != cty.NilVal {
			errs.Add(fmt.Sprintf("Default value is specified for the required attribute %q = %s", a.Name, a.DefaultVal.GoString()), "")
		}
	}

	if a.Constraints.Is(constraint.Integer) && !(a.Type.Equals(cty.Number)) {
		errs.Add(fmt.Sprintf("Integer constraint is specified for non-numeric attribute %q", a.Name), "")
	}
	min := a.MinInclusive

	max := a.MaxInclusive
	skipMinMaxRelativeCheck := false
	for _, v := range []struct {
		name string
		val  cty.Value
	}{{"MinInclusive", min}, {"MaxInclusive", max}} {
		if v.val == cty.NilVal {
			continue
		}

		if (a.Type.IsPrimitiveType() && a.Type == cty.Bool) || (a.Type.IsCapsuleType()) {
			errs.Add(fmt.Sprintf("%s can't be specified for %s %q", v.name, a.Type.FriendlyName(), a.Name), "")
			skipMinMaxRelativeCheck = true
			continue
		}
		if !(v.val.Type().IsPrimitiveType() && v.val.Type() == cty.Number) {
			errs.Add(fmt.Sprintf("%s specified for %q must be a number, not %s", v.name, a.Name, v.val.Type().FriendlyName()), "")
			skipMinMaxRelativeCheck = true
			continue
		}
		if v.val.IsNull() {
			errs.Add(fmt.Sprintf("%s specified for %q must be non-null", v.name, a.Name), "")
			skipMinMaxRelativeCheck = true
			continue
		}

		if !(a.Type.Equals(cty.Number)) {
			// Min is length, must be an num >=0
			num, acc := v.val.AsBigFloat().Int64()
			if acc != big.Exact {
				errs.Add(fmt.Sprintf("%s specified for %q must be an integer", v.name, a.Name), "")
			}
			if num < 0 {
				errs.Add(fmt.Sprintf("%s specified for %q must be >= 0", v.name, a.Name), "")
			}
		}
	}

	if !skipMinMaxRelativeCheck && !min.IsNull() && !max.IsNull() {
		// no errors - values are numbers and can be compared
		if min.LessThanOrEqualTo(max).False() {
			errs.Add(fmt.Sprintf("%q: MinInclusive must be <= MaxInclusive", a.Name), "")
		}
	}

	if len(errs) == 0 {
		if a.DefaultVal != cty.NilVal {
			diags := a.ValidateValue(a.DefaultVal)
			prefix := fmt.Sprintf("Default value for attribute %q: ", a.Name)
			for _, d := range diags {
				if d.Severity == hcl.DiagError {
					d.Summary = prefix + d.Summary
					errs.Append(d)
				}
			}
		}
		if a.ExampleVal != cty.NilVal {
			diags := a.ValidateValue(a.ExampleVal)
			prefix := fmt.Sprintf("Example value for attribute %q: ", a.Name)
			for _, d := range diags {
				if d.Severity == hcl.DiagError {
					d.Summary = prefix + d.Summary
					errs.Append(d)
				}
			}
		}
	}
	return
}
