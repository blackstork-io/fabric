package plugin

import (
	"github.com/Masterminds/semver/v3"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/zclconf/go-cty/cty"
)

// The interface exposed and used by the go-plugin (versioned as 1)
type PluginRPC interface {
	GetPlugins() []Plugin
	Call(args Args) Result
}

// One go-plugin binary can provide multiple plugins and/or one plugin with multiple versions
// Deciding if we actually will do bundling in that way or stick with one plugin - one binary
// is a different question, but the interface wouldn't limit us either way.
type Plugin struct {
	// "content" or "data" for now
	Kind string
	// "text", "plugin_a", etc.
	Name string
	// version of the plugin `Kind Name` that is provided by the current binary
	// github.com/Masterminds/semver provides semver and constraints parsing and validating
	Version semver.Version
	// Specification of the `config` block for this plugin
	// If nil - providing a `config Kind Name` is an error
	ConfigSpec hcldec.Spec
	// Specification of the invocation block's body, i.e. `content text {<spec of what's here>}`
	InvocationSpec hcldec.Spec
}

// We can define arguments in the interface, but:
// 1) too many arguments for a single function call
// 2) net/rpc would still require us to put all args into a single struct
type Args struct {
	Kind string
	Name string
	// useful when multiple plugin versions live in the same binary
	// if not - it's just an extra "1.2.3" string that is unused
	Version semver.Version

	// Result of decoding a config block with ConfigSpec
	Config cty.Value
	// Result of decoding an invocation block with InvocationSpec
	Args cty.Value
	// Passed to content plugins, nil for data plugins, I guess?
	Context map[string]any
}

type Result struct {
	// I'm struggling with correctly typing the result:
	// content blocks return a string
	// data blocks return a map[string]any (gojq breaks if we specify other types),
	// that would be put into the global config
	// Should the plugin be responsible for determining where the result goes,
	// or will we hard-code it in the main app depending on the plugin kind?
	result any
	diags  hcl.Diagnostics
}
