package parser

import (
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"

	"github.com/blackstork-io/fabric/pkg/diagnostics"
	"github.com/blackstork-io/fabric/pkg/parexec"
)

const FabricFileExt = ".fabric"

func FindFiles(rootDir fs.FS, recursive bool) <-chan string {
	paths := make(chan string, 4)
	go func() {
		defer close(paths)
		fs.WalkDir(rootDir, ".", func(path string, d fs.DirEntry, err error) error {
			// TODO: return diags or use better logging
			if err != nil {
				log.Printf("Parse files error: %s; path: %s", err, path)
			}
			if d.IsDir() {
				if !recursive && path != "." {
					return fs.SkipDir
				} else {
					return nil
				}
			}
			if strings.EqualFold(filepath.Ext(path), FabricFileExt) {
				paths <- path
			}
			return nil
		})
	}()
	return paths
}

// func ParseFile(rootDir fs.FS, path string) (diags diagnostics.Diag) {
// 	return nil
// }

type fileParseResult struct {
	file  *hcl.File
	path  string
	diags diagnostics.Diag
}

func parseHcl(bytes []byte, path string) fileParseResult {
	file, diags := hclsyntax.ParseConfig(bytes, path, hcl.InitialPos)
	return fileParseResult{
		file:  file,
		path:  path,
		diags: diagnostics.Diag(diags),
	}
}

func parseFile(rootDir fs.FS, path string) (res fileParseResult) {
	// bytes, err := os.ReadFile(path)
	file, err := rootDir.Open(path)
	if err != nil {
		res.diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "File open error",
			Detail:   fmt.Sprintf("Failed to open file %s: %s", path, err),
			Extra:    err,
		})
		return
	}
	bytes, err := io.ReadAll(file)
	if err != nil {
		res.diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "File read error",
			Detail:   fmt.Sprintf("Failed to read file %s: %s", path, err),
			Extra:    err,
		})
		return
	}
	return parseHcl(bytes, path)
}

func Go(dir string) (diags diagnostics.Diag) {
	dirFs := os.DirFS(dir)
	files := FindFiles(dirFs, true)

	bodies := []hcl.Body{}
	fileMap := map[string]*hcl.File{}

	pe := parexec.New(
		parexec.DiskIOLimiter,
		func(res fileParseResult, _ int) (cmd parexec.Command) {
			if diags.Extend(res.diags) {
				return
			}
			bodies = append(bodies, res.file.Body)
			fileMap[res.path] = res.file
			return
		},
	)

	parexec.MapChan(pe, files, func(path string) fileParseResult {
		return parseFile(dirFs, path)
	})
	// parexec.Map(pe, files, processFile)
	pe.WaitDoneAndLock()
	// if diags.HasErrors() {
	// 	return nil, nil, diags
	// }
	// body = hcl.MergeBodies(bodies)
	return
}
