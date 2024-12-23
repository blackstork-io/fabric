package ast2md

import (
	"bytes"

	"github.com/blackstork-io/fabric/plugin/ast/nodes"
)

var nestedTablesError = []byte("nested tables are not supported")

func (r *Renderer) renderTable(c *nodes.Table, rows []*nodes.Node) (lines linesSlice) {
	if len(rows) == 0 {
		return
	}

	table := make([][]linesSlice, 0, len(rows))
	columnCount := 0
	for _, cellRowN := range rows {
		cellRow := cellRowN.GetChildren()
		row := make([]linesSlice, 0, len(cellRow))
		columnCount = max(columnCount, len(cellRow))
		for _, cell := range cellRow {
			lns := r.Render(cell)
			lns.JoinLines(space)
			row = append(row, lns)
		}
		table = append(table, row)
	}
	if columnCount == 0 {
		return
	}

	if r.inTable > 0 {
		lines.Append(nestedTablesError)
		return
	}
	r.inTable++
	defer func() {
		r.inTable--
	}()

	colWidths := make([]int, columnCount)
	for i := range colWidths {
		colWidths[i] = 1
	}
	maxColWidth := 1
	for _, row := range table {
		for i, cell := range row {
			colWidths[i] = max(colWidths[i], cell.MaxLineLength())
			maxColWidth = max(maxColWidth, colWidths[i])
		}
	}

	alignments := c.Alignments
	if len(alignments) < columnCount {
		alignments = append(alignments, make([]nodes.Alignment, columnCount-len(alignments))...)
	}

	// rendering

	padding := bytes.Repeat(space, maxColWidth)
	lines = renderRow(table[0], colWidths, alignments, padding, false)
	// header separator line
	lines.Extend(
		renderRow(
			nil, colWidths, alignments,
			bytes.Repeat(dash, maxColWidth),
			true,
		),
	)
	for _, row := range table[1:] {
		lines.Extend(renderRow(row, colWidths, alignments, padding, false))
	}
	lines.SetMinMarginBlock(2)
	return
}

func renderRow(row []linesSlice, colWidths []int, alignments []nodes.Alignment, padding []byte, header bool) (line linesSlice) {
	line.Append(pipe)
	for i, width := range colWidths {
		// calculating padding
		contentLen := 0
		if i < len(row) {
			contentLen = row[i].MaxLineLength()
		}
		lAlign := space
		rAlign := space
		rPad := width - contentLen
		lPad := 0
		if header {
			switch alignments[i] {
			case nodes.AlignmentLeft:
				lAlign = colon
			case nodes.AlignmentCenter:
				lAlign = colon
				rAlign = colon
			case nodes.AlignmentRight:
				rAlign = colon
			}
		} else {
			switch alignments[i] {
			case nodes.AlignmentLeft:
			case nodes.AlignmentCenter:
				lPad = rPad / 2
				rPad -= lPad
			case nodes.AlignmentRight:
				lPad = rPad
				rPad = 0
			}
		}

		line.Append(lAlign, padding[:lPad])
		if i < len(row) {
			line.Extend(row[i])
		}
		line.Append(padding[:rPad], rAlign, pipe)
	}
	line.SetMarginBottom(1)
	return
}
