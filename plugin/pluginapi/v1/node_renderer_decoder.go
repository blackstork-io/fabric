package pluginapiv1

import (
	"github.com/blackstork-io/fabric/pkg/utils"
	"github.com/blackstork-io/fabric/plugin"
	astv1 "github.com/blackstork-io/fabric/plugin/ast/v1"
)

func decodePublisherInfo(src *NodeSupportInfo) plugin.PublisherInfo {
	return plugin.PublisherInfo{
		SupportedCustomNodes: src.GetSupportedCustomNodes(),
		UnsupportedNodes:     src.GetUnsupportedNodes(),
	}
}

func decodeRenderNodeParams(src *RenderNodeRequest) *plugin.RenderNodeParams {
	if src == nil {
		return nil
	}
	return &plugin.RenderNodeParams{
		Subtree:         astv1.DecodeNode(src.GetSubtree()),
		NodePath:        astv1.DecodePath(src.GetNodePath()),
		Publisher:       src.GetPublisher(),
		PublisherInfo:   decodePublisherInfo(src.GetPublisherInfo()),
		CustomRenderers: utils.SliceToSet(src.GetCustomNodeRenderers()),
	}
}
