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
	parent   *Node
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

func (n *Node) GetRelativePath(relativeRoot *Node) (path Path) {
	for ; n != nil && n != relativeRoot; n = n.parent {
		path = append(path, n.indexInParent())
	}
	if n != relativeRoot {
		slog.Warn("node is not a descendant of the relative root")
		return nil
	}
	slices.Reverse(path)
	return path
}

func (n *Node) GetAbsolutePath() Path {
	return n.GetRelativePath(nil)
}

// TraversePath traverses the path from the given node.
func (n *Node) TraversePath(path Path) *Node {
	if n == nil {
		return nil
	}
	for _, step := range path {
		if step < 0 || step >= len(n.children) || n == nil {
			return nil
		}
		n = n.children[step]
	}
	return n
}

func (n *Node) SetChildren(children []*Node) {
	for i := range n.children {
		n.children[i].parent = nil
		n.children[i] = nil
	}
	n.children = n.children[:0]
	n.AppendChildren(children...)
}

// AppendChildren appends the given nodes as children of the node.
func (n *Node) AppendChildren(children ...*Node) *Node {
	for _, c := range children {
		if n.validate(c) {
			c.parent = n
			n.children = append(n.children, c)
		}
	}
	return n
}

// InsertAfter inserts the given nodes after the node (as siblings).
// Returns the current node.
func (n *Node) InsertAfter(nodes ...*Node) *Node {
	return n.insert(true, nodes...)
}

// InsertBefore inserts the given nodes before the node (as siblings).
// Returns the current node.
func (n *Node) InsertBefore(nodes ...*Node) *Node {
	return n.insert(false, nodes...)
}

func (n *Node) insert(after bool, nodes ...*Node) *Node {
	if len(nodes) == 0 {
		return n
	}
	idx := n.indexInParent()
	if idx < 0 {
		panic("node is not in the tree")
	}
	if after {
		idx++
	}
	nextSiblings := slices.Clone(n.parent.children[idx:])
	n.parent.children = n.parent.children[:idx]
	n.parent.AppendChildren(nodes...)
	n.parent.AppendChildren(nextSiblings...)
	return n
}

// ReplaceWith replaces the node with the given nodes.
// Returns the replaced node.
func (n *Node) ReplaceWith(nodes ...*Node) *Node {
	idx := n.indexInParent()
	if idx < 0 {
		panic("node is not in the tree")
	}
	nextSiblings := slices.Clone(n.parent.children[idx+1:])
	n.parent.children = n.parent.children[:idx]
	n.parent.AppendChildren(nodes...)
	n.parent.AppendChildren(nextSiblings...)
	n.parent = nil
	return n
}

// NextSibling returns the next sibling of the node or nil if it is the last child.
func (n *Node) NextSibling() *Node {
	idx := n.indexInParent()
	if idx < 0 {
		return nil
	}
	return n.parent.children[idx+1]
}

// PrevSibling returns the previous sibling of the node or nil if it is the first child.
func (n *Node) PrevSibling() *Node {
	idx := n.indexInParent()
	if idx <= 0 {
		return nil
	}
	return n.parent.children[idx-1]
}

// RemoveFromTree removes the node from the tree, returning the node itself.
func (n *Node) RemoveFromTree() *Node {
	return n.ReplaceWith()
}

// IsEmpty returns true if the node has no children.
func (n *Node) IsEmpty() bool {
	return n == nil || len(n.children) == 0
}

func (n *Node) indexInParent() int {
	if n == nil || n.parent == nil {
		return -1
	}
	return slices.Index(n.parent.children, n)
}

func (n *Node) validate(child *Node) bool {
	return n.Content.ValidateChild(child) && child.Content.ValidateParent(n)
}

func (n *Node) Walk(fn func(*Node, Path)) {
	n.walk(Path{}, fn)
}

func (n *Node) walk(p Path, fn func(*Node, Path)) {
	fn(n, p)
	idx := len(p)
	p = append(p, 0)
	for i, c := range n.children {
		p[idx] = i
		c.walk(p, fn)
	}
}

func WalkContent[T NodeContent](n *Node, fn func(T, *Node, Path)) {
	n.Walk(func(n *Node, p Path) {
		if c, ok := n.Content.(T); ok {
			fn(c, n, p)
		}
	})
}
