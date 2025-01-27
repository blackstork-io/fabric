package plugin

import (
	"context"

	"github.com/blackstork-io/fabric/pkg/diagnostics"
	"github.com/blackstork-io/fabric/plugin/ast/nodes"
)

// NodeRenderers is a map of node type to node renderer function
// Key is the node type, without the prefix and plugin name
type NodeRenderers map[string]NodeRendererFunc

type NodeRendererFunc func(ctx context.Context, params *RenderNodeParams) ([]*nodes.Node, diagnostics.Diag)

type RenderNodeParams struct {
	Subtree         *nodes.Node
	NodePath        nodes.Path
	Publisher       string
	PublisherInfo   PublisherInfo
	CustomRenderers map[string]struct{}
}
