package dataspec

import (
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

// Root-only spec
type ObjDumpSpec struct {
	rootSpecSigil
	Doc string
}

func (o *ObjDumpSpec) getSpec() Spec {
	return o
}

func (o *ObjDumpSpec) HcldecSpec() hcldec.Spec {
	return nil
}

func (o *ObjDumpSpec) WriteDoc(w *hclwrite.Body) {
	tokens := comment(nil, o.Doc)
	if len(tokens) != 0 {
		tokens = append(tokens, &hclwrite.Token{
			Type:  hclsyntax.TokenNewline,
			Bytes: []byte{'\n'},
		})
	}
	w.AppendUnstructuredTokens(tokens)
}

func (*ObjDumpSpec) IsEmpty() bool {
	return false
}

func (*ObjDumpSpec) ValidateSpec() (errs []string) {
	return
}
