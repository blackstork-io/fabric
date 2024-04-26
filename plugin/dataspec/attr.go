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

	"github.com/blackstork-io/fabric/pkg/utils"
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
		if length != -1 {
			length := int64(length)
			// length-validating constraints
			if a.Constraints.Is(constraint.NonEmpty) && length == 0 {
				diags = append(diags, &hcl.Diagnostic{
					Severity: hcl.DiagError,
					Summary:  "Attribute must be non-empty",
					Detail:   fmt.Sprintf("The attribute %q can't be an empty %s.", a.Name, a.Type.FriendlyNameForConstraint()),
				})
			}
			if a.MinInclusive != cty.NilVal && a.MaxInclusive != cty.NilVal {
				min := ctyToInt(a.MinInclusive)
				max := ctyToInt(a.MaxInclusive)
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
			} else if a.MinInclusive != cty.NilVal {
				min := ctyToInt(a.MinInclusive)
				if length < min {
					diags = append(diags, &hcl.Diagnostic{
						Severity: hcl.DiagError,
						Summary:  "Attribute length is not in range",
						Detail:   fmt.Sprintf("The length of attribute %q must be >= %d", a.Name, min),
					})
				}
			} else if a.MaxInclusive != cty.NilVal {
				max := ctyToInt(a.MaxInclusive)
				if length > max {
					diags = append(diags, &hcl.Diagnostic{
						Severity: hcl.DiagError,
						Summary:  "Attribute length is not in range",
						Detail:   fmt.Sprintf("The length of attribute %q must be <= %d", a.Name, max),
					})
				}
			}
		}

		// Numeric checks:
		if val.Type().IsPrimitiveType() && val.Type() == cty.Number {
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
			if a.MinInclusive != cty.NilVal && a.MaxInclusive != cty.NilVal {
				if val.GreaterThanOrEqualTo(a.MinInclusive).And(val.LessThanOrEqualTo(a.MaxInclusive)).False() {
					diags = append(diags, &hcl.Diagnostic{
						Severity: hcl.DiagError,
						Summary:  "Attribute is not in range",
						Detail:   fmt.Sprintf("The attribute %q must be in range [%s; %s] (inclusive)", a.Name, a.MinInclusive.AsBigFloat().String(), a.MaxInclusive.AsBigFloat().String()),
					})
				}
			} else if a.MinInclusive != cty.NilVal {
				if val.GreaterThanOrEqualTo(a.MinInclusive).False() {
					diags = append(diags, &hcl.Diagnostic{
						Severity: hcl.DiagError,
						Summary:  "Attribute is not in range",
						Detail:   fmt.Sprintf("The attribute %q must be >= %s", a.Name, a.MinInclusive.AsBigFloat().String()),
					})
				}
			} else if a.MaxInclusive != cty.NilVal {
				if val.LessThanOrEqualTo(a.MaxInclusive).False() {
					diags = append(diags, &hcl.Diagnostic{
						Severity: hcl.DiagError,
						Summary:  "Attribute is not in range",
						Detail:   fmt.Sprintf("The attribute %q must be <= %s", a.Name, a.MaxInclusive.AsBigFloat().String()),
					})
				}
			}
		}
	}
	if !a.OneOf.IsEmpty() {
		found := false
		for _, possibleVal := range a.OneOf {
			if possibleVal.Equals(val).True() {
				found = true
				break
			}
		}
		if !found {
			diags = append(diags, &hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Attribute is not one of the allowed values",
				Detail:   fmt.Sprintf("The attribute %q must be one of: %s", a.Name, a.OneOf),
			})
		}
	}
	return
}

func joinSummaries(diags hcl.Diagnostics) string {
	return strings.Join(
		utils.FnMap(diags, func(diag *hcl.Diagnostic) string { return diag.Summary }),
		"; ",
	)
}

