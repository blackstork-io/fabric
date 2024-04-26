package dataspec_test

import (
	"testing"

	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/stretchr/testify/assert"
	"github.com/zclconf/go-cty/cty"

	"github.com/blackstork-io/fabric/internal/testtools"
	"github.com/blackstork-io/fabric/plugin/dataspec"
	"github.com/blackstork-io/fabric/plugin/dataspec/constraint"
)

func TestValidation(t *testing.T) {
	// t.Parallel()
	for _, tc := range []struct {
		name         string
		obj          *dataspec.AttrSpec
		inputVal     cty.Value
		expectedVal  cty.Value
		asserts      [][]testtools.Assert
		schemaErrors []string
	}{
		{
			name: "basicAttribute",
			obj: &dataspec.AttrSpec{
				Name: "test",
				Type: cty.String,
			},
			expectedVal: cty.NullVal(cty.String),
		},
		{
			name: "requiredAttributeMissing",
			obj: &dataspec.AttrSpec{
				Name:        "test",
				Type:        cty.String,
				ExampleVal:  cty.StringVal("test"),
				Constraints: constraint.Required,
			},
			asserts: [][]testtools.Assert{{
				testtools.IsError,
				testtools.SummaryContains("Missing required argument"),
			}},
		},
		{
			name: "requiredAttributeMissingNonnil",
			obj: &dataspec.AttrSpec{
				Name:        "test",
				Type:        cty.String,
				ExampleVal:  cty.StringVal("test"),
				Constraints: constraint.Required | constraint.NonNull,
			},
			asserts: [][]testtools.Assert{
				{
					testtools.IsError,
					testtools.SummaryContains("Missing required argument"),
				},
				{
					testtools.IsError,
					testtools.SummaryContains("Attribute must be non-null"),
				},
			},
		},
		{
			name: "Null_for_Nonnull_value",
			obj: &dataspec.AttrSpec{
				Name:        "test",
				Type:        cty.String,
				ExampleVal:  cty.StringVal("test"),
				Constraints: constraint.Required | constraint.NonNull,
			},
			inputVal: cty.NullVal(cty.String),
			asserts: [][]testtools.Assert{
				{
					testtools.IsError,
					testtools.SummaryContains("Attribute must be non-null"),
				},
			},
		},
		{
			name: "empty_for_nonnull_value",
			obj: &dataspec.AttrSpec{
				Name:        "test",
				Type:        cty.String,
				ExampleVal:  cty.StringVal("test"),
				Constraints: constraint.Required | constraint.NonNull,
			},
			inputVal:    cty.StringVal(""),
			expectedVal: cty.StringVal(""),
		},
		{
			name: "Nonempty",
			obj: &dataspec.AttrSpec{
				Name:        "test",
				Type:        cty.String,
				ExampleVal:  cty.StringVal("test"),
				Constraints: constraint.Required | constraint.NonEmpty,
			},
			inputVal:    cty.StringVal(""),
			expectedVal: cty.StringVal(""),
			asserts: [][]testtools.Assert{
				{
					testtools.IsError,
					testtools.SummaryContains("Attribute must be non-empty"),
				},
			},
		},
		{
			name: "Nonempty_spaces",
			obj: &dataspec.AttrSpec{
				Name:        "test",
				Type:        cty.String,
				ExampleVal:  cty.StringVal("test"),
				Constraints: constraint.Required | constraint.NonEmpty,
			},
			inputVal:    cty.StringVal("    "),
			expectedVal: cty.StringVal("    "),
		},
		{
			name: "Trimmed_Nonempty_spaces",
			obj: &dataspec.AttrSpec{
				Name:        "test",
				Type:        cty.String,
				ExampleVal:  cty.StringVal("test"),
				Constraints: constraint.Required | constraint.TrimmedNonEmpty,
			},
			inputVal: cty.StringVal("   "),
			asserts: [][]testtools.Assert{
				{
					testtools.IsError,
					testtools.SummaryContains("Attribute must be non-empty"),
				},
			},
		},
		{
			name: "Length_check_min",
			obj: &dataspec.AttrSpec{
				Name:         "test",
				Type:         cty.String,
				ExampleVal:   cty.StringVal("qwertyqwerty"),
				Constraints:  constraint.RequiredMeaningfull,
				MinInclusive: cty.NumberIntVal(10),
			},
			inputVal: cty.StringVal("hello"),
			asserts: [][]testtools.Assert{
				{
					testtools.IsError,
					testtools.SummaryContains("Attribute length is not in range"),
					testtools.DetailContains(">=", "10"),
				},
			},
		},
		{
			name: "Length_check_max",
			obj: &dataspec.AttrSpec{
				Name:         "test",
				Type:         cty.String,
				ExampleVal:   cty.StringVal("12"),
				Constraints:  constraint.RequiredMeaningfull,
				MaxInclusive: cty.NumberIntVal(2),
			},
			inputVal: cty.StringVal("hello"),
			asserts: [][]testtools.Assert{
				{
					testtools.IsError,
					testtools.SummaryContains("Attribute length is not in range"),
					testtools.DetailContains("<=", "2"),
				},
			},
		},
		{
			name: "Length_check_range",
			obj: &dataspec.AttrSpec{
				Name:         "test",
				Type:         cty.String,
				ExampleVal:   cty.StringVal("t"),
				Constraints:  constraint.RequiredMeaningfull,
				MinInclusive: cty.NumberIntVal(1),
				MaxInclusive: cty.NumberIntVal(3),
			},
			inputVal: cty.StringVal("hello"),
			asserts: [][]testtools.Assert{
				{
					testtools.IsError,
					testtools.SummaryContains("Attribute length is not in range"),
					testtools.DetailContains("1", "3"),
				},
			},
		},
		{
			name: "Length_check_range_ok",
			obj: &dataspec.AttrSpec{
				Name:         "test",
				Type:         cty.String,
				ExampleVal:   cty.StringVal("test1"),
				Constraints:  constraint.RequiredMeaningfull,
				MinInclusive: cty.NumberIntVal(5),
				MaxInclusive: cty.NumberIntVal(5),
			},
			inputVal:    cty.StringVal("hello"),
			expectedVal: cty.StringVal("hello"),
		},
		{
			name: "Length_check_range_exact",
			obj: &dataspec.AttrSpec{
				Name:         "test",
				Type:         cty.String,
				Constraints:  constraint.RequiredMeaningfull,
				ExampleVal:   cty.StringVal("test1"),
				MinInclusive: cty.NumberIntVal(5),
				MaxInclusive: cty.NumberIntVal(5),
			},
			inputVal: cty.StringVal("hello_123"),
			asserts: [][]testtools.Assert{
				{
					testtools.IsError,
					testtools.SummaryContains("Attribute length is not in range"),
					testtools.DetailContains("exactly", "5"),
				},
			},
		},
		{
			name: "Length_inverted",
			obj: &dataspec.AttrSpec{
				Name:         "test",
				Type:         cty.String,
				ExampleVal:   cty.StringVal("t"),
				Constraints:  constraint.RequiredMeaningfull,
				MinInclusive: cty.NumberIntVal(10),
				MaxInclusive: cty.NumberIntVal(1),
			},
			schemaErrors: []string{`"test": MaxValInclusive must be <= MaxValInclusive`},
			inputVal:     cty.StringVal("hello_123"),
		},

		{
			name: "MinLength_<0",
			obj: &dataspec.AttrSpec{
				Name:         "test",
				Type:         cty.String,
				ExampleVal:   cty.StringVal("t"),
				Constraints:  constraint.RequiredMeaningfull,
				MinInclusive: cty.NumberIntVal(-1),
			},
			schemaErrors: []string{`MinValInclusive specified for "test" must be >= 0`},
			inputVal:     cty.StringVal("hello_123"),
		},
		{
			name: "MaxLength_<0",
			obj: &dataspec.AttrSpec{
				Name:         "test",
				Type:         cty.String,
				ExampleVal:   cty.StringVal("t"),
				Constraints:  constraint.RequiredMeaningfull,
				MaxInclusive: cty.NumberIntVal(-1),
			},
			schemaErrors: []string{`MaxValInclusive specified for "test" must be >= 0`},
			inputVal:     cty.StringVal("hello_123"),
		},
		{
			name: "oneof",
			obj: &dataspec.AttrSpec{
				Name:        "test",
				Type:        cty.String,
				ExampleVal:  cty.StringVal("a"),
				Constraints: constraint.RequiredMeaningfull,
				OneOf: constraint.OneOf{
					cty.StringVal("a"),
					cty.StringVal("b"),
					cty.StringVal("c"),
				},
			},
			inputVal:    cty.StringVal("a"),
			expectedVal: cty.StringVal("a"),
		},
		{
			name: "oneof_err",
			obj: &dataspec.AttrSpec{
				Name:        "test",
				Type:        cty.String,
				ExampleVal:  cty.StringVal("a"),
				Constraints: constraint.RequiredMeaningfull,
				OneOf: constraint.OneOf{
					cty.StringVal("a"),
					cty.StringVal("b"),
					cty.StringVal("c"),
				},
			},
			inputVal: cty.StringVal("x"),
			asserts: [][]testtools.Assert{{
				testtools.IsError,
				testtools.SummaryContains("not one of"),
				testtools.DetailContains(`"a", "b", "c"`),
			}},
		},
		{
			name: "number_min",
			obj: &dataspec.AttrSpec{
				Name:         "test",
				Type:         cty.Number,
				ExampleVal:   cty.NumberIntVal(1),
				Constraints:  constraint.RequiredMeaningfull,
				MinInclusive: cty.NumberFloatVal(0.5),
			},
			inputVal:    cty.NumberIntVal(1),
			expectedVal: cty.NumberIntVal(1),
		},
		{
			name: "number_min_err",
			obj: &dataspec.AttrSpec{
				Name:         "test",
				Type:         cty.Number,
				ExampleVal:   cty.NumberIntVal(1),
				Constraints:  constraint.RequiredMeaningfull,
				MinInclusive: cty.NumberFloatVal(0.5),
			},
			inputVal: cty.NumberIntVal(0),
			asserts: [][]testtools.Assert{{
				testtools.IsError,
				testtools.SummaryContains("Attribute is not in range"),
				testtools.DetailContains(`>=`),
			}},
		},
		{
			name: "number_max",
			obj: &dataspec.AttrSpec{
				Name:         "test",
				Type:         cty.Number,
				ExampleVal:   cty.NumberIntVal(1),
				Constraints:  constraint.RequiredMeaningfull,
				MaxInclusive: cty.NumberFloatVal(2.5),
			},
			inputVal:    cty.NumberFloatVal(2.3),
			expectedVal: cty.NumberFloatVal(2.3),
		},
		{
			name: "number_max_err",
			obj: &dataspec.AttrSpec{
				Name:         "test",
				Type:         cty.Number,
				ExampleVal:   cty.NumberIntVal(1),
				Constraints:  constraint.RequiredMeaningfull,
				MaxInclusive: cty.NumberFloatVal(2.5),
			},
			inputVal: cty.NumberFloatVal(2.7),
			asserts: [][]testtools.Assert{{
				testtools.IsError,
				testtools.SummaryContains("Attribute is not in range"),
			}},
		},
		{
			name: "number_range",
			obj: &dataspec.AttrSpec{
				Name:         "test",
				Type:         cty.Number,
				ExampleVal:   cty.NumberIntVal(2),
				Constraints:  constraint.RequiredMeaningfull,
				MinInclusive: cty.NumberFloatVal(1.5),
				MaxInclusive: cty.NumberFloatVal(2.5),
			},
			inputVal:    cty.NumberFloatVal(1.7),
			expectedVal: cty.NumberFloatVal(1.7),
		},
		{
			name: "number_range_err",
			obj: &dataspec.AttrSpec{
				Name:         "test",
				Type:         cty.Number,
				ExampleVal:   cty.NumberIntVal(2),
				Constraints:  constraint.RequiredMeaningfull,
				MinInclusive: cty.NumberFloatVal(1.5),
				MaxInclusive: cty.NumberFloatVal(2.5),
			},
			inputVal: cty.NumberFloatVal(4.2),
			asserts: [][]testtools.Assert{{
				testtools.IsError,
				testtools.SummaryContains("Attribute is not in range"),
			}},
		},
		{
			name: "deprecation",
			obj: &dataspec.AttrSpec{
				Name:         "test",
				Type:         cty.Number,
				ExampleVal:   cty.NumberIntVal(2),
				Constraints:  constraint.RequiredMeaningfull,
				MinInclusive: cty.NumberFloatVal(1.5),
				MaxInclusive: cty.NumberFloatVal(2.5),
				Depricated:   "test deprication message",
			},
			inputVal: cty.NumberFloatVal(4.2),
			asserts: [][]testtools.Assert{{
				testtools.IsError,
				testtools.SummaryContains("Attribute is not in range"),
			}, {
				testtools.IsWarning,
				testtools.SummaryContains("Deprecated"),
				testtools.DetailContains("test deprication message"),
			}},
			// TODO: filter deprications from validation
		},
	} {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			spec := dataspec.ObjectSpec{tc.obj}
			if !assert.ElementsMatch(t, tc.schemaErrors, spec.ValidateSpec()) {
				t.FailNow()
			}
			if len(tc.schemaErrors) != 0 {
				return
			}

			body := &hclsyntax.Body{
				Attributes: hclsyntax.Attributes{},
			}
			if tc.inputVal != cty.NilVal {
				body.Attributes[tc.obj.Name] = &hclsyntax.Attribute{
					Name: tc.obj.Name,
					Expr: &hclsyntax.LiteralValueExpr{
						Val: tc.inputVal,
					},
				}
			}
			objVal, diags := dataspec.Decode(body, spec, nil)
			testtools.CompareDiags(t, nil, diags, tc.asserts)
			if diags.HasErrors() {
				return
			}
			val := objVal.GetAttr(tc.obj.Name)
			if !val.RawEquals(tc.expectedVal) {
				t.Fatalf("Values not equal. Expected: %s; Got: %s", tc.expectedVal.GoString(), val.GoString())
			}
		})
	}
}
