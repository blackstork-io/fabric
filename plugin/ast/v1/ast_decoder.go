package astv1

import (
	"log/slog"

	"github.com/blackstork-io/fabric/pkg/utils"
	"github.com/blackstork-io/fabric/plugin/ast/nodes"
)

func DecodeNode(node *Node) *nodes.Node {
	if node == nil {
		return nil
	}
	var res *nodes.Node
	switch node.WhichContent() {
	case Node_Paragraph_case:
		res = nodes.NewNode(
			&nodes.Paragraph{
				IsTextBlock: node.GetParagraph().GetIsTextBlock(),
			},
		)
	case Node_Heading_case:
		res = nodes.NewNode(
			&nodes.Heading{
				Level: int(
					utils.Clamp(1, node.GetHeading().GetLevel(), 6),
				),
			},
		)
	case Node_ThematicBreak_case:
		res = nodes.NewNode(
			&nodes.ThematicBreak{},
		)
	case Node_CodeBlock_case:
		n := node.GetCodeBlock()
		res = nodes.NewNode(
			&nodes.CodeBlock{
				Language: n.GetLanguage(),
				Code:     n.GetCode(),
			},
		)
	case Node_Blockquote_case:
		res = nodes.NewNode(
			&nodes.Blockquote{},
		)
	case Node_List_case:
		n := node.GetList()
		res = nodes.NewNode(
			&nodes.List{
				Marker: byte(n.GetMarker()),
				Start:  n.GetStart(),
			},
		)
	case Node_ListItem_case:
		res = nodes.NewNode(
			&nodes.ListItem{},
		)
	case Node_HtmlBlock_case:
		res = nodes.NewNode(
			&nodes.HTMLBlock{
				HTML: node.GetHtmlBlock().GetHtml(),
			},
		)
	case Node_Text_case:
		n := node.GetText()
		res = nodes.NewNode(
			&nodes.Text{
				Text:          n.GetText(),
				HardLineBreak: n.GetHardLineBreak(),
			},
		)
	case Node_CodeSpan_case:
		res = nodes.NewNode(
			&nodes.CodeSpan{
				Code: node.GetCodeSpan().GetCode(),
			},
		)
	case Node_Emphasis_case:
		res = nodes.NewNode(
			&nodes.Emphasis{
				Level: int(node.GetEmphasis().GetLevel()),
			},
		)
	case Node_Link_case:
		n := node.GetLink()
		res = nodes.NewNode(
			&nodes.Link{
				Destination: n.GetDestination(),
				Title:       n.GetTitle(),
			},
		)
	case Node_Image_case:
		n := node.GetImage()
		res = nodes.NewNode(
			&nodes.Image{
				Source: n.GetSource(),
				Alt:    n.GetAlt(),
			},
		)
	case Node_AutoLink_case:
		res = nodes.NewNode(
			&nodes.AutoLink{
				Value: node.GetAutoLink().GetValue(),
			},
		)
	case Node_HtmlInline_case:
		res = nodes.NewNode(
			&nodes.HTMLInline{
				HTML: node.GetHtmlInline().GetHtml(),
			},
		)
	case Node_Table_case:
		res = nodes.NewNode(
			&nodes.Table{
				Alignments: utils.FnMap(node.GetTable().GetAlignments(), func(a CellAlignment) nodes.Alignment {
					switch a {
					case CellAlignment_CELL_ALIGNMENT_UNSPECIFIED:
						return nodes.AlignmentNone
					case CellAlignment_CELL_ALIGNMENT_LEFT:
						return nodes.AlignmentLeft
					case CellAlignment_CELL_ALIGNMENT_CENTER:
						return nodes.AlignmentCenter
					case CellAlignment_CELL_ALIGNMENT_RIGHT:
						return nodes.AlignmentRight
					default:
						slog.Error("unsupported cell alignment", "alignment", a)
						return nodes.AlignmentNone
					}
				}),
			},
		)
	case Node_TableRow_case:
		res = nodes.NewNode(
			&nodes.TableRow{},
		)
	case Node_TableCell_case:
		res = nodes.NewNode(
			&nodes.TableCell{},
		)
	case Node_TaskCheckbox_case:
		res = nodes.NewNode(
			&nodes.TaskCheckbox{
				Checked: node.GetTaskCheckbox().GetChecked(),
			},
		)
	case Node_Strikethrough_case:
		res = nodes.NewNode(
			&nodes.Strikethrough{},
		)
	case Node_FabricDocument_case:
		res = nodes.NewNode(
			&nodes.FabricDocument{},
		)
	case Node_FabricSection_case:
		res = nodes.NewNode(
			&nodes.FabricSection{},
		)
	case Node_FabricContent_case:
		fc := node.GetFabricContent()
		res = nodes.NewNode(
			&nodes.FabricContent{
				Meta: DecodeMetadata(fc.GetMeta()),
			},
		)
	case Node_Custom_case:
		c := node.GetCustom()
		c.GetScope()
		res = nodes.NewNode(
			&nodes.Custom{
				Data:  c.GetData(),
				Scope: decodeScope(c.GetScope()),
			},
		)
	default:
		slog.Warn("unsupported AST node type", "type", node.WhichContent().String())
	}
	res.AppendChildren(DecodeChildren(node)...)
	return res
}

func decodeScope(s Custom_RenderScope) nodes.RendererScope {
	switch s {
	case Custom_SCOPE_UNSPECIFIED, Custom_SCOPE_NODE:
		return nodes.ScopeNode
	case Custom_SCOPE_CONTENT:
		return nodes.ScopeContent
	case Custom_SCOPE_SECTION:
		return nodes.ScopeSection
	case Custom_SCOPE_DOCUMENT:
		return nodes.ScopeDocument
	default:
		slog.Error("unsupported scope", "scope", s)
		return nodes.ScopeNode
	}
}

func DecodeChildren(node *Node) []*nodes.Node {
	return utils.FnMap(node.GetChildren(), DecodeNode)
}

func DecodeMetadata(meta *FabricContent_Metadata) *nodes.FabricContentMetadata {
	if meta == nil {
		return nil
	}
	return &nodes.FabricContentMetadata{
		Provider: meta.GetProvider(),
		Plugin:   meta.GetPlugin(),
		Version:  meta.GetVersion(),
	}
}

func DecodePath(path *Path) nodes.Path {
	p := path.GetPath()
	if p == nil {
		return nil
	}
	return nodes.Path(utils.FnMap(p, func(i uint32) int { return int(i) }))
}
