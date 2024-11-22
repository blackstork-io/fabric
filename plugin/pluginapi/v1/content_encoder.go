package pluginapiv1

import (
	"fmt"
	"log/slog"

	"github.com/blackstork-io/fabric/pkg/utils"
	"github.com/blackstork-io/fabric/plugin"
	astv1 "github.com/blackstork-io/fabric/plugin/ast/v1"
)

func encodeContentResult(src *plugin.ContentResult) *ContentResult {
	if src == nil {
		return nil
	}
	return &ContentResult{
		Content:  EncodeContent(src.Content),
		Location: encodeLocation(src.Location),
	}
}

func EncodeContent(src plugin.Content) *Content {
	var variant isContent_Value
	switch val := src.(type) {
	case nil:
		return nil
	case *plugin.ContentElement:
		if val == nil {
			break
		}
		el := &ContentElement{
			Markdown: val.AsMarkdownSrc(),
			Meta:     astv1.EncodeMetadata(val.Meta()),
		}
		if val.IsAst() {
			el.Ast = val.AsSerializedNode()
		}
		variant = &Content_Element{
			Element: el,
		}
	case *plugin.ContentSection:
		if val == nil {
			break
		}
		variant = &Content_Section{
			Section: &ContentSection{
				Children: utils.FnMap(val.Children, EncodeContent),
				Meta:     astv1.EncodeMetadata(val.Meta()),
			},
		}
	case *plugin.ContentEmpty:
		variant = &Content_Empty{
			Empty: &ContentEmpty{},
		}
	default:
		slog.Error("unknown content type", "type", fmt.Sprintf("%T", src))
	}

	return &Content{
		Value: variant,
	}
}

func encodeLocation(src *plugin.Location) *Location {
	if src == nil {
		return nil
	}
	return &Location{
		Index:  src.Index,
		Effect: encodeLocationEffect(src.Effect),
	}
}

func encodeLocationEffect(src plugin.LocationEffect) LocationEffect {
	switch src {
	case plugin.LocationEffectBefore:
		return LocationEffect_LOCATION_EFFECT_BEFORE
	case plugin.LocationEffectAfter:
		return LocationEffect_LOCATION_EFFECT_AFTER
	default:
		return LocationEffect_LOCATION_EFFECT_UNSPECIFIED
	}
}
