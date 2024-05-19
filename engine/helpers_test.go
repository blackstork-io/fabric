package engine

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/assert"

	"github.com/blackstork-io/fabric/pkg/diagnostics"
	"github.com/blackstork-io/fabric/pkg/diagnostics/diagtest"
	"github.com/blackstork-io/fabric/plugin"
	"github.com/blackstork-io/fabric/print/mdprint"
)

func renderTest(t *testing.T, testName string, files []string, docName string, expectedResult []string, diagAsserts diagtest.Asserts) {
	t.Helper()
	t.Run(testName, func(t *testing.T) {
		t.Helper()
		sourceDir := fstest.MapFS{}
		for i, content := range files {
			sourceDir[fmt.Sprintf("file_%d.fabric", i)] = &fstest.MapFile{
				Data: []byte(content),
				Mode: 0o777,
			}
		}
		eng := New()
		defer func() {
			eng.Cleanup()
		}()

		var res string
		ctx := context.Background()
		diags := eng.ParseDir(ctx, sourceDir)
		if !diags.HasErrors() {
			if !diags.Extend(eng.LoadPluginResolver(ctx, false)) && !diags.Extend(eng.LoadPluginRunner(ctx)) {
				content, _, diag := eng.RenderContent(ctx, docName)
				diags.Extend(diag)
				res = mdprint.PrintString(content)
			}
		}

		if len(expectedResult) == 0 {
			// so nil == []string{}
			assert.Empty(t, res)
		} else {
			assert.EqualValues(
				t,
				strings.Join(expectedResult, "\n\n"),
				res,
			)
		}
		diagAsserts.AssertMatch(t, diags, eng.FileMap())
	})
}

func fetchDataTest(t *testing.T, testName string, files []string, target string, expectedResult plugin.MapData, diagAsserts diagtest.Asserts) {
	t.Helper()
	t.Run(testName, func(t *testing.T) {
		t.Helper()
		sourceDir := fstest.MapFS{}
		for i, content := range files {
			sourceDir[fmt.Sprintf("file_%d.fabric", i)] = &fstest.MapFile{
				Data: []byte(content),
				Mode: 0o777,
			}
		}
		eng := New()
		defer func() {
			eng.Cleanup()
		}()
		var res plugin.Data
		ctx := context.Background()
		diags := eng.ParseDir(ctx, sourceDir)
		if !diags.HasErrors() {
			if !diags.Extend(eng.LoadPluginResolver(ctx, false)) && !diags.Extend(eng.LoadPluginRunner(ctx)) {
				var diag diagnostics.Diag
				res, diag = eng.FetchData(ctx, target)
				diags.Extend(diag)
			}
		}

		assert.Equal(t, expectedResult, res)
		diagAsserts.AssertMatch(t, diags, eng.FileMap())
	})
}

func lintTest(t *testing.T, fullLint bool, testName string, files []string, diagAsserts diagtest.Asserts) {
	t.Helper()
	t.Run(testName, func(t *testing.T) {
		t.Helper()
		sourceDir := fstest.MapFS{}
		for i, content := range files {
			sourceDir[fmt.Sprintf("file_%d.fabric", i)] = &fstest.MapFile{
				Data: []byte(content),
				Mode: 0o777,
			}
		}
		ctx := context.Background()

		eng := New()
		defer func() {
			eng.Cleanup()
		}()
		diag := []diagnostics.Diag{
			eng.ParseDir(ctx, sourceDir),
			eng.LoadPluginResolver(ctx, false),
			eng.LoadPluginRunner(ctx),
		}
		for _, d := range diag {
			if d.HasErrors() {
				t.Fatalf("Error: %v", d)
			}
		}
		diags := eng.Lint(ctx, fullLint)
		diagAsserts.AssertMatch(t, diags, eng.FileMap())
	})
}

func fullLintTest(t *testing.T, testName string, files []string, diagAsserts [][]diagtest.Assert) {
	t.Helper()
	lintTest(t, true, testName, files, diagAsserts)
}

func limitedLintTest(t *testing.T, testName string, files []string, diagAsserts [][]diagtest.Assert) {
	t.Helper()
	lintTest(t, false, testName, files, diagAsserts)
}
