package diagnostics

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"
)

type PathExtra cty.Path

func (p PathExtra) improveDiagnostic(diag *hcl.Diagnostic) {
	if len(p) == 0 {
		return
	}
	sb := []byte(diag.Detail)
	sb = append(sb, "\nHappened while evaluating value at:\n"...)
	for _, step := range p {
		switch step := step.(type) {
		case cty.GetAttrStep:
			sb = fmt.Appendf(sb, ".%s", step.Name)
			continue
		case cty.IndexStep:
			ty := step.Key.Type()
			if ty.IsPrimitiveType() {
				if ty == cty.String {
					sb = fmt.Appendf(sb, "[%q]", step.Key.AsString())
					continue
				} else if ty == cty.Number {
					i, _ := step.Key.AsBigFloat().Int64()
					sb = fmt.Appendf(sb, "[%d]", i)
					continue
				}
			}
			sb = fmt.Appendf(sb, "[%s]", step.Key.GoString())
			continue
		}
		sb = fmt.Appendf(sb, "[%+v]", step)
	}
	diag.Detail = string(sb)
}

func NewPathExtra(p cty.Path) PathExtra {
	return PathExtra(p.Copy())
}
