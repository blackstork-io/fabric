package ast2md

import (
	"io"
	"slices"
	"unicode/utf8"
)

const (
	// regularNewLine is rendered as a single newline character
	regularNewLine = byte(1 << iota)
	// hardBreakNewLine is rendered as a hard break (literal \ before newline)
	hardBreakNewLine
	// marginNewLine is rendered as a regular newline, but can be collapsed with other margins
	marginNewLine
)

// A fragment of a line or a whole line
type fragment struct {
	pref      [][]byte
	content   [][]byte
	lineBreak byte
}

var _ io.WriterTo = &fragment{}

func (l *fragment) isHardBreak() bool {
	return l.lineBreak&hardBreakNewLine == hardBreakNewLine
}

func (l *fragment) setHardBreak(val bool) *fragment {
	if val {
		l.setMargin(false)
		l.lineBreak |= hardBreakNewLine
	} else {
		l.lineBreak &= ^hardBreakNewLine
	}
	return l
}

func (l *fragment) setMargin(val bool) *fragment {
	if val {
		l.lineBreak |= marginNewLine
	} else {
		l.lineBreak &= ^marginNewLine
	}
	return l
}

func (l *fragment) isMargin() bool {
	return l.lineBreak&marginNewLine == marginNewLine
}

func (l *fragment) isNewLine() bool {
	return l.lineBreak != 0
}

func (l *fragment) setNewLine(val bool) *fragment {
	if val {
		l.setMargin(false)
		l.lineBreak |= regularNewLine
	} else {
		l.lineBreak = 0
	}
	return l
}

func (l *fragment) writePrefix(w io.Writer) (n int64, err error) {
	var m int
	for i := len(l.pref) - 1; i >= 0; i-- {
		m, err = w.Write(l.pref[i])
		n += int64(m)
		if err != nil {
			return
		}
	}
	return
}

func (l *fragment) WriteTo(w io.Writer) (n int64, err error) {
	var m int
	n, err = l.writePrefix(w)
	if err != nil {
		return
	}
	for _, c := range l.content {
		m, err = w.Write(c)
		n += int64(m)
		if err != nil {
			return
		}
	}
	if l.isHardBreak() {
		m, err = w.Write(hardBreak)
		n += int64(m)
		if err != nil {
			return
		}
	} else if l.isNewLine() {
		m, err = w.Write(newLine)
		n += int64(m)
		if err != nil {
			return
		}
	}
	return
}

func (l *fragment) append(b ...[]byte) *fragment {
	l.content = append(l.content, b...)
	return l
}

func (l *fragment) length() int {
	res := 0
	for _, c := range l.content {
		res += utf8.RuneCount(c)
	}
	return res
}

func (l *fragment) surround(with []byte) *fragment {
	l.content = slices.Insert(l.content, 0, with)
	l.content = append(l.content, with)
	return l
}

func (l *fragment) isEmpty() bool {
	if l.isMargin() {
		return true
	}
	for _, c := range l.content {
		if len(c) != 0 {
			return false
		}
	}
	return true
}
