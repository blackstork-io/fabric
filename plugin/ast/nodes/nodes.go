package nodes

import (
	"fmt"
	"log/slog"
	"slices"
)

// Node is a node in the fabric AST.
type Node struct {
	Content  NodeContent
	Children []*Node
	Parent   *Node
}

// NewNode creates a new node with the given content and children.
// Children could be another node, NodeContent instance or a string/[]byte
// (which will be converted to a Text node).
func NewNode(content NodeContent, children ...any) *Node {
	n := &Node{
		Content:  content,
		Children: make([]*Node, 0, len(children)),
	}
	for _, c := range children {
		var child *Node
		switch c := c.(type) {
		case *Node:
			child = c
		case string:
			child = &Node{
				Content: &Text{Text: []byte(c)},
			}
		case []byte:
			child = &Node{
				Content: &Text{Text: c},
			}
		case NodeContent:
			child = &Node{
				Content: c,
			}
		default:
			slog.Error("invalid node child type", "type", fmt.Sprintf("%T", c))
			continue
		}
		n.AppendChildren(child)
	}
	return n
}

func (n *Node) GetChildren() []*Node {
	return n.Children
}

func (n *Node) SetChildren(children []*Node) {
	for _, c := range children {
		c.Parent = n
	}
	n.Children = children
}

func (n *Node) AppendChildren(children ...*Node) {
	for _, c := range children {
		c.Parent = n
	}
	n.Children = append(n.Children, children...)
}

func (n *Node) RemoveFromTree() *Node {
	if n.Parent != nil {
		n.Parent.Children = slices.DeleteFunc(n.Parent.Children, func(c *Node) bool {
			return c == n
		})
		n.Parent = nil
	}
	return n
}
