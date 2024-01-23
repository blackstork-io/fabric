package plugin

import (
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
	// The namespace used during installation of the plugins
	Namespace string
	// "content" or "data" for now
	Kind string
	// "text", "plugin_a", etc.
	Name string
	// version of the plugin `Kind Name` that is provided by the current binary
	Version Version
	// Specification of the `config` block for this plugin
	// If nil - providing a `config Kind Name` is an error
	ConfigSpec hcldec.Spec
	// Specification of the invocation block's body, i.e. `content text {<spec of what's here>}`
	InvocationSpec hcldec.Spec
}

type Args struct {
	// Specifies which kind, name and version of plugin to execute
	Kind    string
	Name    string
	Version Version

	// Result of decoding a config block with ConfigSpec
	Config cty.Value
	// Result of decoding an invocation block with InvocationSpec
	Args cty.Value
	// Passed to content plugins, nil for data plugins
	Context map[string]any
}

type Result struct {
	// `content` plugins return a markdown string
	// `data` plugins return a map[string]any that would be put into the global config
	Result any
	Diags  hcl.Diagnostics
}
