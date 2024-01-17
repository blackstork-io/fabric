package parser

import (
	"strings"

	"github.com/sanity-io/litter"
	"golang.org/x/exp/maps"

	"github.com/blackstork-io/fabric/parser/definitions"
	"github.com/blackstork-io/fabric/parser/evaluation"
	"github.com/blackstork-io/fabric/pkg/diagnostics"
)

type MockCaller struct{}

func (c *MockCaller) dumpContext(context map[string]any) string {
	return litter.Sdump("Context:", context)
}

func (c *MockCaller) dumpConfig(config evaluation.Configuration) string {
	switch c := config.(type) {
	case *definitions.ConfigPtr:
		if c == nil {
			return "NoConfig"
		}
		attrs, _ := c.Cfg.Body.JustAttributes()
		return litter.Sdump("ConfigPtr", maps.Keys(attrs))
	case *definitions.Config:
		if c == nil {
			return "NoConfig"
		}
		attrs, _ := c.Block.Body.JustAttributes()
		return litter.Sdump("Config", maps.Keys(attrs))
	case nil:
		return "NoConfig"
	default:
		return "UnknownConfig " + litter.Sdump(c)
	}
}

func (c *MockCaller) dumpInvocation(invoke evaluation.Invocation) string {
	switch inv := invoke.(type) {
	case *evaluation.BlockInvocation:
		if inv == nil {
			return "NoConfig"
		}
		attrStringed := map[string]string{}
		attrs, _ := inv.Body.JustAttributes()
		for k, v := range attrs {
			val, _ := v.Expr.Value(nil)
			attrStringed[k] = val.GoString()
		}

		return litter.Sdump("BlockInvocation", attrStringed)
	case *definitions.TitleInvocation:
		if inv == nil {
			return "NoInvocation"
		}
		val, _ := inv.Expression.Value(nil)
		return litter.Sdump("TitleInvocation", val.GoString())
	case nil:
		return "NoInvocation"
	default:
		return "UnknownInvocation " + litter.Sdump(inv)
	}
}

// CallContent implements PluginCaller.
func (c *MockCaller) CallContent(name string, config evaluation.Configuration, invocation evaluation.Invocation, context map[string]any) (result string, diag diagnostics.Diag) {
	dump := []string{
		"Call to content:",
	}
	dump = append(dump, c.dumpConfig(config))
	dump = append(dump, c.dumpInvocation(invocation))
	dump = append(dump, c.dumpContext(context))
	return strings.Join(dump, "\n") + "\n\n", nil
}

// CallData implements PluginCaller.
func (c *MockCaller) CallData(name string, config evaluation.Configuration, invocation evaluation.Invocation) (result map[string]any, diag diagnostics.Diag) {
	dump := []string{
		"Call to data:",
	}
	dump = append(dump, c.dumpConfig(config))
	dump = append(dump, c.dumpInvocation(invocation))
	return map[string]any{"dumpResult": strings.Join(dump, "\n")}, nil
}

var _ PluginCaller = (*MockCaller)(nil)