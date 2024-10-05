package pluginapiv1

import (
	"fmt"
	"log/slog"

	"github.com/blackstork-io/fabric/pkg/utils"
	"github.com/blackstork-io/fabric/plugin"
)

func decodeContentResult(src *ContentResult) *plugin.ContentResult {
	if src == nil {
		return nil
	}
	return &plugin.ContentResult{
		Content:  decodeContent(src.Content),
		Location: decodeLocation(src.Location),
	}
}

func decodeContent(src *Content) plugin.Content {
	switch val := src.GetValue().(type) {
	case *Content_Element:
		return plugin.NewElementFromMarkdownAndAST(val.Element.GetMarkdown(), val.Element.GetAst())
	case *Content_Section:
		return &plugin.ContentSection{
			Children: utils.FnMap(val.Section.GetChildren(), decodeContent),
		}
	case *Content_Empty:
		return &plugin.ContentEmpty{}
	case nil:
		return nil
	default:
		slog.Error("unknown content type", "type", fmt.Sprintf("%T", src))
		return nil
	}
}

func decodeLocation(src *Location) *plugin.Location {
	if src == nil {
		return nil
	}
	return &plugin.Location{
		Index:  src.GetIndex(),
		Effect: decodeLocationEffect(src.GetEffect()),
	}
}

func decodeLocationEffect(src LocationEffect) plugin.LocationEffect {
	switch src {
	case LocationEffect_LOCATION_EFFECT_BEFORE:
		return plugin.LocationEffectBefore
	case LocationEffect_LOCATION_EFFECT_AFTER:
		return plugin.LocationEffectAfter
	default:
		return plugin.LocationEffectUnspecified
	}
}
