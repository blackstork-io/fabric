package nodes

import "google.golang.org/protobuf/types/known/anypb"

// nodeContentSigil marks a struct as an AST node's content.
// It should be embedded in NodeKind structs
type nodeContentSigil struct{}

func (s *nodeContentSigil) isNodeContent() nodeContentSigil {
	return *s
}

// NodeContent is implemented by all possible node types in the fabric AST.
type NodeContent interface {
	isNodeContent() nodeContentSigil
}

type Emphasis struct {
	nodeContentSigil
	Level int
}

func Bold() *Emphasis {
	return &Emphasis{Level: 2}
}

func Italics() *Emphasis {
	return &Emphasis{Level: 1}
}

type Document struct {
	nodeContentSigil
	// TODO: add metadata for plugin-generated documents
}

type Paragraph struct {
	nodeContentSigil
	// Goldmark's TextBlocks are akin to Paragraphs, except
	// they are rendered without a wrapping <p> tag.
	IsTextBlock bool
}

type Heading struct {
	nodeContentSigil
	Level int
}

type ThematicBreak struct {
	nodeContentSigil
}

type CodeBlock struct {
	nodeContentSigil
	// Language is nil for indented code blocks
	// and empty []byte for fenced code blocks without language
	Language []byte
	Code     []byte
}

type Blockquote struct {
	nodeContentSigil
}

type List struct {
	nodeContentSigil
	Start  uint32
	Marker byte
	Items  [][]*Node
}

type HTMLBlock struct {
	nodeContentSigil
	HTML []byte
}

type Text struct {
	// Also covers String
	nodeContentSigil
	Text []byte
	// If true - the text ends with a hard line break
	HardLineBreak bool
}

type CodeSpan struct {
	nodeContentSigil
	Code []byte
}

type Link struct {
	nodeContentSigil
	Destination []byte
	Title       []byte
}

type Image struct {
	nodeContentSigil
	Source []byte
	Alt    []byte
}

// LinkOrImage is implemented by Link and Image.
// It is used to render links and images in a generic way.
type LinkOrImage interface {
	Url() []byte
	TitleOrAlt() []byte
}

func (l *Link) Url() []byte {
	return l.Destination
}

func (l *Link) TitleOrAlt() []byte {
	return l.Title
}

func (i *Image) Url() []byte {
	return i.Source
}

func (i *Image) TitleOrAlt() []byte {
	return i.Alt
}

type AutoLink struct {
	nodeContentSigil
	Value []byte
}

type HTMLInline struct {
	nodeContentSigil
	HTML []byte
}

type Table struct {
	nodeContentSigil
	Alignments []Alignment
	// Cells: [row][column][cell content nodes]
	Cells [][][]*Node
	// first row is always a header
}

// Alignment represents the alignment of a table cell.
type Alignment int

const (
	AlignmentNone Alignment = iota // must be 0, code elsewhere depends on this
	AlignmentLeft
	AlignmentCenter
	AlignmentRight
)

type TaskCheckbox struct {
	nodeContentSigil
	Checked bool
}

type Strikethrough struct {
	nodeContentSigil
}

// Custom content node.
// Should be converted to normal AST nodes by the plugin prior to rendering.
type Custom struct {
	nodeContentSigil
	Data *anypb.Any
}
