package ast

import (
	"bytes"
	"regexp"

	astv1 "github.com/blackstork-io/fabric/plugin/ast/v1"
)

// inlines

func CodeSpan(code []byte) *astv1.Node_CodeSpan {
	return &astv1.Node_CodeSpan{
		CodeSpan: &astv1.CodeSpan{
			Base: &astv1.BaseNode{
				Children: []*astv1.Node{{
					Kind: &astv1.Node_Text{
						Text: &astv1.Text{
							Base:    &astv1.BaseNode{},
							Segment: code,
							Raw:     true,
						},
					},
				}},
			},
		},
	}
}

func emphasis(level int64, children []astv1.InlineContent) *astv1.Node_Emphasis {
	return &astv1.Node_Emphasis{
		Emphasis: &astv1.Emphasis{
			Base: &astv1.BaseNode{
				Children: astv1.Inlines.ExtendNodes(children, nil),
			},
			Level: level,
		},
	}
}

func Emphasis(children ...astv1.InlineContent) *astv1.Node_Emphasis {
	return emphasis(1, children)
}

var Italic = Emphasis

func StrongEmphasis(children ...astv1.InlineContent) *astv1.Node_Emphasis {
	return emphasis(2, children)
}

var Bold = StrongEmphasis

func Strikethrough(children ...astv1.InlineContent) *astv1.Node_Strikethrough {
	return &astv1.Node_Strikethrough{
		Strikethrough: &astv1.Strikethrough{
			Base: &astv1.BaseNode{
				Children: astv1.Inlines.ExtendNodes(children, nil),
			},
		},
	}
}

func Link(text ...astv1.InlineContent) *astv1.LinkOrImage {
	return &astv1.LinkOrImage{
		Base: &astv1.BaseNode{
			Children: astv1.Inlines.ExtendNodes(text, nil),
		},
		IsImage: false,
	}
}

func Image(alt ...astv1.InlineContent) *astv1.LinkOrImage {
	return &astv1.LinkOrImage{
		Base: &astv1.BaseNode{
			Children: astv1.Inlines.ExtendNodes(alt, nil),
		},
		IsImage: true,
	}
}

func AutoLink(url string) *astv1.Node_AutoLink {
	before, after, found := bytes.Cut([]byte(url), []byte("://"))
	if found {
		return &astv1.Node_AutoLink{
			AutoLink: &astv1.AutoLink{
				Base:     &astv1.BaseNode{},
				Type:     astv1.AutoLinkType_AUTO_LINK_TYPE_URL,
				Protocol: before,
				Value:    after,
			},
		}
	} else {
		return &astv1.Node_AutoLink{
			AutoLink: &astv1.AutoLink{
				Base:  &astv1.BaseNode{},
				Type:  astv1.AutoLinkType_AUTO_LINK_TYPE_EMAIL,
				Value: before,
			},
		}
	}
}

func AutoLinkEmail(email string) *astv1.Node_AutoLink {
	return &astv1.Node_AutoLink{
		AutoLink: &astv1.AutoLink{
			Base:  &astv1.BaseNode{},
			Type:  astv1.AutoLinkType_AUTO_LINK_TYPE_EMAIL,
			Value: []byte(email),
		},
	}
}

func InlineHTML(html string) *astv1.Node_RawHtml {
	return &astv1.Node_RawHtml{
		RawHtml: &astv1.RawHTML{
			Base:     &astv1.BaseNode{},
			Segments: [][]byte{[]byte(html)},
		},
	}
}

func LineBreak() *astv1.Node_Text {
	return &astv1.Node_Text{
		Text: &astv1.Text{
			Base:          &astv1.BaseNode{},
			Segment:       nil,
			HardLineBreak: true,
		},
	}
}

var textHardBreakRegexp = regexp.MustCompile(`( {2,}\n *|\\\n *)`)

func splitBytes(re *regexp.Regexp, s []byte, n int) [][]byte {
	if n == 0 {
		return nil
	}

	if len(s) == 0 {
		return [][]byte{nil}
	}

	matches := re.FindAllIndex(s, n)
	subBytes := make([][]byte, 0, len(matches))

	beg := 0
	end := 0
	for _, match := range matches {
		if n > 0 && len(subBytes) >= n-1 {
			break
		}

		end = match[0]
		if match[1] != 0 {
			subBytes = append(subBytes, s[beg:end])
		}
		beg = match[1]
	}

	if end != len(s) {
		subBytes = append(subBytes, s[beg:])
	}

	return subBytes
}

