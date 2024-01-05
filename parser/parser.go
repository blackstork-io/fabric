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

// Calls fn with paths to every *.fabric files and collects errors into the returned diags.
func FindFabricFiles(rootDir fs.FS, recursive bool, fn func(path string)) (diags diagnostics.Diag) {
	err := fs.WalkDir(rootDir, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			diags.Append(&hcl.Diagnostic{
				Severity: hcl.DiagWarning,
				Summary:  "Directory traversal error",
				Detail: fmt.Sprintf(
					"Error while looking at '%s': %s",
					path, err,
				),
				Extra: err,
			})
			return nil
		}
		if d.IsDir() {
			if !recursive && path != "." {
				return fs.SkipDir
			}
			return nil
		}
		if strings.EqualFold(filepath.Ext(path), FabricFileExt) {
			fn(path)
		}
		return nil
	})
	if err != nil {
		diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagWarning,
			Summary:  "fs.WalkDir error",
			Detail:   err.Error(),
			Extra:    err,
		})
	}
	return
}

type fileParseResult struct {
	diagnostics.Diag
	file   *hcl.File
	path   string
	blocks *DefinedBlocks
}

const (
	BlockKindDocument = "document"
	BlockKindConfig   = "config"
	BlockKindContent  = "content"
	BlockKindData     = "data"
	BlockKindRef      = "ref"
	BlockKindSection  = "section"
)

func parseHclBytes(bytes []byte, path string) (res fileParseResult) {
	file, diag := hclsyntax.ParseConfig(bytes, path, hcl.InitialPos)
	res.file = file
	res.path = path
	res.Diag = diagnostics.Diag(diag)
	if res.HasErrors() {
		return
	}

	body, ok := res.file.Body.(*hclsyntax.Body)
	if !ok {
		res.Add("Failed to parse", fmt.Sprintf("Can't inspect body of the file '%s'", path))
		return
	}

	blocks, diags := parseBlockDefinitions(body)
	res.Extend(diags)
	res.blocks = blocks

	return
}

func readFile(rootDir fs.FS, path string) ([]byte, error) {
	file, err := rootDir.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()
	bytes, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	return bytes, nil
}

func readFabricFile(dirFs fs.FS, path string) ([]byte, diagnostics.Diag) {
	bytes, err := readFile(dirFs, path)
	if err != nil {
		return nil, diagnostics.FromHcl(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "File read error",
			Detail: fmt.Sprintf(
				"Error while looking at '%s': %s",
				path, err,
			),
			Extra: err,
		})
	}
	return bytes, nil
}

type DirParseResult struct {
	diagnostics.Diag
	Blocks  *DefinedBlocks
	FileMap map[string]*hcl.File
}

func ParseDir(dir string) DirParseResult {
	result := DirParseResult{
		Blocks:  NewDefinedBlocks(),
		FileMap: map[string]*hcl.File{},
	}

	// Collects parsed results
	parsePE := parexec.New(
		parexec.CPULimiter,
		func(res fileParseResult, _ int) (cmd parexec.Command) {
			result.Extend(res.Diag)
			result.Extend(result.Blocks.Merge(res.blocks))
			result.FileMap[res.path] = res.file
			return
		},
	)

	// Schedules read file to be parsed
	goParseHCL := parexec.GoWithArgs(parsePE, func(bytes []byte, path string) fileParseResult {
		res := parseHclBytes(bytes, path)
		return res
	})

	dirFs := os.DirFS(dir)
	var readDiags diagnostics.Diag

	readPE := parexec.New(
		parexec.DiskIOLimiter,
		func(diags diagnostics.Diag, _ int) (_ parexec.Command) {
			readDiags.Extend(diags)
			return
		},
	)

	// Reads files in readPE and shedules them to be parsed in parsePE
	goReadFabricFile := parexec.GoWithArg(readPE, func(path string) diagnostics.Diag {
		bytes, diag := readFabricFile(dirFs, path)
		if !diag.HasErrors() {
			goParseHCL(bytes, path)
		}
		return diag
	})

	// Walks the given dir and schedules files to be read in readPE
	readPE.Go(func() diagnostics.Diag {
		return FindFabricFiles(dirFs, true, goReadFabricFile)
	})

	// All files have been read
	readPE.WaitDoneAndLock()
	// All files have been parsed
	parsePE.WaitDoneAndLock()

	// prepending read diags, since they logically happen earlier
	result.Diag = append(readDiags, result.Diag...)
	return result
}
