package astv1

import "slices"

type Content interface {
	ExtendNodes(nodes []*Node) []*Node
}

type Contents []Content

func (c Contents) ExtendNodes(nodes []*Node) []*Node {
	for _, content := range c {
		nodes = content.ExtendNodes(nodes)
	}
	return nodes
}

// InlineContent represents any inline node or something convertible to inline nodes.
type InlineContent interface {
	Content
	isInline()
}

// Inlines is a list of InlineContent.
type Inlines []InlineContent

func (n Inlines) isInline() {}
func (n Inlines) ExtendNodes(nodes []*Node) []*Node {
	for _, inlines := range n {
		nodes = inlines.ExtendNodes(nodes)
	}
	return nodes
}

func (n *Node_CodeSpan) isInline() {}
func (n *Node_CodeSpan) ExtendNodes(nodes []*Node) []*Node {
	return append(nodes, &Node{Kind: n})
}

func (n *Node_Emphasis) isInline() {}
func (n *Node_Emphasis) ExtendNodes(nodes []*Node) []*Node {
	return append(nodes, &Node{Kind: n})
}

func (n *Node_LinkOrImage) isInline() {}
func (n *Node_LinkOrImage) ExtendNodes(nodes []*Node) []*Node {
	return append(nodes, &Node{Kind: n})
}

func (n *Node_LinkOrImage) SetDestination(url string) *Node_LinkOrImage {
	n.LinkOrImage.Destination = []byte(url)
	return n
}

func (n *Node_LinkOrImage) SetTitle(title string) *Node_LinkOrImage {
	n.LinkOrImage.Title = []byte(title)
	return n
}

func (n *Node_AutoLink) isInline() {}
func (n *Node_AutoLink) ExtendNodes(nodes []*Node) []*Node {
	return append(nodes, &Node{Kind: n})
}

func (n *Node_RawHtml) isInline() {}
func (n *Node_RawHtml) ExtendNodes(nodes []*Node) []*Node {
	return append(nodes, &Node{Kind: n})
}

func (n *Node_Text) isInline() {}
func (n *Node_Text) ExtendNodes(nodes []*Node) []*Node {
	return append(nodes, &Node{Kind: n})
}

func (n *Node_String_) isInline() {}
func (n *Node_String_) ExtendNodes(nodes []*Node) []*Node {
	return append(nodes, &Node{Kind: n})
}

func (n *Node_Strikethrough) isInline() {}
func (n *Node_Strikethrough) ExtendNodes(nodes []*Node) []*Node {
	return append(nodes, &Node{Kind: n})
}

// checkbox is inline content, but it can be used only in a list item.

// BlockContent represents any block node or something convertible to block nodes.
type BlockContent interface {
	Content
	isBlock()
}
type Blocks []BlockContent

func (b Blocks) isBlock() {}
func (b Blocks) ExtendNodes(nodes []*Node) []*Node {
	for _, inlines := range b {
		nodes = inlines.ExtendNodes(nodes)
	}
	return nodes
}

func (n *Node_Blockquote) isBlock() {}
func (n *Node_Blockquote) ExtendNodes(nodes []*Node) []*Node {
	return append(nodes, &Node{Kind: n})
}

func (n *Node_List) isBlock() {}
func (n *Node_List) ExtendNodes(nodes []*Node) []*Node {
	if n == nil || n.List == nil || len(n.List.GetBase().GetChildren()) == 0 {
		return nodes
	}
	return append(nodes, &Node{
		Kind: n,
	})
}

// SetStart sets the start number of the list (only works on ordered lists, "." and ")" markers).
func (n *Node_List) SetStart(start uint32) *Node_List {
	switch n.List.GetMarker() {
	case '.', ')':
		n.List.Start = start
	}
	return n
}

func prependCheckbox(checked bool, content []BlockContent) (itemChildren []*Node) {
	itemChildren = Blocks.ExtendNodes(content, nil)

	checkbox := &Node{
		Kind: &Node_TaskCheckbox{
			TaskCheckbox: &TaskCheckbox{
				IsChecked: checked,
			},
		},
	}

	if len(itemChildren) > 0 {
		if text, ok := itemChildren[0].GetKind().(*Node_Paragraph); ok {
			text.Paragraph.Base.Children = slices.Insert(text.Paragraph.GetBase().GetChildren(), 0, checkbox)
			return
		}
	}
	// Checkbox was not prepended to existing paragraph, create a new paragraph.
	itemChildren = slices.Insert(itemChildren, 0, &Node{
		Kind: &Node_Paragraph{
			Paragraph: &Paragraph{
				Base: &BaseNode{
					Children: []*Node{checkbox},
				},
			},
		},
	})
	return
}

func (n *Node_List) appendItem(children []*Node) {
	n.List.Base.Children = append(n.List.Base.Children, &Node{
		Kind: &Node_ListItem{
			ListItem: &ListItem{
				Base: &BaseNode{
					Children: children,
				},
			},
		},
	})
}

