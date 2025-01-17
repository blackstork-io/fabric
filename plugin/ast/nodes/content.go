package nodes

import (
	"fmt"
	"log/slog"
	"strings"

	"google.golang.org/protobuf/types/known/anypb"

	"github.com/blackstork-io/fabric/plugin/plugindata"
)

// nodeContentSigil marks a struct as an AST node's content.
// It should be embedded in NodeKind structs
type nodeContentSigil struct{}

func (s *nodeContentSigil) isNodeContent() nodeContentSigil {
	return *s
}

func (s *nodeContentSigil) ValidateParent(parent *Node) bool {
	return true
}

func (s *nodeContentSigil) ValidateChild(child *Node) bool {
	return true
}

// NodeContent is implemented by all possible node types in the fabric AST.
type NodeContent interface {
	isNodeContent() nodeContentSigil
	ValidateParent(parent *Node) bool
	ValidateChild(child *Node) bool
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

func (s *ThematicBreak) ValidateChild(child *Node) bool {
	slog.Error(
		"ThematicBreak cannot have children, attempted to add",
		"child", fmt.Sprintf("%T", child.Content),
	)
	return false
}

type CodeBlock struct {
	nodeContentSigil
	// Language is nil for indented code blocks
	// and empty []byte for fenced code blocks without language
	Language []byte
	Code     []byte
}

func (s *CodeBlock) ValidateChild(child *Node) bool {
	slog.Error(
		"CodeBlock cannot have children, attempted to add",
		"child", fmt.Sprintf("%T", child.Content),
	)
	return false
}

type Blockquote struct {
	nodeContentSigil
}

type List struct {
	nodeContentSigil
	Start  uint32
	Marker byte
}

func (s *List) ValidateChild(child *Node) bool {
	if _, ok := child.Content.(*ListItem); !ok {
		slog.Error(
			"List can only contain ListItems, attempted to add",
			"child", fmt.Sprintf("%T", child.Content),
		)
		return false
	}
	return true
}

type ListItem struct {
	nodeContentSigil
}

func (s *ListItem) ValidateParent(parent *Node) bool {
	if _, ok := parent.Content.(*List); !ok {
		slog.Error(
			"ListItem can only be contained by List, attempted to add to",
			"parent", fmt.Sprintf("%T", parent),
		)
		return false
	}
	return true
}

type HTMLBlock struct {
	nodeContentSigil
	HTML []byte
}

func (s *HTMLBlock) ValidateChild(child *Node) bool {
	slog.Error(
		"HTMLBlock cannot have children, attempted to add",
		"child", fmt.Sprintf("%T", child.Content),
	)
	return false
}

type Text struct {
	// Also covers String
	nodeContentSigil
	Text []byte
	// If true - the text ends with a hard line break
	HardLineBreak bool
}

func (s *Text) ValidateChild(child *Node) bool {
	slog.Error(
		"Text cannot have children, attempted to add",
		"child", fmt.Sprintf("%T", child.Content),
	)
	return false
}

type CodeSpan struct {
	nodeContentSigil
	Code []byte
}

func (s *CodeSpan) ValidateChild(child *Node) bool {
	slog.Error(
		"CodeSpan cannot have children, attempted to add",
		"child", fmt.Sprintf("%T", child.Content),
	)
	return false
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

func (s *AutoLink) ValidateChild(child *Node) bool {
	slog.Error(
		"AutoLink cannot have children, attempted to add",
		"child", fmt.Sprintf("%T", child.Content),
	)
	return false
}

type HTMLInline struct {
	nodeContentSigil
	HTML []byte
}

func (s *HTMLInline) ValidateChild(child *Node) bool {
	slog.Error(
		"HTMLInline cannot have children, attempted to add",
		"child", fmt.Sprintf("%T", child.Content),
	)
	return false
}

type Table struct {
	nodeContentSigil
	Alignments []Alignment
	// first row is always a header
}

func (s *Table) ValidateChild(child *Node) bool {
	if _, ok := child.Content.(*TableRow); !ok {
		slog.Error(
			"Table can only contain TableRows, attempted to add",
			"child", fmt.Sprintf("%T", child.Content),
		)
		return false
	}
	return true
}

type TableRow struct {
	nodeContentSigil
}

func (s *TableRow) ValidateParent(parent *Node) bool {
	if _, ok := parent.Content.(*Table); !ok {
		slog.Error(
			"TableRow can only be contained by Table, attempted to add to",
			"parent", fmt.Sprintf("%T", parent),
		)
		return false
	}
	return true
}

func (s *TableRow) ValidateChild(child *Node) bool {
	if _, ok := child.Content.(*TableCell); !ok {
		slog.Error(
			"TableRow can only contain TableCells, attempted to add",
			"child", fmt.Sprintf("%T", child.Content),
		)
		return false
	}
	return true
}

type TableCell struct {
	nodeContentSigil
}

func (s *TableCell) ValidateParent(parent *Node) bool {
	if _, ok := parent.Content.(*TableRow); !ok {
		slog.Error(
			"TableCell can only be contained by TableRow, attempted to add to",
			"parent", fmt.Sprintf("%T", parent),
		)
		return false
	}
	return true
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

func (s *TaskCheckbox) ValidateParent(parent *Node) bool {
	if _, ok := parent.Content.(*ListItem); !ok {
		slog.Error(
			"TaskCheckbox can only be contained by ListItem, attempted to add to",
			"parent", fmt.Sprintf("%T", parent),
		)
		return false
	}
	if len(parent.children) > 0 {
		slog.Error(
			"TaskCheckbox must be the first child of a ListItem",
		)
		return false
	}
	return true
}

type Strikethrough struct {
	nodeContentSigil
}

type RendererScope int

const (
	// Renderer receives only the current node
	ScopeNode RendererScope = iota
	// Renderer receives the entire parent content node
	ScopeContent
	// Renderer receives the entire parent section or document (if not within a section)
	ScopeSection
	// Renderer receives the entire parent document
	ScopeDocument
)

// Custom content node.
// Should be converted to normal AST nodes by the plugin prior to rendering.
type Custom struct {
	nodeContentSigil
	Data  *anypb.Any
	Scope RendererScope
}

// Prefix format: types.blackstork.io/fabric/v1/custom_nodes/<plugin_name>/<node_type>
const CustomNodeTypeURLPrefix = "types.blackstork.io/fabric/v1/custom_nodes/"

// GetStrippedNodeType returns local component of the custom node type url
func (c *Custom) GetStrippedNodeType() string {
	if c == nil {
		return ""
	}
	typeUrl := c.Data.GetTypeUrl()
	if strings.HasPrefix(typeUrl, CustomNodeTypeURLPrefix) {
		_, typeUrl, _ = strings.Cut(typeUrl[len(CustomNodeTypeURLPrefix):], "/")
	}
	return strings.Trim(typeUrl, "/")
}

// GetPluginName returns the plugin name from the custom node type url
func (c *Custom) GetPluginName() string {
	if c == nil {
		return ""
	}
	typeUrl := c.Data.GetTypeUrl()
	if strings.HasPrefix(typeUrl, CustomNodeTypeURLPrefix) {
		pluginName, _, _ := strings.Cut(typeUrl[len(CustomNodeTypeURLPrefix):], "/")
		return strings.Trim(pluginName, "/")
	}
	return ""
}

type FabricDocument struct {
	nodeContentSigil
}

func (s *FabricDocument) ValidateParent(parent *Node) bool {
	slog.Warn(
		"FabricDocument should not be nested, attempted to add to",
		"parent", fmt.Sprintf("%T", parent),
	)
	return true
}

func (s *FabricDocument) ValidateChild(child *Node) bool {
	switch child.Content.(type) {
	case *FabricSection, *FabricContent:
	default:
		slog.Warn(
			"FabricDocument can only contain FabricSections and FabricContent, attempted to add",
			"child", fmt.Sprintf("%T", child.Content),
		)
	}
	return true
}

type FabricSection struct {
	nodeContentSigil
}

func (s *FabricSection) ValidateParent(parent *Node) bool {
	switch parent.Content.(type) {
	case *FabricDocument, *FabricSection:
	default:
		slog.Warn("FabricSection can only be contained by FabricDocument or FabricSection, attempted to add to",
			"parent", fmt.Sprintf("%T", parent.Content),
		)
	}
	return true
}

func (s *FabricSection) ValidateChild(child *Node) bool {
	switch child.Content.(type) {
	case *FabricSection, *FabricContent:
	default:
		slog.Warn(
			"FabricSection can only contain FabricSections and FabricContent, attempted to add",
			"child", fmt.Sprintf("%T", child.Content),
		)
	}
	return true
}

type FabricContent struct {
	nodeContentSigil
	Meta *FabricContentMetadata
}

type FabricContentMetadata struct {
	// ie "blackstork/builtin"
	Provider string
	// ie "title"
	Plugin  string
	Version string
}

var _ plugindata.Convertible = (*FabricContentMetadata)(nil)

func (m *FabricContentMetadata) AsPluginData() plugindata.Data {
	if m == nil {
		return nil
	}
	return plugindata.Map{
		"provider": plugindata.String(m.Provider),
		"plugin":   plugindata.String(m.Plugin),
		"version":  plugindata.String(m.Version),
	}
}

func (c *FabricContent) ValidateParent(parent *Node) bool {
	switch parent.Content.(type) {
	case *FabricDocument, *FabricSection:
	default:
		slog.Warn("FabricContent can only be contained by FabricDocument or FabricSection, attempted to add to",
			"parent", fmt.Sprintf("%T", parent.Content),
		)
	}
	return true
}
