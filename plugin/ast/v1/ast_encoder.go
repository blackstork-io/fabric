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
	res := new(Node)

	switch n := node.Content.(type) {
	case *nodes.Paragraph:
		p := new(Paragraph)
		p.SetIsTextBlock(n.IsTextBlock)
		res.SetParagraph(p)
	case *nodes.Heading:
		h := new(Heading)
		h.SetLevel(uint32(n.Level))
		res.SetHeading(h)
	case *nodes.ThematicBreak:
		res.SetThematicBreak(new(ThematicBreak))
	case *nodes.CodeBlock:
		res.SetCodeBlock(
			CodeBlock_builder{
				Code: n.Code,
				// builder preserves nil value here
				Language: n.Language,
			}.Build(),
		)
	case *nodes.Blockquote:
		res.SetBlockquote(new(Blockquote))
	case *nodes.List:
		l := new(List)
		l.SetMarker(uint32(n.Marker))
		l.SetStart(n.Start)
		res.SetList(l)
	case *nodes.ListItem:
		res.SetListItem(new(ListItem))
	case *nodes.HTMLBlock:
		hb := new(HTMLBlock)
		hb.SetHtml(n.HTML)
		res.SetHtmlBlock(hb)
	case *nodes.Text:
		t := new(Text)
		t.SetHardLineBreak(n.HardLineBreak)
		t.SetText(n.Text)
		res.SetText(t)
	case *nodes.CodeSpan:
		cs := new(CodeSpan)
		cs.SetCode(n.Code)
		res.SetCodeSpan(cs)
	case *nodes.Emphasis:
		em := new(Emphasis)
		em.SetLevel(int64(n.Level))
		res.SetEmphasis(em)
	case *nodes.Link:
		l := new(Link)
		l.SetDestination(n.Destination)
		l.SetTitle(n.Title)
		res.SetLink(l)
	case *nodes.Image:
		i := new(Image)
		i.SetSource(n.Source)
		i.SetAlt(n.Alt)
		res.SetImage(i)
	case *nodes.AutoLink:
		al := new(AutoLink)
		al.SetValue(n.Value)
		res.SetAutoLink(al)
	case *nodes.HTMLInline:
		hi := new(HTMLInline)
		hi.SetHtml(n.HTML)
		res.SetHtmlInline(hi)
	case *nodes.Table:
		t := new(Table)
		t.SetAlignments(utils.FnMap(n.Alignments, func(a nodes.Alignment) CellAlignment {
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
		}))
		res.SetTable(t)
	case *nodes.TableRow:
		res.SetTableRow(new(TableRow))
	case *nodes.TableCell:
		res.SetTableCell(new(TableCell))
	case *nodes.TaskCheckbox:
		cb := new(TaskCheckbox)
		cb.SetChecked(n.Checked)
		res.SetTaskCheckbox(cb)
	case *nodes.Strikethrough:
		res.SetStrikethrough(new(Strikethrough))

	case *nodes.FabricDocument:
		res.SetFabricDocument(new(FabricDocument))
	case *nodes.FabricSection:
		res.SetFabricSection(new(FabricSection))
	case *nodes.FabricContent:
		res.SetFabricContent(
			FabricContent_builder{
				Meta: EncodeMetadata(n.Meta),
			}.Build(),
		)
	case *nodes.Custom:
		c := new(Custom)
		c.SetData(n.Data)
		c.SetScope(encodeScope(n.Scope))
		res.SetCustom(c)
	default:
		slog.Warn("unsupported node type", "type", fmt.Sprintf("%T", n))
	}
	res.SetChildren(EncodeChildren(node))
	return res
}

func encodeScope(scope nodes.RendererScope) Custom_RenderScope {
	switch scope {
	case nodes.ScopeNode:
		return Custom_RENDER_SCOPE_NODE
	case nodes.ScopeContent:
		return Custom_RENDER_SCOPE_CONTENT
	case nodes.ScopeSection:
		return Custom_RENDER_SCOPE_SECTION
	case nodes.ScopeDocument:
		return Custom_RENDER_SCOPE_DOCUMENT
	default:
		slog.Error("unsupported renderer scope", "scope", scope)
		return Custom_RENDER_SCOPE_NODE
	}
}

func EncodeChildren(node *nodes.Node) []*Node {
	return utils.FnMap(node.GetChildren(), EncodeNode)
}

func EncodeMetadata(meta *nodes.FabricContentMetadata) *FabricContent_Metadata {
	if meta == nil {
		return nil
	}
	m := new(FabricContent_Metadata)
	m.SetProvider(meta.Provider)
	m.SetPlugin(meta.Plugin)
	m.SetVersion(meta.Version)
	return m
}

func EncodePath(path nodes.Path) *Path {
	if path == nil {
		return nil
	}
	p := new(Path)
	p.SetPath(utils.FnMap(path, func(i int) uint32 { return uint32(i) }))
	return p
}
