package dataspec

import (
	"bufio"
	"slices"
	"strings"
	"unicode"

	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
)

func exampleValueForType(ty cty.Type) cty.Value {
	switch {
	case ty.Equals(cty.NilType):
		// just something to get out of the bad situation
		return cty.StringVal("<unspecified default value>")
	case ty.Equals(cty.Bool):
		return cty.True
	case ty.Equals(cty.String):
		return cty.StringVal("some string")
	case ty.Equals(cty.Number):
		return cty.NumberIntVal(42)
	case ty.IsObjectType():
		v := map[string]cty.Value{}
		for name, typ := range ty.AttributeTypes() {
			v[name] = exampleValueForType(typ)
		}
		return cty.ObjectVal(v)
	case ty.IsListType():
		elem := exampleValueForType(ty.ElementType())
		return cty.ListVal([]cty.Value{
			elem, elem,
		})
	case ty.IsMapType():
		elem := exampleValueForType(ty.ElementType())
		return cty.MapVal(map[string]cty.Value{
			"key1": elem,
			"key2": elem,
		})
	case ty.IsSetType():
		elem := exampleValueForType(ty.ElementType())
		return cty.SetVal([]cty.Value{
			elem,
		})
	case ty.IsTupleType():
		v := []cty.Value{}
		for _, innerTy := range ty.TupleElementTypes(){
			v = append(v, exampleValueForType(innerTy))
		}
		return cty.TupleVal(v)
	default:
		return cty.NullVal(ty)
	}
}

func comment(tokens hclwrite.Tokens, text string) hclwrite.Tokens {
	s := bufio.NewScanner(strings.NewReader(strings.ReplaceAll(text, "\t", "    ")))
	var lines []string
	var commonWhitespace string
	lenNonempty := 0

	for s.Scan() {
		line := s.Text()
		if strings.TrimSpace(line) == "" {
			if lenNonempty != 0 {
				lines = append(lines, commonWhitespace)
			}
			continue
		}

		whitespace := 0
		for whitespace < len(line) && line[whitespace] == ' ' {
			whitespace++
		}
		if lenNonempty == 0 {
			commonWhitespace = strings.Repeat(" ", whitespace)
		} else if whitespace < len(commonWhitespace) {
			commonWhitespace = commonWhitespace[:whitespace]
		}
		lines = append(lines, line)
		lenNonempty = len(lines)
	}
	tokens = slices.Grow(tokens, lenNonempty)

	var sb strings.Builder

	for _, line := range lines[:lenNonempty] {
		sb.WriteString("# ")
		sb.WriteString(strings.TrimRightFunc(line[len(commonWhitespace):], unicode.IsSpace))
		sb.WriteByte('\n')
		tokens = append(tokens, &hclwrite.Token{
			Type:  hclsyntax.TokenComment,
			Bytes: []byte(sb.String()),
		})
		sb.Reset()
	}
	return tokens
}

func appendCommentNewLine(tokens hclwrite.Tokens) hclwrite.Tokens {
	return append(tokens, &hclwrite.Token{
		Type:  hclsyntax.TokenComment,
		Bytes: []byte("#\n"),
	})
}
