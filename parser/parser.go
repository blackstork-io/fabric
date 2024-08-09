package parser

import (
	"fmt"
	"io"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"

	"github.com/blackstork-io/fabric/parser/definitions"
	"github.com/blackstork-io/fabric/pkg/diagnostics"
	"github.com/blackstork-io/fabric/pkg/parexec"
	"github.com/blackstork-io/fabric/pkg/utils"
)

// FS-level parsing of fabric files

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

func parseHclBytes(bytes []byte, path string) (res fileParseResult) {
	file, diag := hclsyntax.ParseConfig(bytes, path, hcl.InitialPos)
	res.file = file
	res.path = path
	res.Diag = diagnostics.Diag(diag)
	if res.HasErrors() {
		return
	}

	body := utils.ToHclsyntaxBody(res.file.Body)

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

func ParseDir(dir fs.FS) (*DefinedBlocks, map[string]*hcl.File, diagnostics.Diag) {
	blocks := NewDefinedBlocks()
	fileMap := map[string]*hcl.File{}
	var parseDiags diagnostics.Diag

	// Collects parsed results
	parsePE := parexec.New(
		parexec.CPULimiter,
		func(res fileParseResult, _ int) (cmd parexec.Command) {
			parseDiags.Extend(res.Diag)
			parseDiags.Extend(blocks.Merge(res.blocks))
			fileMap[res.path] = res.file
			return
		},
	)

	// Schedules read file to be parsed
	goParseHCL := parexec.GoWithArgs(parsePE, func(bytes []byte, path string) fileParseResult {
		res := parseHclBytes(bytes, path)
		return res
	})

	var readDiags diagnostics.Diag

	readPE := parexec.New(
		parexec.DiskIOLimiter,
		func(diags diagnostics.Diag, _ int) (_ parexec.Command) {
			readDiags.Extend(diags)
			return
		},
	)

	// Reads files in readPE and schedules them to be parsed in parsePE
	goReadFabricFile := parexec.GoWithArg(readPE, func(path string) diagnostics.Diag {
		bytes, diag := readFabricFile(dir, path)
		if !diag.HasErrors() {
			goParseHCL(bytes, path)
		}
		return diag
	})

	// Walks the given dir and schedules files to be read in readPE
	readPE.Go(func() diagnostics.Diag {
		return FindFabricFiles(dir, true, goReadFabricFile)
	})

	// All files have been read
	readPE.WaitDoneAndLock()
	// All files have been parsed
	parsePE.WaitDoneAndLock()

	// prepending read diags, since they logically happen earlier
	diags := append(readDiags, parseDiags...) //nolint:gocritic

	if len(fileMap) == 0 {
		diags.Add(
			"No valid fabric files found",
			fmt.Sprintf("No valid *.fabric files found at '%s'", dir),
		)
	}

	return blocks, fileMap, diags
}

func parseBlockDefinitions(body *hclsyntax.Body) (res *DefinedBlocks, diags diagnostics.Diag) {
	res = NewDefinedBlocks()

	for _, block := range body.Blocks {
		switch block.Type {
		case definitions.BlockKindData, definitions.BlockKindContent, definitions.BlockKindPublish:
			plugin, dgs := definitions.DefinePlugin(block, true)
			if diags.Extend(dgs) {
				continue
			}
			key := plugin.GetKey()
			if key == nil {
				panic("unable to get the key of the top-level block")
			}
			diags.Append(AddIfMissing(res.Plugins, *key, plugin))
		case definitions.BlockKindDocument:
			blk, dgs := definitions.DefineDocument(block)
			if diags.Extend(dgs) {
				continue
			}
			diags.Append(AddIfMissing(res.Documents, blk.Name, blk))
		case definitions.BlockKindSection:
			blk, dgs := definitions.DefineSection(block, true)
			if diags.Extend(dgs) {
				continue
			}
			diags.Append(AddIfMissing(res.Sections, blk.Name(), blk))
		case definitions.BlockKindConfig:
			cfg, dgs := definitions.DefineConfig(block)
			if diags.Extend(dgs) {
				continue
			}
			key := cfg.GetKey()
			if key == nil {
				panic("unable to get the key of the top-level block")
			}
			diags.Append(AddIfMissing(res.Config, *key, cfg))
		case definitions.BlockKindGlobalConfig:
			if res.GlobalConfig != nil {
				origRng := res.GlobalConfig.GetHCLBlock().DefRange()
				diags.Append(&hcl.Diagnostic{
					Severity: hcl.DiagError,
					Summary:  "Global config redefinition",
					Detail: fmt.Sprintf(
						"Global config must be defined at most once. Original definition at %s:%d",
						origRng.Filename, origRng.Start.Line,
					),
					Subject: block.DefRange().Ptr(),
				})
				continue
			}
			cfg, diag := definitions.DefineGlobalConfig(block)
			if diags.Extend(diag) {
				continue
			}
			res.GlobalConfig = cfg
		default:
			diags.Append(definitions.NewNestingDiag(
				"Top level of fabric document",
				block,
				body,
				[]string{
					definitions.BlockKindData,
					definitions.BlockKindContent,
					definitions.BlockKindPublish,
					definitions.BlockKindDocument,
					definitions.BlockKindSection,
					definitions.BlockKindConfig,
					definitions.BlockKindGlobalConfig,
				}))
		}
	}
	return
}
