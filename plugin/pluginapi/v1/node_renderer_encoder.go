package pluginapiv1

import "github.com/blackstork-io/fabric/plugin"

func encodePublisherInfo(src plugin.PublisherInfo) *NodeSupportInfo {
	return &NodeSupportInfo{
		SupportedCustomNodes: src.SupportedCustomNodes,
		UnsupportedNodes:     src.UnsupportedNodes,
	}
}
