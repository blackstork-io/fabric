package nodes

import (
	"fmt"
	"log/slog"
	"slices"
)

// Node is a node in the fabric AST.
type Node struct {
	Content  NodeContent
	children []*Node
	Parent   *Node
}

// NewNode creates a new node with the given content and children.
// Children could be another node, NodeContent instance or a string/[]byte
// (which will be converted to a Text node).
func NewNode(content NodeContent, children ...any) *Node {
	n := &Node{
		Content:  content,
		children: make([]*Node, 0, len(children)),
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
	return n.children
}

func (n *Node) SetChildren(children []*Node) {
	clear(n.children)
	n.children = n.children[:0]
	n.AppendChildren(children...)
}

func (n *Node) AppendChildren(children ...*Node) *Node {
	for _, c := range children {
		if n.validate(c) {
			c.Parent = n
			n.children = append(n.children, c)
		}
	}
	return n
}

func (n *Node) RemoveFromTree() *Node {
	if n.Parent != nil {
		idx := slices.Index(n.Parent.children, n)
		if idx >= 0 {
			n.Parent.children = slices.Delete(n.Parent.children, idx, idx+1)
		}
		n.Parent = nil
	}
	return n
}

func (n *Node) validate(child *Node) bool {
	if v, ok := n.Content.(ChildValidator); ok {
		if !v.ValidateChild(child) {
			return false
		}
	}
	if v, ok := child.Content.(ParentValidator); ok {
		if !v.ValidateParent(n) {
			return false
		}
	}
	return true
}