func convertSoftBreaks(txt []byte) (res []*astv1.Text) {
	split := bytes.Split(txt, []byte("\n"))
	res = make([]*astv1.Text, 0, len(split))
	res = append(res, &astv1.Text{
		Segment: split[0],
	})
	for i, s := range split[1:] {
		res[i].SoftLineBreak = true
		res = append(res, &astv1.Text{
			Segment: s,
		})
	}
	return
}

func convertHardBreaks(txt []byte) (res []*astv1.Text) {
	split := splitBytes(textHardBreakRegexp, txt, -1)
	res = make([]*astv1.Text, 0, len(split))
	res = append(res, convertSoftBreaks(split[0])...)
	for i, s := range split[1:] {
		res[i].HardLineBreak = true
		res = append(res, convertSoftBreaks(s)...)
	}
	return
}

func Text(text string) (res astv1.Inlines) {
	txt := convertHardBreaks([]byte(text))
	res = make(astv1.Inlines, len(txt))
	for i, t := range txt {
		res[i] = &astv1.Node_Text{
			Text: t,
		}
	}
	return
}

// container blocks
func Blockquote(children ...astv1.BlockContent) *astv1.Node_Blockquote {
	return &astv1.Node_Blockquote{
		Blockquote: &astv1.Blockquote{
			Base: &astv1.BaseNode{
				Children: astv1.Blocks.ExtendNodes(children, nil),
			},
		},
	}
}

type listMarker uint32

const (
	Period listMarker = '.'
	Paren  listMarker = ')'
	Star   listMarker = '*'
	Plus   listMarker = '+'
	Hyphen listMarker = '-'
)

func List(marker listMarker) *astv1.Node_List {
	return &astv1.Node_List{
		List: &astv1.List{
			Base:   &astv1.BaseNode{},
			Marker: uint32(marker),
		},
	}
}

func ThematicBreak() *astv1.Node_ThematicBreak {
	return &astv1.Node_ThematicBreak{
		ThematicBreak: &astv1.ThematicBreak{
			Base: &astv1.BaseNode{},
		},
	}
}

func Header(level uint32, children ...astv1.InlineContent) *astv1.Node_Heading {
	return &astv1.Node_Heading{
		Heading: &astv1.Heading{
			Base: &astv1.BaseNode{
				Children: astv1.Inlines.ExtendNodes(children, nil),
			},
			Level: max(1, min(level, 6)),
		},
	}
}

func IndentedCodeBlock(code []byte) *astv1.Node_CodeBlock {
	return &astv1.Node_CodeBlock{
		CodeBlock: &astv1.CodeBlock{
			Base:  &astv1.BaseNode{},
			Lines: bytes.Split(code, []byte("\n")),
		},
	}
}

func FencedCodeBlock(code []byte) *astv1.Node_FencedCodeBlock {
	return &astv1.Node_FencedCodeBlock{
		FencedCodeBlock: &astv1.FencedCodeBlock{
			Base:  &astv1.BaseNode{},
			Info:  nil,
			Lines: bytes.Split(code, []byte("\n")),
		},
	}
}

func HTMLBlock(html []byte) *astv1.Node_HtmlBlock {
	return &astv1.Node_HtmlBlock{
		HtmlBlock: &astv1.HTMLBlock{
			Base: &astv1.BaseNode{},
			// setting the value to most general type, hopefully renderers don't care
			Type:  astv1.HTMLBlockType_HTML_BLOCK_TYPE_7,
			Lines: bytes.Split(html, []byte("\n")),
		},
	}
}

func Paragraph(children ...astv1.InlineContent) *astv1.Node_Paragraph {
	return &astv1.Node_Paragraph{
		Paragraph: &astv1.Paragraph{
			Base: &astv1.BaseNode{
				Children: astv1.Inlines.ExtendNodes(children, nil),
			},
		},
	}
}

var (
	AlignLeft   = astv1.CellAlignment_CELL_ALIGNMENT_LEFT
	AlignRight  = astv1.CellAlignment_CELL_ALIGNMENT_RIGHT
	AlignCenter = astv1.CellAlignment_CELL_ALIGNMENT_CENTER
	AlignNone   = astv1.CellAlignment_CELL_ALIGNMENT_NONE
)

func Table() *astv1.Node_Table {
	return &astv1.Node_Table{
		Table: &astv1.Table{
			Base: &astv1.BaseNode{},
		},
	}
}
