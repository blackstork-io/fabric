package plugin

import (
	"fmt"

	"github.com/blackstork-io/fabric/parser/definitions"
	"github.com/blackstork-io/fabric/plugin/ast"
	"github.com/blackstork-io/fabric/plugin/ast/nodes"
	"github.com/blackstork-io/fabric/plugin/plugindata"
)

type Content interface {
	plugindata.Convertible
	AsNode() *nodes.Node
	IsEmpty() bool
}

type ContentSectionOrDoc struct {
	Children   []Content
	pluginData plugindata.Map
	isDoc      bool
}

func NewSection(meta *definitions.MetaBlock, childCount int) *ContentSectionOrDoc {
	return newSection(meta, childCount)
}

func NewDocument(childCount int) *ContentSectionOrDoc {
	doc := newSection(nil, childCount)
	doc.isDoc = true
	return doc
}

func newSection(meta *definitions.MetaBlock, childCount int) *ContentSectionOrDoc {
	res := &ContentSectionOrDoc{
		Children: make([]Content, 0, childCount),
	}
	res.pluginData = plugindata.Map{
		"type":     plugindata.String("section"),
		"children": plugindata.List{},
	}
	if meta != nil {
		res.pluginData["meta"] = meta.AsPluginData()
	}
	return res
}

func (c *ContentSectionOrDoc) AppendChild(child Content) {
	if child == nil {
		return
	}
	c.Children = append(c.Children, child)
	c.pluginData["children"] = append(c.pluginData["children"].(plugindata.List), child.AsPluginData())
}

func (c *ContentSectionOrDoc) AsPluginData() plugindata.Data {
	if c == nil {
		return nil
	}
	return c.pluginData
}

// IsEmpty returns true if the section does not contain children
func (c *ContentSectionOrDoc) IsEmpty() bool {
	for _, child := range c.Children {
		if !child.IsEmpty() {
			return false
		}
	}
	return true
}

func (c *ContentSectionOrDoc) AsNode() *nodes.Node {
	if c == nil {
		return nil
	}
	var content *nodes.Node
	if c.isDoc {
		content = nodes.NewNode(&nodes.FabricDocument{})
	} else {
		content = nodes.NewNode(&nodes.FabricSection{})
	}
	for _, child := range c.Children {
		if child.IsEmpty() {
			continue
		}
		content.AppendChildren(child.AsNode())
	}
	return content
}

type ContentElement struct {
	node       *nodes.Node
	pluginData plugindata.Map
}

func NewElementFromNode(node *nodes.Node) *ContentElement {
	if _, ok := node.Content.(*nodes.FabricContent); !ok {
		return NewElement(node)
	}
	res := &ContentElement{
		node: node,
		pluginData: plugindata.Map{
			"type":     plugindata.String("element"),
			"markdown": plugindata.String(string(ast.AST2Md(node))),
		},
	}
	return res
}

// NewElement is the preferred way to create a new content element.
// It accepts a list of AST nodes to build the content element.
func NewElement(content ...*nodes.Node) *ContentElement {
	n := nodes.NewNode(&nodes.FabricContent{})
	n.SetChildren(content)
	return NewElementFromNode(n)
}

// NewElementFromMarkdown creates a new content element from a markdown string.
//
// Deprecated: opt in to working with the new AST by using [NewElement] instead.
func NewElementFromMarkdown(source string) *ContentElement {
	content := nodes.NewNode(&nodes.FabricContent{})
	content.SetChildren(ast.Markdown2AST([]byte(source)))
	return NewElementFromNode(content)
}

func (c *ContentElement) SetPluginMeta(meta *nodes.FabricContentMetadata) {
	if c == nil || meta == nil {
		return
	}
	c.pluginData["meta"] = meta.AsPluginData()
	c.node.Content.(*nodes.FabricContent).Meta = meta
	nodes.WalkContent(c.node, func(c *nodes.Custom, _ *nodes.Node, _ nodes.Path) {
		c.Data.TypeUrl = fmt.Sprintf(
			"%s%s/%s", nodes.CustomNodeTypeURLPrefix,
			meta.Plugin,
			c.GetStrippedNodeType(),
		)
	})
}

// AsMarkdownSrc returns the markdown source of the content element.
//
// Deprecated: opt in to working with the new AST by using .AsNode()
func (c *ContentElement) AsMarkdownSrc() string {
	return string(c.pluginData["markdown"].(plugindata.String))
}

func (c *ContentElement) AsNode() *nodes.Node {
	if c == nil {
		return nil
	}
	return c.node
}

func (c *ContentElement) IsEmpty() bool {
	return c == nil
}

var emptyContent = plugindata.Map{
	"type": plugindata.String("empty"),
	"meta": nil,
}

func (c *ContentElement) AsPluginData() plugindata.Data {
	if c == nil {
		return emptyContent
	}
	return c.pluginData
}
