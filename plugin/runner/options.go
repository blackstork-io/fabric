package runner

import "github.com/blackstork-io/fabric/plugin"

type options struct {
	pluginDir  string
	versionMap VersionMap
	builtin    []*plugin.Schema
}

var defaultOptions = options{
	pluginDir: "./plugins",
}

type Option func(*options)

func WithPluginDir(dir string) Option {
	return func(o *options) {
		o.pluginDir = dir
	}
}

func WithPluginVersions(m VersionMap) Option {
	return func(o *options) {
		o.versionMap = m
	}
}

func WithBuiltIn(builtin ...*plugin.Schema) Option {
	return func(o *options) {
		o.builtin = builtin
	}
}
