package parser

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"

	"github.com/blackstork-io/fabric/pkg/diagnostics"
	"github.com/blackstork-io/fabric/pkg/parexec"
)

const FabricFileExt = ".fabric"

type FindFilesResult struct {
	Path string
	Err  error
}

func FindFiles(rootDir fs.FS, recursive bool) <-chan FindFilesResult {
	results := make(chan FindFilesResult, 4) //nolint:gomnd
	go func() {
		defer close(results)
		err := fs.WalkDir(rootDir, ".", func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				results <- FindFilesResult{
					Path: path,
					Err:  err,
				}
				return nil //nolint:nilerr
			}
			if d.IsDir() {
				if !recursive && path != "." {
					return fs.SkipDir
				}
				return nil
			}
			if strings.EqualFold(filepath.Ext(path), FabricFileExt) {
				results <- FindFilesResult{
					Path: path,
				}
			}
			return nil
		})
		if err != nil {
			results <- FindFilesResult{
				Path: "",
				Err:  err,
			}
		}
	}()
	return results
}

type parseResult struct {
	diagnostics.Diag
	file *hcl.File
	path string
}

func parseHcl(bytes []byte, path string) parseResult {
	file, diags := hclsyntax.ParseConfig(bytes, path, hcl.InitialPos)
	return parseResult{
		Diag: diagnostics.Diag(diags),
		file: file,
		path: path,
	}
}

func readFile(rootDir fs.FS, path string) ([]byte, error) {
	file, err := rootDir.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	bytes, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	return bytes, nil
}

func Go(dir string) (diags diagnostics.Diag) {
	dirFs := os.DirFS(dir)
	files := FindFiles(dirFs, true)

	bodies := []hcl.Body{}
	fileMap := map[string]*hcl.File{}

	pe := parexec.New(
		parexec.DiskIOLimiter,
		func(res parseResult, _ int) (cmd parexec.Command) {
			if res.path != "" {
				for i := range res.Diag {
					res.Diag[i].Detail = fmt.Sprintf(
						"Error while looking at '%s': %s",
						res.path, res.Diag[i].Detail,
					)
				}
			}
			if diags.Extend(res.Diag) {
				return
			}
			bodies = append(bodies, res.file.Body)
			fileMap[res.path] = res.file
			return
		},
	)

	parexec.MapChan(pe, files, func(foundFile FindFilesResult) (res parseResult) {
		res.path = foundFile.Path
		if foundFile.Err != nil {
			res.Append(&hcl.Diagnostic{
				Severity: hcl.DiagWarning,
				Summary:  "Directory traversal error",
				Detail:   foundFile.Err.Error(),
				Extra:    foundFile.Err,
			})
			return
		}
		bytes, err := readFile(dirFs, res.path)
		if err != nil {
			res.Append(&hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "File read error",
				Detail:   err.Error(),
				Extra:    err,
			})
			return
		}
		return parseHcl(bytes, res.path)
	})
	pe.WaitDoneAndLock()

	return
}
