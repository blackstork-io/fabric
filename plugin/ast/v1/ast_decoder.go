package astv1

import (
	"fmt"
	"log/slog"

	"github.com/blackstork-io/fabric/pkg/utils"
	"github.com/blackstork-io/fabric/plugin/ast/nodes"
)

func DecodeNode(node *Node) *nodes.Node {
	if node == nil {
		return nil
	}
	var res *nodes.Node
	switch n := node.GetContent().(type) {
	case *Node_Document:
		res = nodes.NewNode(
			&nodes.Document{},
		)
	case *Node_Paragraph:
		res = nodes.NewNode(
			&nodes.Paragraph{
				IsTextBlock: n.Paragraph.GetIsTextBlock(),
			},
		)
	case *Node_Heading:
		res = nodes.NewNode(
			&nodes.Heading{
				Level: int(
					utils.Clamp(1, n.Heading.GetLevel(), 6),
				),
			},
		)
	case *Node_ThematicBreak:
		res = nodes.NewNode(
			&nodes.ThematicBreak{},
		)
	case *Node_CodeBlock:
		res = nodes.NewNode(
			&nodes.CodeBlock{
				Language: n.CodeBlock.GetLanguage(),
				Code:     n.CodeBlock.GetCode(),
			},
		)
	case *Node_Blockquote:
		res = nodes.NewNode(
			&nodes.Blockquote{},
		)
	case *Node_List:
		res = nodes.NewNode(
			&nodes.List{
				Marker: byte(n.List.GetMarker()),
				Start:  n.List.GetStart(),
			},
		)
	case *Node_ListItem:
		res = nodes.NewNode(
			&nodes.ListItem{},
		)
	case *Node_HtmlBlock:
		res = nodes.NewNode(
			&nodes.HTMLBlock{
				HTML: n.HtmlBlock.GetHtml(),
			},
		)
	case *Node_Text:
		res = nodes.NewNode(
			&nodes.Text{
				Text:          n.Text.GetText(),
				HardLineBreak: n.Text.GetHardLineBreak(),
			},
		)
	case *Node_CodeSpan:
		res = nodes.NewNode(
			&nodes.CodeSpan{
				Code: n.CodeSpan.GetCode(),
			},
		)
	case *Node_Emphasis:
		res = nodes.NewNode(
			&nodes.Emphasis{
				Level: int(n.Emphasis.GetLevel()),
			},
		)
	case *Node_Link:
		res = nodes.NewNode(
			&nodes.Link{
				Destination: n.Link.GetDestination(),
				Title:       n.Link.GetTitle(),
			},
		)
	case *Node_Image:
		res = nodes.NewNode(
			&nodes.Image{
				Source: n.Image.GetSource(),
				Alt:    n.Image.GetAlt(),
			},
		)
	case *Node_AutoLink:
		res = nodes.NewNode(
			&nodes.AutoLink{
				Value: n.AutoLink.GetValue(),
			},
		)
	case *Node_HtmlInline:
		res = nodes.NewNode(
			&nodes.HTMLInline{
				HTML: n.HtmlInline.GetHtml(),
			},
		)
	case *Node_Table:
		res = nodes.NewNode(
			&nodes.Table{
				Alignments: utils.FnMap(n.Table.GetAlignments(), func(a CellAlignment) nodes.Alignment {
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
	case *Node_TableRow:
		res = nodes.NewNode(
			&nodes.TableRow{},
		)
	case *Node_TableCell:
		res = nodes.NewNode(
			&nodes.TableCell{},
		)
	case *Node_TaskCheckbox:
		res = nodes.NewNode(
			&nodes.TaskCheckbox{
				Checked: n.TaskCheckbox.GetChecked(),
			},
		)
	case *Node_Strikethrough:
		res = nodes.NewNode(
			&nodes.Strikethrough{},
		)
	case *Node_Custom:
		res = nodes.NewNode(
			&nodes.Custom{
				Data: n.Custom.GetData(),
			},
		)
	default:
		panic(fmt.Errorf("unsupported node type: %T", n))
	}
	res.AppendChildren(DecodeChildren(node)...)
	return res
}

func DecodeChildren(node *Node) []*nodes.Node {
	return utils.FnMap(node.GetChildren(), DecodeNode)
}

func DecodeMetadata(meta *Metadata) *nodes.ContentMeta {
	if meta == nil {
		return nil
	}
	return &nodes.ContentMeta{
		Provider: meta.GetProvider(),
		Plugin:   meta.GetPlugin(),
		Version:  meta.GetVersion(),
	}
}
