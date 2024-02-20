package basic

import (
	"context"
	"math/rand"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/zclconf/go-cty/cty"

	"github.com/blackstork-io/fabric/plugin"
)

const (
	// defaultMin is the default minimum value for random number generation
	defaultMin = 0
	// defaultMax is the default maximum value for random number generation
	defaultMax = 100
)

// makeRandomNumbersDataSource creates a new data source for generating random numbers
func makeRandomNumbersDataSource() *plugin.DataSource {
	return &plugin.DataSource{
		// Config is optional, we can define the schema for the config that is reusable for this data source
		Config: hcldec.ObjectSpec{
			"min": &hcldec.AttrSpec{
				Name:     "min",
				Required: false,
				Type:     cty.Number,
			},
			"max": &hcldec.AttrSpec{
				Name:     "max",
				Required: false,
				Type:     cty.Number,
			},
		},
		// We define the schema for the arguments
		Args: hcldec.ObjectSpec{
			"length": &hcldec.AttrSpec{
				Name:     "length",
				Required: true,
				Type:     cty.Number,
			},
		},
		// Optional: We can also define the schema for the config
		DataFunc: fetchRandomNumbers,
	}
}

func fetchRandomNumbers(ctx context.Context, params *plugin.RetrieveDataParams) (plugin.Data, hcl.Diagnostics) {
	min := params.Config.GetAttr("min")
	max := params.Config.GetAttr("max")
	if min.IsNull() {
		min = cty.NumberIntVal(defaultMin)
	}
	if max.IsNull() {
		max = cty.NumberIntVal(defaultMax)
	}
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
