package astv1_test

// import (
// 	"bytes"
// 	"testing"

// 	"github.com/blackstork-io/goldmark-markdown/pkg/mdexamples"
// 	"github.com/yuin/goldmark"
// 	"github.com/yuin/goldmark/extension"
// 	"github.com/yuin/goldmark/text"

// 	astv1 "github.com/blackstork-io/fabric/plugin/ast/v1"
// )

// func roundtrip(t *testing.T, source []byte) {
// 	t.Helper()
// 	md := goldmark.New(
// 		goldmark.WithExtensions(
// 			extension.Table,
// 			extension.Strikethrough,
// 			extension.TaskList,
// 		),
// 	)

// 	tree := md.Parser().Parse(text.NewReader(source))

// 	// roundtrip
// 	encTree, err := astv1.Encode(tree, source)
// 	if err != nil {
// 		t.Fatalf("encode: %v", err)
// 	}
// 	decTree, decSrc, err := astv1.Decode(encTree)
// 	if err != nil {
// 		t.Fatalf("decode: %v", err)
// 	}

// 	// assume that trees are identical if their html render is identical
// 	// can't use equality checks because the segments hold no reference to *which*
// 	// source they are referring to
// 	var bufOrig, bufRoundtrip bytes.Buffer

// 	err = md.Renderer().Render(&bufOrig, source, tree)
// 	if err != nil {
// 		t.Fatalf("render original: %v", err)
// 	}

// 	err = md.Renderer().Render(&bufRoundtrip, decSrc, decTree)
// 	if err != nil {
// 		t.Fatalf("render decoded: %v", err)
// 	}

// 	if !bytes.Equal(bufOrig.Bytes(), bufRoundtrip.Bytes()) {
// 		t.Errorf("rendered html differs")
// 		t.Logf("---------- original -----------\n%s\n\n", bufOrig.String())
// 		t.Logf("---------- roundtrip ----------\n%s\n\n", bufRoundtrip.String())
// 	}
// }

// func TestSpecExamplesRoundtrip(t *testing.T) {
// 	for _, exFile := range mdexamples.ReadAllSpecExamples() {
// 		for _, ex := range exFile.Examples {
// 			t.Run(exFile.Name+":"+ex.Link, func(t *testing.T) {
// 				t.Parallel()
// 				// Skipped tests are ok to use, we're only interested in the roundtrip
// 				roundtrip(t, ex.Markdown)
// 			})
// 		}
// 	}
// }

// func TestDocumentsRoundtrip(t *testing.T) {
// 	for _, ex := range mdexamples.ReadAllDocumentExamples() {
// 		t.Run(ex.Name, func(t *testing.T) {
// 			t.Parallel()
// 			roundtrip(t, ex.Data)
// 		})
// 	}
// }

// func TestFuzzCase(t *testing.T) {
// 	t.Skip("Expected to fail: bugs in goldmark")
// 	roundtrip(t, []byte("* 0\n-|\n\t0"))
// 	roundtrip(t, []byte(">  ```\n>\t0"))
// }

// func FuzzEncoder(f *testing.F) {
// 	// for _, exFile := range test.ReadAllSpecExamples() {
// 	// 	for _, ex := range exFile.Examples {
// 	// 		f.Add(ex.Markdown)
// 	// 	}
// 	// }
// 	for _, ex := range mdexamples.ReadAllDocumentExamples() {
// 		f.Add(ex.Data)
// 	}
// 	f.Fuzz(func(t *testing.T, data []byte) {
// 		roundtrip(t, data)
// 	})
// }
