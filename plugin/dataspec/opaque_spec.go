package dataspec

import (
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"

	"github.com/blackstork-io/fabric/pkg/diagnostics"
)

// OpaqueSpec adds an ability to use any hcldec.Spec
// (without automatic documentation generation).
type OpaqueSpec struct {
	Spec hcldec.Spec
	Doc  string
}

func (o *OpaqueSpec) getSpec() Spec {
	return o
}

func (o *OpaqueSpec) HcldecSpec() hcldec.Spec {
	return o.Spec
}

func (o *OpaqueSpec) WriteDoc(w *hclwrite.Body) {
	tokens := comment(nil, o.Doc)
	if len(tokens) != 0 {
		tokens = append(tokens, &hclwrite.Token{
			Type:  hclsyntax.TokenNewline,
			Bytes: []byte{'\n'},
		})
	}
	w.AppendUnstructuredTokens(tokens)
}

func (*OpaqueSpec) ValidateSpec() (errs diagnostics.Diag) {
	return
}
