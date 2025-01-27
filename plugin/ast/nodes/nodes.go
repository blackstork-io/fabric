package nodes

import (
	"fmt"
	"log/slog"
	"slices"

	"iter"
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

// Clone creates a deep copy of the whole tree and a shallow copy of the content.
// Changes to the node (eg replacing children) will not affect the original node.
// Changes to the content will affect the original node, create a copy of the content if needed.
func Clone(root *Node) (clone *Node) {
	for root.parent != nil {
		root = root.parent
	}
	return cloneChildren(nil, []*Node{root})[0]
}

func cloneChildren(parent *Node, children []*Node) []*Node {
	clone := make([]*Node, len(children))
	for i, n := range children {
		clone[i] = &Node{
			Content: n.Content,
			parent:  parent,
		}
		if len(n.children) > 0 {
			clone[i].children = cloneChildren(clone[i], n.children)
		}
	}
	return clone
}

func (n *Node) GetRelativePath(relativeRoot *Node) (path Path) {
	path = Path{}
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
	if n == nil {
		return nil
	}
	idx := n.indexInParent()
	if idx < 0 {
		panic("node is not in the tree or root node")
	}
	// performing the replacement with a temporary slice
	// so that the validation logic gets the proper insertion context
	nextSiblings := n.parent.children[idx+1:]
	nodes = slices.Grow(nodes, len(nextSiblings))
	nodes = append(nodes, nextSiblings...)
	n.parent.children = n.parent.children[:idx]
	n.parent.AppendChildren(nodes...)
	n.parent = nil
	return n
}

// NextSibling returns the next sibling of the node or nil if it is the last child.
func (n *Node) NextSibling() *Node {
	if n == nil || n.parent == nil {
		return nil
	}
	children := n.parent.children
	idx := slices.Index(children, n) + 1
	if idx <= 0 || idx >= len(children) {
		return nil
	}
	return children[idx]
}

// PrevSibling returns the previous sibling of the node or nil if it is the first child.
func (n *Node) PrevSibling() *Node {
	if n == nil || n.parent == nil {
		return nil
	}
	children := n.parent.children
	idx := slices.Index(children, n) - 1
	if idx < 0 || idx >= len(children) {
		return nil
	}
	return children[idx]
}

// Parent returns the parent of the node or nil if it is the root node.
func (n *Node) Parent() *Node {
	return n.parent
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

// Cursor represents a position in the AST. It can be used only within a single iteration step.
type Cursor struct {
	path       Path
	isEntering bool
	node       *Node
}

func newCursor(n *Node, startingPath Path) *Cursor {
	return &Cursor{
		path:       startingPath,
		isEntering: true,
		node:       n,
	}
}

func (c *Cursor) dfs(yield func(*Cursor) bool) {
	// TODO: test heavily
	initPath := c.path
	for i, step := range initPath {
		c.path = c.path[:i]
		if !yield(c) {
			return
		}
		c.node = c.node.children[step]
	}
dfs_loop:
	for {
		if !yield(c) {
			return
		}
		if len(c.node.children) != 0 {
			c.node = c.node.children[0]
			c.path = append(c.path, 0)
			continue
		}

		c.isEntering = false

		if !yield(c) {
			return
		}

		curPath := c.path
		for {
			if len(curPath) == len(initPath) {
				// no more nodes, exit
				break dfs_loop
			}
			parent := c.node.parent
			curPath[len(curPath)-1]++
			children := parent.children
			if curPath[len(curPath)-1] < len(children) && curPath[len(curPath)-1] > 0 {
				c.node = children[curPath[len(curPath)-1]]
				c.isEntering = true
				break
			}

			// done with parent, leave it
			c.node = parent
			curPath = curPath[:len(curPath)-1]
			c.path = curPath
			if !yield(c) {
				return
			}
		}
	}
	for i := len(c.path) - 1; i >= 0; i-- {
		c.path = c.path[:i]
		c.node = c.node.parent
		if !yield(c) {
			return
		}
	}
}

// Walk traverses the AST.
// Cursor should not leave the loop body, tree should not be modified during the iteration.
func (n *Node) Walk() iter.Seq[*Cursor] {
	return newCursor(n, Path{}).dfs
}

// Walk traverses the AST according to the given path, and then walks all of the children.
func (n *Node) WalkFromPath(path Path) iter.Seq[*Cursor] {
	return newCursor(n, path).dfs
}

// Path gets the path to the current node.
func (c *Cursor) Path() Path {
	return c.path.Clone()
}

// PathLen gets the length of the path to the current node.
func (c *Cursor) PathLen() int {
	return len(c.path)
}

// Node gets the current node.
func (c *Cursor) Node() *Node {
	return c.node
}

// Content gets the content of the current node.
func (c *Cursor) Content() NodeContent {
	return c.node.Content
}

// IsEntering returns true if the cursor is entering the node.
func (c *Cursor) IsEntering() bool {
	return c.isEntering
}

// IsLeaving returns true if the cursor is leaving the node.
func (c *Cursor) IsLeaving() bool {
	return !c.isEntering
}

// WalkContent traverses the AST in depth-first order, returning only nodes with the given content type.
// Cursor should not leave the loop body, tree should not be modified during the iteration.
func WalkContent[T NodeContent](n *Node) iter.Seq2[*Cursor, T] {
	return func(yield func(*Cursor, T) bool) {
		for c := range n.Walk() {
			if content, ok := c.Content().(T); ok && !yield(c, content) {
				return
			}
		}
	}
}
