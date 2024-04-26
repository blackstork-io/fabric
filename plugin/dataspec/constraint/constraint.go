package constraint

import (
	"bytes"
	"fmt"

	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
)

type Constraints uint32

func (c Constraints) Is(test Constraints) bool {
	return (c & test) == test
}

const (
	// Attribute can't be left out (but can be null, empty, etc.)
	Required Constraints = (1 << iota)
	// Attribute can't be null
	NonNull
	// Attribute can't be an empty collections or an empty string
	NonEmpty
	// If an attribute is a string - preprocess it with strings.TrimSpace
	TrimSpace
	// Attribute must be a whole integer number
	Integer
)

const (
	// Attribute will be trimmed and later checked to be non-empty
	// If attribute is not a string - trim operation does nothing
	TrimmedNonEmpty Constraints = NonEmpty | TrimSpace

	// Attribute is required, non-null, non-empty, strings are trimmed
	RequiredMeaningfull Constraints = Required | NonNull | TrimmedNonEmpty
)

type OneOf []cty.Value

var _ fmt.Stringer = OneOf{}

func (o OneOf) String() string {
	f := hclwrite.NewFile()
	f.Body().SetAttributeValue("t", cty.TupleVal(o))

	b := bytes.TrimSpace(hclwrite.Format(f.Bytes()))
	b = b[len("t = [") : len(b)-len("]")]
	return string(bytes.TrimSpace(b))
}

func (o OneOf) IsEmpty() bool {
	return len(o) == 0
}
