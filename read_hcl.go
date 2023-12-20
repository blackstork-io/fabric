package main

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"github.com/blackstork-io/fabric/pkg/diagnostics"
	"github.com/blackstork-io/fabric/pkg/parexec"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

type fileParseResult struct {
	file  *hcl.File
	diags diagnostics.Diagnostics
}

func readFile(path string) (bytes []byte, err error) {
	bytes, err = os.ReadFile(path)
	if err != nil {
		err = fmt.Errorf("failed to read file '%s': %w", path, err)
	}
	return
}

func parseHcl(bytes []byte, filename string) *fileParseResult {
	file, diags := hclsyntax.ParseConfig(bytes, filename, hcl.InitialPos)
	return &fileParseResult{
		file:  file,
		diags: diagnostics.Diagnostics(diags),
	}
}

func processFile(path string) *fileParseResult {
	bytes, err := readFile(path)
	if err != nil {
		diag := diagnostics.FromErr(err, "File read error")
		diag.Subject = &hcl.Range{Filename: path}
		return &fileParseResult{diags: []*hcl.Diagnostic{diag}}
	}
	return parseHcl(bytes, path)
}

func readAndParse(files []string) (body hcl.Body, fileMap map[string]*hcl.File, diags diagnostics.Diagnostics) {
	slices.Sort(files)
	bodies := make([]hcl.Body, len(files))
	fileMap = make(map[string]*hcl.File, len(files))

	pe := parexec.New(
		parexec.NewLimiter(min(len(files), 4)),
		func(res *fileParseResult, idx int) (cmd parexec.Command) {
			if diags.Extend(res.diags) {
				return
			}
			bodies[idx] = res.file.Body
			fileMap[files[idx]] = res.file
			return
		},
	)
	parexec.Map(pe, files, processFile)
	pe.WaitDoneAndLock()
	if diags.HasErrors() {
		return nil, nil, diags
	}
	body = hcl.MergeBodies(bodies)
	return
}

func fromDisk() (body hcl.Body, fileMap map[string]*hcl.File, diags diagnostics.Diagnostics) {
	// TODO: replace with filepath.WalkDir()
	files, err := filepath.Glob(path + "*.fabric")
	if diags.FromErr(err, "Can't find files") {
		return
	}
	if len(files) == 0 {
		diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Failed to load files",
			Detail:   fmt.Sprintf("no *.fabric files found at %s", path),
		})
		return
	}
	return readAndParse(files)
}
