package basic

import (
	"context"
	"math/rand"

	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"

	"github.com/blackstork-io/fabric/plugin"
	"github.com/blackstork-io/fabric/plugin/dataspec"
	"github.com/blackstork-io/fabric/plugin/dataspec/constraint"
)

// makeRandomNumbersDataSource creates a new data source for generating random numbers
func makeRandomNumbersDataSource() *plugin.DataSource {
	return &plugin.DataSource{
		// Config is optional, we can define the schema for the config that is reusable for this data source
		Config: dataspec.ObjectSpec{
			&dataspec.AttrSpec{
				Name:       "min",
				Type:       cty.Number,
				Doc:        `Lower bound (inclusive)`,
				DefaultVal: cty.NumberIntVal(0),
			},
			&dataspec.AttrSpec{
				Name:       "max",
				Type:       cty.Number,
				Doc:        `Upper bound (inclusive)`,
				DefaultVal: cty.NumberIntVal(100),
			},
		},
		// We define the schema for the arguments
		Args: dataspec.ObjectSpec{
			&dataspec.AttrSpec{
				Name:        "length",
				Constraints: constraint.Integer | constraint.Required,
				Type:        cty.Number,
				ExampleVal:  cty.NumberIntVal(10),
			},
		},
		// Optional: We can also define the schema for the config
		DataFunc: fetchRandomNumbers,
	}
}

func fetchRandomNumbers(ctx context.Context, params *plugin.RetrieveDataParams) (plugin.Data, hcl.Diagnostics) {
	min := params.Config.GetAttr("min")
	max := params.Config.GetAttr("max")

	// validating the arguments
	length := params.Args.GetAttr("length")
	if length.IsNull() {
		return nil, hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "Failed to parse arguments",
			Detail:   "length is required",
		}}
	}
	if min.GreaterThan(max).True() {
		return nil, hcl.Diagnostics{{
			Severity: hcl.DiagError,
			Summary:  "Failed to parse config",
			Detail:   "min is greater than max",
		}}
	}

	lengthInt, _ := length.AsBigFloat().Int64()
	minInt, _ := min.AsBigFloat().Int64()
	maxInt, _ := max.AsBigFloat().Int64()

	data := make(plugin.ListData, lengthInt)
	for i := int64(0); i < lengthInt; i++ {
		n := rand.Int63() % (maxInt - minInt + 1) //nolint:G404 // weak rng is ok here
		data[i] = plugin.NumberData(n + minInt)
	}
	return data, nil
}