// AppendItem appends a list item to the list.
func (n *Node_List) AppendItem(content ...BlockContent) *Node_List {
	n.appendItem(Blocks.ExtendNodes(content, nil))
	return n
}

// AppendTaskItem appends a task list item to the list.
func (n *Node_List) AppendTaskItem(checked bool, content ...BlockContent) *Node_List {
	n.appendItem(prependCheckbox(checked, content))
	return n
}

// Thematic breaks
func (n *Node_ThematicBreak) isBlock() {}

func (n *Node_ThematicBreak) ExtendNodes(nodes []*Node) []*Node {
	return append(nodes, &Node{Kind: n})
}

func (n *Node_Heading) isBlock() {}

func (n *Node_Heading) ExtendNodes(nodes []*Node) []*Node {
	return append(nodes, &Node{Kind: n})
}

func (n *Node_CodeBlock) isBlock() {}

func (n *Node_CodeBlock) ExtendNodes(nodes []*Node) []*Node {
	return append(nodes, &Node{Kind: n})
}

func (n *Node_FencedCodeBlock) isBlock() {}

func (n *Node_FencedCodeBlock) ExtendNodes(nodes []*Node) []*Node {
	return append(nodes, &Node{Kind: n})
}

func (n *Node_FencedCodeBlock) SetLanguage(language string) *Node_FencedCodeBlock {
	n.FencedCodeBlock.Info = &Text{
		Segment: []byte(language),
	}
	return n
}

func (n *Node_HtmlBlock) isBlock() {}
func (n *Node_HtmlBlock) ExtendNodes(nodes []*Node) []*Node {
	return append(nodes, &Node{Kind: n})
}

func (n *Node_Paragraph) isBlock() {}
func (n *Node_Paragraph) ExtendNodes(nodes []*Node) []*Node {
	return append(nodes, &Node{Kind: n})
}

func (n *Node_Table) isBlock() {}
func (n *Node_Table) ExtendNodes(nodes []*Node) []*Node {
	if n == nil {
		return nodes
	}
	rows := n.Table.GetBase().GetChildren()
	if len(rows) == 0 {
		return nodes
	}
	headerCellCount := len(rows[0].GetTableRow().GetBase().GetChildren())
	if headerCellCount == 0 {
		// must have at least one table cell
		return nodes
	}
	unspecifiedAlignments := headerCellCount - len(n.Table.Alignments)
	if unspecifiedAlignments <= 0 {
		// Trim the alignments to the number of header cells
		n.Table.Alignments = n.Table.Alignments[:len(n.Table.Alignments)+unspecifiedAlignments]
	} else {
		// Extend the alignments to the number of header cells, filling missing alignments with NONE
		n.Table.Alignments = slices.Grow(n.Table.Alignments, unspecifiedAlignments)
		for range unspecifiedAlignments {
			n.Table.Alignments = append(n.Table.Alignments, CellAlignment_CELL_ALIGNMENT_NONE)
		}
	}
	for rowIdx, row := range rows {
		tr := row.GetTableRow()
		if tr == nil {
			return nodes
		}
		tr.IsHeader = rowIdx == 0
		cells := tr.GetBase().GetChildren()
		tr.Alignments = n.Table.Alignments[:min(len(cells), len(n.Table.Alignments))]
		for i, cell := range cells {
			c := cell.GetTableCell()
			if c == nil {
				return nodes
			}
			if i < len(tr.GetAlignments()) {
				c.Alignment = tr.GetAlignments()[i]
			} else {
				c.Alignment = CellAlignment_CELL_ALIGNMENT_NONE
			}
		}
	}
	return append(nodes, &Node{Kind: n})
}

func (n *Node_Table) SetColumnAlignments(alignments ...CellAlignment) *Node_Table {
	n.Table.Alignments = alignments
	return n
}

func (n *Node_Table) AppendRow(cells ...[]InlineContent) *Node_Table {
	if len(cells) == 0 {
		return n
	}
	cellNodes := make([]*Node, len(cells))
	for i, cellContent := range cells {
		if cellContent == nil {
			continue
		}
		cellNodes[i] = &Node{
			Kind: &Node_TableCell{
				TableCell: &TableCell{
					Base: &BaseNode{
						Children: Inlines.ExtendNodes(cellContent, nil),
					},
				},
			},
		}
	}

	n.Table.Base.Children = append(n.Table.Base.Children, &Node{
		Kind: &Node_TableRow{
			TableRow: &TableRow{
				Base: &BaseNode{
					Children: cellNodes,
				},
				IsHeader: false,
			},
		},
	})
	return n
}

func NewContent(nodes ...BlockContent) *FabricContentNode {
	var children []*Node
	for _, node := range nodes {
		children = node.ExtendNodes(children)
	}
	return &FabricContentNode{
		Root: &BaseNode{
			Children:           children,
			BlankPreviousLines: true,
		},
	}
}
