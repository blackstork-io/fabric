package astv1

import (
	"fmt"
	"log/slog"

	"github.com/blackstork-io/fabric/pkg/utils"
	"github.com/blackstork-io/fabric/plugin/ast/nodes"
)

func EncodeNode(node *nodes.Node) *Node {
	if node == nil {
		return nil
	}
	var res *Node
	switch n := node.Content.(type) {
	case *nodes.Document:
		res = &Node{
			Content: &Node_Document{},
		}
	case *nodes.Paragraph:
		res = &Node{
			Content: &Node_Paragraph{
				Paragraph: &Paragraph{
					IsTextBlock: n.IsTextBlock,
				},
			},
		}
	case *nodes.Heading:
		res = &Node{
			Content: &Node_Heading{
				Heading: &Heading{
					Level: uint32(n.Level),
				},
			},
		}
	case *nodes.ThematicBreak:
		res = &Node{
			Content: &Node_ThematicBreak{},
		}
	case *nodes.CodeBlock:
		res = &Node{
			Content: &Node_CodeBlock{
				CodeBlock: &CodeBlock{
					Language: n.Language,
					Code:     n.Code,
				},
			},
		}
	case *nodes.Blockquote:
		res = &Node{
			Content: &Node_Blockquote{},
		}
	case *nodes.List:
		res = &Node{
			Content: &Node_List{
				List: &List{
					Marker: uint32(n.Marker),
					Start:  n.Start,
				},
			},
		}
	case *nodes.ListItem:
		res = &Node{
			Content: &Node_ListItem{},
		}
	case *nodes.HTMLBlock:
		res = &Node{
			Content: &Node_HtmlBlock{
				HtmlBlock: &HTMLBlock{
					Html: n.HTML,
				},
			},
		}
	case *nodes.Text:
		res = &Node{
			Content: &Node_Text{
				Text: &Text{
					Text:          n.Text,
					HardLineBreak: n.HardLineBreak,
				},
			},
		}
	case *nodes.CodeSpan:
		res = &Node{
			Content: &Node_CodeSpan{
				CodeSpan: &CodeSpan{
					Code: n.Code,
				},
			},
		}
	case *nodes.Emphasis:
		res = &Node{
			Content: &Node_Emphasis{
				Emphasis: &Emphasis{
					Level: int64(n.Level),
				},
			},
		}
	case *nodes.Link:
		res = &Node{
			Content: &Node_Link{
				Link: &Link{
					Destination: n.Destination,
					Title:       n.Title,
				},
			},
		}
	case *nodes.Image:
		res = &Node{
			Content: &Node_Image{
				Image: &Image{
					Source: n.Source,
					Alt:    n.Alt,
				},
			},
		}
	case *nodes.AutoLink:
		res = &Node{
			Content: &Node_AutoLink{
				AutoLink: &AutoLink{
					Value: n.Value,
				},
			},
		}
	case *nodes.HTMLInline:
		res = &Node{
			Content: &Node_HtmlInline{
				HtmlInline: &HTMLInline{
					Html: n.HTML,
				},
			},
		}
	case *nodes.Table:
		res = &Node{
			Content: &Node_Table{
				Table: &Table{
					Alignments: utils.FnMap(n.Alignments, func(a nodes.Alignment) CellAlignment {
						switch a {
						case nodes.AlignmentNone:
							return CellAlignment_CELL_ALIGNMENT_UNSPECIFIED
						case nodes.AlignmentLeft:
							return CellAlignment_CELL_ALIGNMENT_LEFT
						case nodes.AlignmentCenter:
							return CellAlignment_CELL_ALIGNMENT_CENTER
						case nodes.AlignmentRight:
							return CellAlignment_CELL_ALIGNMENT_RIGHT
						default:
							slog.Error("unsupported cell alignment", "alignment", a)
							return CellAlignment_CELL_ALIGNMENT_UNSPECIFIED
						}
					}),
				},
			},
		}
	case *nodes.TableRow:
		res = &Node{
			Content: &Node_TableRow{},
		}
	case *nodes.TableCell:
		res = &Node{
			Content: &Node_TableCell{},
		}
	case *nodes.TaskCheckbox:
		res = &Node{
			Content: &Node_TaskCheckbox{
				TaskCheckbox: &TaskCheckbox{
					Checked: n.Checked,
				},
			},
		}
	case *nodes.Strikethrough:
		res = &Node{
			Content: &Node_Strikethrough{},
		}
	case *nodes.Custom:
		res = &Node{
			Content: &Node_Custom{
				Custom: &Custom{
					Data: n.Data,
				},
			},
		}
	default:
		panic(fmt.Errorf("unsupported node type: %T", n))
	}
	res.Children = EncodeChildren(node)
	return res
}

func EncodeChildren(node *nodes.Node) []*Node {
	return utils.FnMap(node.GetChildren(), EncodeNode)
}

func EncodeMetadata(meta *nodes.ContentMeta) *Metadata {
	if meta == nil {
		return nil
	}
	return &Metadata{
		Provider: meta.Provider,
		Plugin:   meta.Plugin,
		Version:  meta.Version,
	}
}
