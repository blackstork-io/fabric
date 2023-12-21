package main

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

func Surround(surround string, elems ...string) []string {
	res := make([]string, len(elems))
	for i := range res {
		res[i] = fmt.Sprintf("%s%s%s", surround, elems[i], surround)
	}
	return res
}

func Join(sep string, elems ...string) string {
	return strings.Join(elems, sep)
}

func JoinSurround(sep, surround string, elems ...string) string {
	if len(elems) == 0 {
		return ""
	}
	var b strings.Builder
	resLen := len(sep) * (len(elems) - 1)
	resLen += len(surround) * 2 * len(elems)
	for _, e := range elems {
		resLen += len(e)
	}
	b.Grow(resLen)

	b.WriteString(surround)
	b.WriteString(elems[0])
	b.WriteString(surround)
	for _, e := range elems[1:] {
		b.WriteString(sep)
		b.WriteString(surround)
		b.WriteString(e)
		b.WriteString(surround)
	}
	return b.String()
}

// Proper unicode aware capitalization function. If something is wrong – just returns string as is.
func CapitalizeFirstLetter(s string) string {
	r, offset := utf8.DecodeRuneInString(s)
	if r == utf8.RuneError {
		// do nothing, since this is just a cosmetic function
		return s
	}
	upperR := unicode.ToUpper(r)
	if r == upperR {
		// upper and lower case letters are identical, do not realloc
		return s
	}
	capitalRuneLen := utf8.RuneLen(upperR)
	if capitalRuneLen == -1 {
		return s
	}

	var b strings.Builder

	b.Grow(capitalRuneLen + len(s) - offset)
	b.WriteRune(upperR)
	b.WriteString(s[offset:])
	return b.String()
}
