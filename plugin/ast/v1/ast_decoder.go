package astv1

import (
	"fmt"
	"math"

	"github.com/yuin/goldmark/ast"
	east "github.com/yuin/goldmark/extension/ast"
	"github.com/yuin/goldmark/text"

	"github.com/blackstork-io/fabric/plugin/ast/astsrc"
	"github.com/blackstork-io/fabric/plugin/ast/nodes"
)

type DecoderOption interface {
	apply(dec *decoder)
}

type (
	AttributeDecoder func(*Attribute) (ast.Attribute, error)
	NodeDecoder      func(*Node) (ast.Node, error)
)

func (f AttributeDecoder) apply(dec *decoder) {
	if f != nil {
		dec.attributeDecoder = f
	}
}

func (f NodeDecoder) apply(dec *decoder) {
	if f != nil {
		dec.nodeDecoder = f
	}
}

func defaultNodeDecoder(node *Node) (ast.Node, error) {
	return nil, fmt.Errorf("%w: %T", ErrUnsupportedNodeType, node)
}

func Decode(root *Node, opts ...DecoderOption) (node ast.Node, source astsrc.ASTSource, err error) {
	defer func() {
		err = recoverBubbleError(err, recover())
	}()
	dec := &decoder{
		attributeDecoder: DefaultAttributeDecoder,
		nodeDecoder:      defaultNodeDecoder,
	}
	for _, opt := range opts {
		opt.apply(dec)
	}

	node = dec.decodeNode(root)
	source = dec.source
	return
}

type decoder struct {
	source           astsrc.ASTSource
	attributeDecoder AttributeDecoder
	nodeDecoder      NodeDecoder
}

func DefaultAttributeDecoder(attr *Attribute) (res ast.Attribute, err error) {
	switch val := attr.GetValue().(type) {
	case nil:
	case *Attribute_Bytes:
		res.Value = val.Bytes
	case *Attribute_Str:
		res.Value = val.Str
	default:
		err = fmt.Errorf("%w: %T", ErrUnsupportedAttributeType, val)
		return
	}
	res.Name = attr.GetName()
	return
}

func (d *decoder) decodeBaseNode(base *BaseNode, node ast.Node) {
	if base == nil {
		return
	}
	for _, encAttr := range base.GetAttributes() {
		attr, err := d.attributeDecoder(encAttr)
		if err != nil {
			bubbleUp(err)
		}
		node.SetAttribute(attr.Name, attr.Value)
	}

	for _, encChild := range base.GetChildren() {
		node.AppendChild(node, d.decodeNode(encChild))
	}
}

func (d *decoder) decodeText(txt *Text) *ast.Text {
	if txt == nil {
		return nil
	}
	res := ast.NewText()
	res.Segment = d.source.Append(txt.GetSegment())
	res.SetSoftLineBreak(txt.GetSoftLineBreak())
	res.SetHardLineBreak(txt.GetHardLineBreak())
	res.SetRaw(txt.GetRaw())
	return res
}

