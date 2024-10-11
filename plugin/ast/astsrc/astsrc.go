package astsrc

import (
	"fmt"

	"github.com/yuin/goldmark/text"
)

// ASTSource holds the source of the markdown AST (a read-only byte slice).
type ASTSource []byte

// Append appends bytes to the source and returns corresponding segment.
func (s *ASTSource) Append(data []byte) text.Segment {
	start := len(*s)
	*s = append(*s, data...)
	return text.NewSegment(start, len(*s))
}

// AppendMultiple appends multiple byte slices to the source and returns corresponding segments.
func (s *ASTSource) AppendMultiple(data [][]byte) *text.Segments {
	values := make([]text.Segment, len(data))
	for i, segment := range data {
		values[i] = s.Append(segment)
	}
	res := text.NewSegments()
	res.AppendAll(values)
	return res
}

// AppendString appends a string to the source and returns corresponding segment.
func (s *ASTSource) AppendString(data string) text.Segment {
	return s.Append([]byte(data))
}

// Appendf appends formatted string to the source and returns corresponding segment.
func (s *ASTSource) Appendf(format string, args ...interface{}) text.Segment {
	start := len(*s)
	*s = fmt.Appendf(*s, format, args...)
	return text.NewSegment(start, len(*s))
}

// AsBytes returns the source as a byte slice.
// Returned bytes should be treated as read-only or modified with care,
// ensuring that the offsets are not changed.
func (s ASTSource) AsBytes() []byte {
	return s
}
