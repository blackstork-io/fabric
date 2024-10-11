package astv1

import (
	"errors"
	"fmt"

	"github.com/yuin/goldmark/ast"
	east "github.com/yuin/goldmark/extension/ast"
	"github.com/yuin/goldmark/text"

	"github.com/blackstork-io/fabric/plugin/ast/nodes"
)

// ErrSkipAttribute should be returned from AttributeEncoder if the
// attribute needs to be be skipped.
var ErrSkipAttribute = errors.New("skip")

var (
	ErrUnsupportedNodeType      = errors.New("unsupported node type")
	ErrUnsupportedAttributeType = errors.New("unsupported attribute type")
	ErrUnsupportedAlignment     = errors.New("unsupported alignment")
)

type (
	AttributeEncoder func(*ast.Attribute) (*Attribute, error)
	NodeEncoder      func(ast.Node) (isNode_Kind, error)
)

func (f AttributeEncoder) apply(enc *encoder) {
	if f != nil {
		enc.attributeEncoder = f
	}
}

func (f NodeEncoder) apply(enc *encoder) {
	if f != nil {
		enc.nodeEncoder = f
	}
}

type EncoderOption interface {
	apply(enc *encoder)
}

type encoder struct {
	attributeEncoder AttributeEncoder
	nodeEncoder      NodeEncoder
	source           []byte
}

func defaultNodeEncoder(node ast.Node) (isNode_Kind, error) {
	return nil, fmt.Errorf("%w: %s", ErrUnsupportedNodeType, node.Kind())
}

func Encode(root ast.Node, source []byte, opts ...EncoderOption) (node *Node, err error) {
	defer func() {
		err = recoverBubbleError(err, recover())
	}()
	enc := &encoder{
		attributeEncoder: DefaultAttributeEncoder,
		source:           source,
		nodeEncoder:      defaultNodeEncoder,
	}
	for _, opt := range opts {
		opt.apply(enc)
	}
	node = enc.encodeNode(root)
	return
}

func (e *encoder) encodeSegment(seg text.Segment) []byte {
	if seg.Start > seg.Stop || seg.Start < 0 || seg.Stop > len(e.source) {
		// Invalid segment.
		return nil
	}
	// Note: zero-length segments (seg.Start == seg.Stop) are valid,
	// and even could be meaningful if padding is > 0.
	// Render them using seg.Value.
	return seg.Value(e.source)
}

func (e *encoder) encodeBaseNode(node *ast.BaseNode) *BaseNode {
	var res BaseNode
	attrs := node.Attributes()
	res.Attributes = make([]*Attribute, 0, len(attrs))
	for _, attr := range attrs {
		encoded, err := e.attributeEncoder(&attr)
		if err != nil {
			bubbleUp(err)
		}
		res.Attributes = append(res.Attributes, encoded)
	}
	res.Children = make([]*Node, 0, node.ChildCount())
	for child := node.FirstChild(); child != nil; child = child.NextSibling() {
		res.Children = append(res.Children, e.encodeNode(child))
	}
	return &res
}

func (e *encoder) encodeBaseBlock(node *ast.BaseBlock) (res *BaseNode) {
	res = e.encodeBaseNode(&node.BaseNode)
	res.BlankPreviousLines = node.HasBlankPreviousLines()
	return
}