func (a *AttrSpec) ValidateSpec() (errs []string) {
	if a.Constraints.Is(constraint.Required) {
		//
		// if a.ExampleVal == cty.NilVal {
		// 	errs = append(errs, fmt.Sprintf("Missing example value on required attibute %q", a.Name))
		// }
		if a.DefaultVal != cty.NilVal {
			errs = append(errs, fmt.Sprintf("Default value is specified for the required attribute %q = %s", a.Name, a.DefaultVal.GoString()))
		}
	}

	if a.Constraints.Is(constraint.Integer) && !(a.Type.Equals(cty.Number)) {
		errs = append(errs, fmt.Sprintf("Integer constraint is specified for non-numeric attribute %q", a.Name))
	}

	if a.MinInclusive != cty.NilVal {
		if (a.Type.IsPrimitiveType() && a.Type == cty.Bool) || (a.Type.IsCapsuleType()) {
			errs = append(errs, fmt.Sprintf("MinValInclusive can't be specified for %s %q", a.Type.FriendlyName(), a.Name))
		}
		if !(a.MinInclusive.Type().IsPrimitiveType() && a.MinInclusive.Type() == cty.Number) {
			errs = append(errs, fmt.Sprintf("MinValInclusive specified for %q must be a number, not %s", a.Name, a.MinInclusive.Type().FriendlyName()))
		}

		if !(a.Type.Equals(cty.Number)) {
			// Min is length, must be an int >=0
			int, acc := a.MinInclusive.AsBigFloat().Int64()
			if acc != big.Exact {
				errs = append(errs, fmt.Sprintf("MinValInclusive specified for %q must be an integer", a.Name))
			}
			if int < 0 {
				errs = append(errs, fmt.Sprintf("MinValInclusive specified for %q must be >= 0", a.Name))
			}
		}

	}
	if a.MaxInclusive != cty.NilVal {
		if (a.Type.IsPrimitiveType() && a.Type == cty.Bool) || (a.Type.IsCapsuleType()) {
			errs = append(errs, fmt.Sprintf("MaxValInclusive can't be specified for %s %q", a.Type.FriendlyName(), a.Name))
		}
		if !(a.MaxInclusive.Type().IsPrimitiveType() && a.MaxInclusive.Type() == cty.Number) {
			errs = append(errs, fmt.Sprintf("MaxValInclusive specified for %q must be a number, not %s", a.Name, a.MaxInclusive.Type().FriendlyName()))
		}
		if !(a.Type.Equals(cty.Number)) {
			// Max is length, must be an int >=0
			int, acc := a.MaxInclusive.AsBigFloat().Int64()
			if acc != big.Exact {
				errs = append(errs, fmt.Sprintf("MaxValInclusive specified for %q must be an integer", a.Name))
			}
			if int < 0 {
				errs = append(errs, fmt.Sprintf("MaxValInclusive specified for %q must be >= 0", a.Name))
			}
		}
	}
	if len(errs) == 0 && a.MinInclusive != cty.NilVal && a.MaxInclusive != cty.NilVal {
		// no errors - values are numbers and can be compared
		if a.MinInclusive.LessThanOrEqualTo(a.MaxInclusive).False() {
			errs = append(errs, fmt.Sprintf("%q: MaxValInclusive must be <= MaxValInclusive", a.Name))
		}
	}

	if len(errs) == 0 {
		if a.DefaultVal != cty.NilVal {
			diags := a.ValidateValue(a.DefaultVal)
			if diags.HasErrors() {
				errs = append(errs,
					fmt.Sprintf("Default value for attribute %q (%s) failed validation: %s",
						a.Name,
						a.DefaultVal.GoString(),
						joinSummaries(diags),
					))
			}
		}
		if a.ExampleVal != cty.NilVal {
			diags := a.ValidateValue(a.ExampleVal)
			if diags.HasErrors() {
				errs = append(errs,
					fmt.Sprintf("Example value for attribute %q (%s) failed validation: %s",
						a.Name,
						a.ExampleVal.GoString(),
						joinSummaries(diags),
					))
			}
		}
	}
	return
}
