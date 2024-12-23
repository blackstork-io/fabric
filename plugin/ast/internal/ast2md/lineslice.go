package ast2md

import (
	"bytes"
	"io"
	"slices"
)

// Collection of fragments. Each fragment is it's own line
type linesSlice []*fragment

// WriteTo implements io.WriterTo.
func (l linesSlice) WriteTo(w io.Writer) (n int64, err error) {
	if len(l) == 0 {
		return
	}
	var m int64
	for _, l := range l {
		m, err = l.WriteTo(w)
		n += m
		if err != nil {
			return
		}
	}
	return
}

// first returns the first fragment in the slice. If the slice is empty, a new fragment is appended
func (l *linesSlice) first() *fragment {
	if len(*l) == 0 {
		*l = append(*l, &fragment{})
	}
	return (*l)[0]
}

// last returns the last fragment in the slice. If the slice is empty, a new fragment is appended
func (l *linesSlice) last() *fragment {
	if len(*l) == 0 {
		*l = append(*l, &fragment{})
	}
	return (*l)[len(*l)-1]
}

// SetMarginBlock sets the margin on the top and bottom of lines
func (l *linesSlice) SetMarginBlock(n int) *linesSlice {
	return l.SetMarginTop(n).SetMarginBottom(n)
}

// SetMinMarginBlock sets the margin on the top and bottom of lines to at least n
func (l *linesSlice) SetMinMarginBlock(n int) *linesSlice {
	return l.SetMinMarginTop(n).SetMinMarginBottom(n)
}

// MarginTop returns the number of margin lines at the top of the block
func (l *linesSlice) MarginTop() (n int) {
	for _, line := range *l {
		if line.isMargin() {
			n += 1
		} else {
			break
		}
	}
	return
}

// MarginBottom returns the number of margin lines at the bottom of the block
// Does not include the last newline of the block
func (l *linesSlice) MarginBottom() (n int) {
	for i := len(*l) - 1; i >= 0; i-- {
		if (*l)[i].isMargin() {
			n += 1
		} else {
			break
		}
	}
	return
}

// SetMarginTop sets the number of margin lines at the top of the block
func (l *linesSlice) SetMarginTop(n int) *linesSlice {
	delta := n - l.MarginTop()
	if delta > 0 {
		*l = slices.Insert(*l, 0, make([]*fragment, delta)...)
		for i := 0; i < delta; i++ {
			(*l)[i] = new(fragment).setMargin(true)
		}
	} else if delta < 0 {
		*l = slices.Delete(*l, 0, -delta)
	}
	return l
}

// SetMinMarginTop sets the number of margin lines at the top of the block to at least n
func (l *linesSlice) SetMinMarginTop(n int) *linesSlice {
	if l.MarginTop() < n {
		l.SetMarginTop(n)
	}
	return l
}

// SetMarginBottom sets the number of margin lines at the bottom of the block
func (l *linesSlice) SetMarginBottom(n int) *linesSlice {
	delta := n - l.MarginBottom()
	if delta > 0 {
		if !l.last().isNewLine() {
			l.last().setNewLine(true)
		}
		delta--
	}
	initLen := len(*l)
	if delta > 0 {
		*l = append(*l, make([]*fragment, delta)...)
		for i := initLen; i < delta+initLen; i++ {
			(*l)[i] = new(fragment).setMargin(true)
		}
	} else if delta < 0 {
		*l = (*l)[:initLen+delta]
	}
	return l
}

// SetMinMarginBottom sets the number of margin lines at the bottom of the block to at least n
func (l *linesSlice) SetMinMarginBottom(n int) *linesSlice {
	if l.MarginBottom() < n {
		l.SetMarginBottom(n)
	}
	return l
}

// Append appends data to the lines. New lines are replaced with spaces
func (l *linesSlice) Append(data ...[]byte) *linesSlice {
	last := l.last()
	if last.isNewLine() {
		last = &fragment{}
		*l = append(*l, last)
	}
	for _, b := range data {
		if len(b) != 0 {
			last.append(bytes.ReplaceAll(b, newLine, space))
		}
	}
	return l
}

// Surround wraps nonempty lines with the given prefix and suffix
func (l *linesSlice) Surround(with []byte) *linesSlice {
	for _, line := range *l {
		if !line.isEmpty() {
			line.surround(with)
		}
	}

	return l
}

// RemoveEmptyLines removes empty lines in the linesSlice
func (l *linesSlice) RemoveEmptyLines() *linesSlice {
	*l = slices.DeleteFunc(*l, func(l *fragment) bool {
		return l.isEmpty() || l.isMargin()
	})
	return l
}