func (e *encoder) encodeNode(node ast.Node) *Node {
	var kind isNode_Kind

	switch n := node.(type) {
	case *ast.Document:
		kind = &Node_Document{
			Document: &Document{
				Base: e.encodeBaseBlock(&n.BaseBlock),
			},
		}
	case *ast.TextBlock:
		kind = &Node_TextBlock{
			TextBlock: &TextBlock{
				Base: e.encodeBaseBlock(&n.BaseBlock),
			},
		}
	case *ast.Paragraph:
		kind = &Node_Paragraph{
			Paragraph: &Paragraph{
				Base: e.encodeBaseBlock(&n.BaseBlock),
			},
		}
	case *ast.Heading:
		kind = &Node_Heading{
			Heading: &Heading{
				Base:  e.encodeBaseBlock(&n.BaseBlock),
				Level: uint32(n.Level), //nolint:gosec // Level is bounded by 1-6.
			},
		}
	case *ast.ThematicBreak:
		kind = &Node_ThematicBreak{
			ThematicBreak: &ThematicBreak{
				Base: e.encodeBaseBlock(&n.BaseBlock),
			},
		}
	case *ast.CodeBlock:
		kind = &Node_CodeBlock{
			CodeBlock: &CodeBlock{
				Base:  e.encodeBaseBlock(&n.BaseBlock),
				Lines: e.encodeSegments(n.Lines()),
			},
		}
	case *ast.FencedCodeBlock:
		var txt *Text
		if n.Info != nil {
			txt = &Text{
				Segment: e.encodeSegment(n.Info.Segment),
				Raw:     n.Info.IsRaw(),
			}
		}
		kind = &Node_FencedCodeBlock{
			FencedCodeBlock: &FencedCodeBlock{
				Base:  e.encodeBaseBlock(&n.BaseBlock),
				Info:  txt,
				Lines: e.encodeSegments(n.Lines()),
			},
		}
	case *ast.Blockquote:
		kind = &Node_Blockquote{
			Blockquote: &Blockquote{
				Base: e.encodeBaseBlock(&n.BaseBlock),
			},
		}
	case *ast.List:
		kind = &Node_List{
			List: &List{
				Base:    e.encodeBaseBlock(&n.BaseBlock),
				Marker:  uint32(n.Marker),
				IsTight: n.IsTight,
				Start:   uint32(n.Start), //nolint:gosec // Start numbers must be nine digits or less (https://spec.commonmark.org/0.31.2/#example-265)
			},
		}
	case *ast.ListItem:
		kind = &Node_ListItem{
			ListItem: &ListItem{
				Base:   e.encodeBaseBlock(&n.BaseBlock),
				Offset: int64(n.Offset),
			},
		}
	case *ast.HTMLBlock:
		// Doing encoding manually to prevent a change in the internal representation
		// of HTMLBlockType enum from breaking our encoding.
		kind = &Node_HtmlBlock{
			HtmlBlock: &HTMLBlock{
				Base:        e.encodeBaseBlock(&n.BaseBlock),
				Type:        encodeHTMLBlockType(n.HTMLBlockType),
				Lines:       e.encodeSegments(n.Lines()),
				ClosureLine: e.encodeSegment(n.ClosureLine),
			},
		}
	case *ast.Text:
		kind = &Node_Text{
			Text: &Text{
				Base:          e.encodeBaseNode(&n.BaseNode),
				Segment:       e.encodeSegment(n.Segment),
				SoftLineBreak: n.SoftLineBreak(),
				HardLineBreak: n.HardLineBreak(),
				Raw:           n.IsRaw(),
			},
		}
	case *ast.String:
		kind = &Node_String_{
			String_: &String{
				Base:  e.encodeBaseNode(&n.BaseNode),
				Value: n.Value,
				Raw:   n.IsRaw(),
				Code:  n.IsCode(),
			},
		}
	case *ast.CodeSpan:
		kind = &Node_CodeSpan{
			CodeSpan: &CodeSpan{
				Base: e.encodeBaseNode(&n.BaseNode),
			},
		}
	case *ast.Emphasis:
		kind = &Node_Emphasis{
			Emphasis: &Emphasis{
				Base:  e.encodeBaseNode(&n.BaseNode),
				Level: int64(n.Level),
			},
		}
	case *ast.Link:
		kind = &Node_LinkOrImage{
			LinkOrImage: &LinkOrImage{
				Base:        e.encodeBaseNode(&n.BaseNode),
				Destination: n.Destination,
				Title:       n.Title,
				IsImage:     false,
			},
		}
	case *ast.Image:
		kind = &Node_LinkOrImage{
			LinkOrImage: &LinkOrImage{
				Base:        e.encodeBaseNode(&n.BaseNode),
				Destination: n.Destination,
				Title:       n.Title,
				IsImage:     true,
			},
		}
	case *ast.AutoLink:
		kind = &Node_AutoLink{
			AutoLink: &AutoLink{
				Base:     e.encodeBaseNode(&n.BaseNode),
				Type:     encodeAutoLinkType(n.AutoLinkType),
				Protocol: n.Protocol,
				Value:    n.Label(e.source),
			},
		}
	case *ast.RawHTML:
		kind = &Node_RawHtml{
			RawHtml: &RawHTML{
				Base:     e.encodeBaseNode(&n.BaseNode),
				Segments: e.encodeSegments(n.Segments),
			},
		}
	// GitHub Flavored Markdown
	case *east.Table:
		kind = &Node_Table{
			Table: &Table{
				Base:       e.encodeBaseBlock(&n.BaseBlock),
				Alignments: encodeCellAlignments(n.Alignments),
			},
		}

	case *east.TableRow:
		kind = &Node_TableRow{
			TableRow: &TableRow{
				Base:       e.encodeBaseBlock(&n.BaseBlock),
				Alignments: encodeCellAlignments(n.Alignments),
				IsHeader:   false,
			},
		}
	case *east.TableHeader:
		kind = &Node_TableRow{
			TableRow: &TableRow{
				Base:       e.encodeBaseBlock(&n.BaseBlock),
				Alignments: encodeCellAlignments(n.Alignments),
				IsHeader:   true,
			},
		}
	case *east.TableCell:
		kind = &Node_TableCell{
			TableCell: &TableCell{
				Base:      e.encodeBaseBlock(&n.BaseBlock),
				Alignment: encodeCellAlignment(n.Alignment),
			},
		}
	case *east.Strikethrough:
		kind = &Node_Strikethrough{
			Strikethrough: &Strikethrough{
				Base: e.encodeBaseNode(&n.BaseNode),
			},
		}
	case *east.TaskCheckBox:
		kind = &Node_TaskCheckbox{
			TaskCheckbox: &TaskCheckbox{
				Base:      e.encodeBaseNode(&n.BaseNode),
				IsChecked: n.IsChecked,
			},
		}
	case *nodes.FabricContentNode:
		kind = &Node_ContentNode{
			ContentNode: &FabricContentNode{
				Metadata: EncodeMetadata(n.Meta),
				Root:     e.encodeBaseBlock(&n.BaseBlock),
			},
		}
	case *nodes.CustomBlock:
		kind = &Node_Custom{
			Custom: &CustomNode{
				IsInline: false,
				Data:     n.Data,
			},
		}
	case *nodes.CustomInline:
		kind = &Node_Custom{
			Custom: &CustomNode{
				IsInline: true,
				Data:     n.Data,
			},
		}
	default:
		var err error
		kind, err = e.nodeEncoder(node)
		if err != nil {
			bubbleUp(err)
		}
	}

	return &Node{
		Kind: kind,
	}
}

func (e *encoder) encodeSegments(seg *text.Segments) [][]byte {
	res := make([][]byte, seg.Len())
	for i := range res {
		res[i] = e.encodeSegment(seg.At(i))
	}
	return res
}

func DefaultAttributeEncoder(attr *ast.Attribute) (*Attribute, error) {
	var res isAttribute_Value
	switch val := attr.Value.(type) {
	case []byte:
		res = &Attribute_Bytes{
			Bytes: val,
		}
	case string:
		res = &Attribute_Str{
			Str: val,
		}
	case uint, uint8, uint16, uint32, uintptr, uint64, int, int8, int16, int32, int64:
		res = &Attribute_Str{
			Str: fmt.Sprintf("%d", val),
		}
	case bool:
		res = &Attribute_Bytes{
			Bytes: attr.Name,
		}
	case float32, float64:
		res = &Attribute_Str{
			Str: fmt.Sprintf("%f", val),
		}
	default:
		return nil, fmt.Errorf("%w: %T", ErrUnsupportedAttributeType, attr.Value)
	}
	return &Attribute{
		Name:  attr.Name,
		Value: res,
	}, nil
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