func (d *decoder) decodeNode(node *Node) (res ast.Node) {
	var base *BaseNode
	switch val := node.GetKind().(type) {
	case *Node_Document:
		base = val.Document.GetBase()
		res = ast.NewDocument()
		res.SetBlankPreviousLines(true)
	case *Node_TextBlock:
		base = val.TextBlock.GetBase()
		res = ast.NewTextBlock()
	case *Node_Paragraph:
		base = val.Paragraph.GetBase()
		res = ast.NewParagraph()
	case *Node_Heading:
		base = val.Heading.GetBase()
		res = ast.NewHeading(int(val.Heading.GetLevel()))
	case *Node_ThematicBreak:
		base = val.ThematicBreak.GetBase()
		res = ast.NewThematicBreak()
	case *Node_CodeBlock:
		base = val.CodeBlock.GetBase()
		codeBlock := ast.NewCodeBlock()
		codeBlock.SetLines(d.source.AppendMultiple(val.CodeBlock.GetLines()))
		res = codeBlock
	case *Node_FencedCodeBlock:
		base = val.FencedCodeBlock.GetBase()
		fencedCodeBlock := ast.NewFencedCodeBlock(d.decodeText(val.FencedCodeBlock.GetInfo()))
		fencedCodeBlock.SetLines(d.source.AppendMultiple(val.FencedCodeBlock.GetLines()))
		res = fencedCodeBlock
	case *Node_Blockquote:
		base = val.Blockquote.GetBase()
		res = ast.NewBlockquote()
	case *Node_List:
		base = val.List.GetBase()
		marker := val.List.GetMarker()
		if marker == 0 || marker > math.MaxUint8 {
			bubbleUp(fmt.Errorf("invalid marker character code: %d", marker))
		}
		list := ast.NewList(byte(marker))
		list.IsTight = val.List.GetIsTight()
		list.Start = int(val.List.GetStart())
		res = list
	case *Node_ListItem:
		base = val.ListItem.GetBase()
		res = ast.NewListItem(int(val.ListItem.GetOffset()))
	case *Node_HtmlBlock:
		base = val.HtmlBlock.GetBase()
		htmlBlock := ast.NewHTMLBlock(val.HtmlBlock.GetType().decode())
		if closure := val.HtmlBlock.GetClosureLine(); closure != nil {
			htmlBlock.ClosureLine = d.source.Append(val.HtmlBlock.GetClosureLine())
		} else {
			htmlBlock.ClosureLine = text.NewSegment(-1, -1)
		}
		htmlBlock.SetLines(d.source.AppendMultiple(val.HtmlBlock.GetLines()))
		res = htmlBlock
	case *Node_Text:
		base = val.Text.GetBase()
		res = d.decodeText(val.Text)
	case *Node_String_:
		base = val.String_.GetBase()
		str := ast.NewString(val.String_.GetValue())
		str.SetRaw(val.String_.GetRaw())
		str.SetCode(val.String_.GetCode())
		res = str
	case *Node_CodeSpan:
		base = val.CodeSpan.GetBase()
		res = ast.NewCodeSpan()
	case *Node_Emphasis:
		base = val.Emphasis.GetBase()
		res = ast.NewEmphasis(int(val.Emphasis.GetLevel()))
	case *Node_LinkOrImage:
		base = val.LinkOrImage.GetBase()
		tRes := ast.NewLink()
		tRes.Title = val.LinkOrImage.GetTitle()
		tRes.Destination = val.LinkOrImage.GetDestination()
		if val.LinkOrImage.GetIsImage() {
			res = ast.NewImage(tRes)
		} else {
			res = tRes
		}
	case *Node_AutoLink:
		base = val.AutoLink.GetBase()
		tRes := ast.NewAutoLink(val.AutoLink.GetType().decode(), d.decodeText(&Text{
			Segment: val.AutoLink.GetValue(),
		}))
		tRes.Protocol = val.AutoLink.GetProtocol()
		res = tRes

	case *Node_RawHtml:
		base = val.RawHtml.GetBase()
		tRes := ast.NewRawHTML()
		tRes.Segments = d.source.AppendMultiple(val.RawHtml.GetSegments())
		res = tRes

	case *Node_Table:
		base = val.Table.GetBase()
		tRes := east.NewTable()
		tRes.Alignments = decodeCellAlignments(val.Table.GetAlignments())
		res = tRes

	case *Node_TableRow:
		base = val.TableRow.GetBase()
		tRes := east.NewTableRow(decodeCellAlignments(val.TableRow.GetAlignments()))
		if val.TableRow.GetIsHeader() {
			res = east.NewTableHeader(tRes)
		} else {
			res = tRes
		}
	case *Node_TableCell:
		base = val.TableCell.GetBase()
		tRes := east.NewTableCell()
		tRes.Alignment = val.TableCell.GetAlignment().decode()
		res = tRes

	case *Node_Strikethrough:
		base = val.Strikethrough.GetBase()
		res = east.NewStrikethrough()

	case *Node_TaskCheckbox:
		base = val.TaskCheckbox.GetBase()
		res = east.NewTaskCheckBox(val.TaskCheckbox.GetIsChecked())

	case *Node_ContentNode:
		res = &nodes.FabricContentNode{
			Meta: DecodeMetadata(val.ContentNode.GetMetadata()),
		}
		base = val.ContentNode.GetRoot()
	case *Node_Custom:
		res = nodes.NewCustomNode(val.Custom.GetIsInline(), val.Custom.GetData())
		base = &BaseNode{
			BlankPreviousLines: val.Custom.GetBlankPreviousLines(),
		}
	default:
		bubbleUp(fmt.Errorf("Unsupported block kind: %T", node.GetKind()))
	}
	switch res.Type() {
	case ast.TypeBlock:
		res.SetBlankPreviousLines(base.GetBlankPreviousLines())
		fallthrough
	case ast.TypeInline, ast.TypeDocument:
		d.decodeBaseNode(base, res)
	default:
		bubbleUp(fmt.Errorf("Unsupported block type: %+v", res.Type()))
	}
	return
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
