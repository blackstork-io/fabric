package eval

import (
	"github.com/blackstork-io/fabric/plugin"
)

type DataSources interface {
	DataSource(name string) (*plugin.DataSource, bool)
}

type ContentProviders interface {
	ContentProvider(name string) (*plugin.ContentProvider, bool)
}

type Publishers interface {
	Publisher(name string) (*plugin.Publisher, bool)
}

type NodeRenderers interface {
	NodeRenderer(customNodeType string) (plugin.NodeRendererFunc, bool)
	AllNodeRenderers() map[string]struct{}
}

type Plugins interface {
	DataSources
	ContentProviders
	Publishers
	NodeRenderers
}
