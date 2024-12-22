package ast2md

import (
	"bytes"

	"github.com/blackstork-io/fabric/plugin/ast/nodes"
)

func (r *Renderer) renderCodeBlock(c *nodes.CodeBlock) (lines linesSlice) {
	fenceChar := byte('`')
	if bytes.IndexByte(c.Language, fenceChar) != -1 {
		fenceChar = '~'
	}

	fence := makeShortestFence(c.Code, fenceChar, 3)
	lines.Append(fence)
	if len(c.Language) > 0 {
		lines.Append(c.Language)
	}
	lines.SetMarginBottom(1).
		AppendLines(c.Code).
		SetMinMarginBottom(1).
		Append(fence).
		SetMarginBottom(2)

	return
}

func (r *Renderer) renderCodeSpan(c *nodes.CodeSpan) (lines linesSlice) {
	code := bytes.ReplaceAll(c.Code, newLine, space)
	needsSpaces := (
	// needs to be surrounded by spaces if
	len(code) != 0 && (
	// it starts OR ends with a backtick
	(code[0] == '`' || code[len(code)-1] == '`') ||
		// or if it starts AND ends with a space,
		(code[0] == ' ' && code[len(code)-1] == ' ' &&
			// AND not entirely made of spaces
			len(bytes.Trim(code, " ")) != 0)))

	fence := makeShortestFence(c.Code, '`', 1)

	if r.inTable > 0 {
		code = pipeEscapeRe.ReplaceAll(code, []byte("$1\\$2"))
	}
	if needsSpaces {
		lines.Append(fence, space, code, space, fence)
	} else {
		lines.Append(fence, code, fence)
	}
	return
}

// findRuns returns a map (set) of the lengths of runs of char in src.
func findRuns(src []byte, char byte) map[int]struct{} {
	runs := make(map[int]struct{})
	idx := 0
	for {
		offset := bytes.IndexByte(src[idx:], char)
		if offset == -1 {
			break
		}
		idx += offset

		start := idx
		idx++
		for idx < len(src) && src[idx] == char {
			idx++
		}
		runs[idx-start] = struct{}{}
	}
	return runs
}

// makeShortestFence returns the shortest fence that can contain the code/code block.
func makeShortestFence(code []byte, fenceChar byte, minFenceLen int) (fence []byte) {
	runs := findRuns(code, fenceChar)
	for {
		if _, found := runs[minFenceLen]; !found {
			break
		}
		minFenceLen++
	}
	return bytes.Repeat([]byte{fenceChar}, minFenceLen)
}