// MaxLineLength returns the length of the longest line in the slice
func (l *linesSlice) MaxLineLength() int {
	maxLen := 0
	for _, line := range *l {
		maxLen = max(maxLen, line.length())
	}
	return maxLen
}

// AppendLines appends data (possibly containing newlines) to the lines.
func (l *linesSlice) AppendLines(dataWithNl []byte) *linesSlice {
	last := l.last()
	for len(dataWithNl) != 0 {
		if last.isNewLine() {
			last = new(fragment)
			*l = append(*l, last)
		}
		nlPos := bytes.IndexByte(dataWithNl, '\n')
		if nlPos != -1 {
			last.append(dataWithNl[:nlPos]).setNewLine(true)
			dataWithNl = dataWithNl[nlPos+1:]
		} else {
			last.append(dataWithNl)
			break
		}
	}
	return l
}

// AppendHardBreak appends a hard break to the last line
func (l *linesSlice) AppendHardBreak() *linesSlice {
	l.last().setHardBreak(true)
	return l
}

// IsEmpty returns true if all lines in the slice are empty
func (l *linesSlice) IsEmpty() bool {
	for _, line := range *l {
		if !line.isEmpty() {
			return false
		}
	}
	return true
}

// PrependPrefixes prepends a prefix to all lines in the slice. First line gets its own prefix.
// Margins on the block level (top and bottom row) don't get prefixed
func (l *linesSlice) PrependPrefixes(firstLinePrefix, prefix []byte) *linesSlice {
	start := l.MarginTop()
	end := len(*l) - l.MarginBottom()
	if start < end {
		(*l)[start].pref = append((*l)[start].pref, firstLinePrefix)
		start++
		for i := start; i < end; i++ {
			(*l)[i].pref = append((*l)[i].pref, prefix)
		}
	} else {
		// completely empty block, add an empty line with firstLinePrefix to preserve
		emptyPrefixLine := new(fragment)
		emptyPrefixLine.pref = append(emptyPrefixLine.pref, firstLinePrefix)
		emptyPrefixLine.setNewLine(true)
		*l = append(*l, emptyPrefixLine)
	}
	return l
}

// PrependPrefix prepends a prefix to all lines in the slice.
// Margins on the block level (top and bottom row) don't get prefixed
func (l *linesSlice) PrependPrefix(prefix []byte) *linesSlice {
	return l.PrependPrefixes(prefix, prefix)
}

// TrimEmptyLines removes empty lines (and margins) from the beginning and end of the slice
func (l *linesSlice) TrimEmptyLines() *linesSlice {
	lines := *l
	var startIdx, endIdx int
	for startIdx = 0; startIdx < len(lines); startIdx++ {
		if !lines[startIdx].isEmpty() {
			break
		}
	}
	for endIdx = len(lines) - 1; endIdx > startIdx; endIdx-- {
		if !lines[endIdx].isEmpty() {
			break
		}
	}
	endIdx++
	*l = lines[startIdx:endIdx]
	l.last().setNewLine(false)
	return l
}

// ClearPrefixes removes all prefixes from the lines
func (l *linesSlice) ClearPrefixes() *linesSlice {
	for _, line := range *l {
		line.pref = nil
	}
	return l
}

// Extend appends the other lines to the end of the current lines, collapsing margins
func (l *linesSlice) Extend(other linesSlice) *linesSlice {
	if len(other) == 0 {
		return l
	}
	lines := *l
	if len(lines) == 0 {
		*l = other
		return l
	}
	last := lines.last()
	bot := lines.MarginBottom()
	top := other.MarginTop()
	if top != 0 && bot == 0 && !last.isNewLine() {
		last.setNewLine(true)
	}
	delta := top - bot - 1
	if delta > 0 {
		other = other[top-delta:]
	} else {
		other = other[top:]
	}

	for _, line := range other {
		switch {
		case len(line.pref) != 0:
			last.setNewLine(true)
			fallthrough
		case last.isNewLine():
			lines = append(lines, line)
			last = line
		default:
			last.content = append(last.content, line.content...)
			last.lineBreak = line.lineBreak
		}
	}
	*l = lines
	return l
}

// JoinLines collapses lines into a single line with the given separator
// Prefixes (except the first line prefix) are discarded
func (l *linesSlice) JoinLines(with []byte) *linesSlice {
	first := l.first()
	for _, line := range (*l)[1:] {
		if line.isEmpty() {
			continue
		}
		if len(with) != 0 {
			first.append(with)
		}
		first.append(line.content...)
	}
	clear((*l)[1:])
	*l = (*l)[:1]
	return l
}
