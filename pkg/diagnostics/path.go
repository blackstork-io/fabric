package diagnostics

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"
)

type PathExtra cty.Path

func (p PathExtra) String() string {
	if len(p) == 0 {
		return ""
	}
	var sb []byte
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
	return string(sb)
}

func (p PathExtra) improveDiagnostic(diag *hcl.Diagnostic) {
	if len(p) == 0 {
		return
	}
	diag.Detail = fmt.Sprintf("%s\nHappened while evaluating value at:\n%s", diag.Detail, p)
}

func NewPathExtra(p cty.Path) PathExtra {
	return PathExtra(p.Copy())
}
