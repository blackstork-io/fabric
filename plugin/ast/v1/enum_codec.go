package astv1

import (
	"fmt"

	"github.com/yuin/goldmark/ast"
	east "github.com/yuin/goldmark/extension/ast"
)

func encodeHTMLBlockType(ty ast.HTMLBlockType) HTMLBlockType {
	switch ty {
	case ast.HTMLBlockType1:
		return HTMLBlockType_HTML_BLOCK_TYPE_1
	case ast.HTMLBlockType2:
		return HTMLBlockType_HTML_BLOCK_TYPE_2
	case ast.HTMLBlockType3:
		return HTMLBlockType_HTML_BLOCK_TYPE_3
	case ast.HTMLBlockType4:
		return HTMLBlockType_HTML_BLOCK_TYPE_4
	case ast.HTMLBlockType5:
		return HTMLBlockType_HTML_BLOCK_TYPE_5
	case ast.HTMLBlockType6:
		return HTMLBlockType_HTML_BLOCK_TYPE_6
	case ast.HTMLBlockType7:
		return HTMLBlockType_HTML_BLOCK_TYPE_7
	default:
		panic(bubbleWrap(fmt.Errorf("unsupported HTML block type: %v", ty)))
	}
}

func (bt HTMLBlockType) decode() ast.HTMLBlockType {
	switch bt {
	case HTMLBlockType_HTML_BLOCK_TYPE_1:
		return ast.HTMLBlockType1
	case HTMLBlockType_HTML_BLOCK_TYPE_2:
		return ast.HTMLBlockType2
	case HTMLBlockType_HTML_BLOCK_TYPE_3:
		return ast.HTMLBlockType3
	case HTMLBlockType_HTML_BLOCK_TYPE_4:
		return ast.HTMLBlockType4
	case HTMLBlockType_HTML_BLOCK_TYPE_5:
		return ast.HTMLBlockType5
	case HTMLBlockType_HTML_BLOCK_TYPE_6:
		return ast.HTMLBlockType6
	case HTMLBlockType_HTML_BLOCK_TYPE_7:
		return ast.HTMLBlockType7
	case HTMLBlockType_HTML_BLOCK_TYPE_UNSPECIFIED:
		fallthrough
	default:
		panic(bubbleWrap(fmt.Errorf("unsupported HTML block type: %v", bt)))
	}
}

func encodeAutoLinkType(ty ast.AutoLinkType) AutoLinkType {
	switch ty {
	case ast.AutoLinkURL:
		return AutoLinkType_AUTO_LINK_TYPE_URL
	case ast.AutoLinkEmail:
		return AutoLinkType_AUTO_LINK_TYPE_EMAIL
	default:
		panic(bubbleWrap(fmt.Errorf("unsupported auto link type: %v", ty)))
	}
}

func (ty AutoLinkType) decode() ast.AutoLinkType {
	switch ty {
	case AutoLinkType_AUTO_LINK_TYPE_URL:
		return ast.AutoLinkURL
	case AutoLinkType_AUTO_LINK_TYPE_EMAIL:
		return ast.AutoLinkEmail
	case AutoLinkType_AUTO_LINK_TYPE_UNSPECIFIED:
		fallthrough
	default:
		panic(bubbleWrap(fmt.Errorf("unsupported auto link type: %v", ty)))
	}
}

func encodeCellAlignment(alignment east.Alignment) CellAlignment {
	switch alignment {
	case east.AlignLeft:
		return CellAlignment_CELL_ALIGNMENT_LEFT
	case east.AlignRight:
		return CellAlignment_CELL_ALIGNMENT_RIGHT
	case east.AlignCenter:
		return CellAlignment_CELL_ALIGNMENT_CENTER
	case east.AlignNone:
		return CellAlignment_CELL_ALIGNMENT_NONE
	default:
		panic(bubbleWrap(fmt.Errorf("unsupported cell alignment: %v", alignment)))
	}
}

func (ca CellAlignment) decode() east.Alignment {
	switch ca {
	case CellAlignment_CELL_ALIGNMENT_NONE:
		return east.AlignNone
	case CellAlignment_CELL_ALIGNMENT_LEFT:
		return east.AlignLeft
	case CellAlignment_CELL_ALIGNMENT_RIGHT:
		return east.AlignRight
	case CellAlignment_CELL_ALIGNMENT_CENTER:
		return east.AlignCenter
	case CellAlignment_CELL_ALIGNMENT_UNSPECIFIED:
		fallthrough
	default:
		panic(bubbleWrap(fmt.Errorf("unsupported cell alignment: %v", ca)))
	}
}

func encodeCellAlignments(align []east.Alignment) []CellAlignment {
	res := make([]CellAlignment, len(align))
	for i, a := range align {
		res[i] = encodeCellAlignment(a)
	}
	return res
}

func decodeCellAlignments(align []CellAlignment) []east.Alignment {
	res := make([]east.Alignment, len(align))
	for i, a := range align {
		res[i] = a.decode()
	}
	return res
}
