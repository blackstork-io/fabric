package definitions

import (
	"context"
	"fmt"
	"strings"

	"github.com/gobwas/glob"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/zclconf/go-cty/cty"

	"github.com/blackstork-io/fabric/cmd/fabctx"
	"github.com/blackstork-io/fabric/pkg/diagnostics"
)

const patternName = "expose_env_vars_with_pattern"

var DefaultEnvVarsPattern = glob.MustCompile("FABRIC_*")

type GlobalConfigDefinition struct {
	block *hclsyntax.Block
}

func (g *GlobalConfigDefinition) GetHCLBlock() *hcl.Block {
	return g.block.AsHCLBlock()
}

func (g *GlobalConfigDefinition) Parse(ctx context.Context) (cfg *GlobalConfig, diags diagnostics.Diag) {
	var globalCfg GlobalConfig
	var diag diagnostics.Diag
	var body hcl.Body

	evalCtx := fabctx.GetEvalContext(ctx)

	globalCfg.EnvVarsPattern, body, diag = g.parseEnvVarPattern(evalCtx, g.block.Body)
	diags.Extend(diag)
	diags.Extend(gohcl.DecodeBody(body, evalCtx, &globalCfg))

	if diags.HasErrors() {
		return
	}
	return &globalCfg, diags
}

func (g *GlobalConfigDefinition) parseEnvVarPattern(evalCtx *hcl.EvalContext, body hcl.Body) (pat glob.Glob, rest hcl.Body, diags diagnostics.Diag) {
	attr, found := g.block.Body.Attributes[patternName]
	if !found {
		return DefaultEnvVarsPattern, body, nil
	}
	defer func() {
		if diags.HasErrors() {
			pat = nil
		}
		diags.Refine(diagnostics.DefaultSubject(attr.Expr.Range()))
	}()

	val, rest, diag := hcldec.PartialDecode(
		g.block.Body,
		&hcldec.AttrSpec{
			Name:     patternName,
			Type:     cty.String,
			Required: true,
		},
		evalCtx,
	)
	if diags.Extend(diag) {
		return
	}
	if val.IsNull() {
		return
	}
	strVal := val.AsString()

	trimmedStr := strings.TrimSpace(strVal)
	if trimmedStr != strVal {
		diags.AddWarn(
			fmt.Sprintf("%q contains a whitespace", patternName),
			"Leading and trailing whitespaces are ignored",
		)
	}
	if trimmedStr == "" {
		return
	}
	var err error
	pat, err = glob.Compile(trimmedStr)
	if err != nil {
		diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  fmt.Sprintf("Failed to parse %q", patternName),
			Detail:   err.Error(),
		})
	}
	return
}

func DefineGlobalConfig(block *hclsyntax.Block) (config *GlobalConfigDefinition, diags diagnostics.Diag) {
	return &GlobalConfigDefinition{
		block: block,
	}, nil
}

type GlobalConfig struct {
	CacheDir       string            `hcl:"cache_dir,optional"`
	PluginRegistry *PluginRegistry   `hcl:"plugin_registry,block"`
	PluginVersions map[string]string `hcl:"plugin_versions,optional"`
	EnvVarsPattern glob.Glob
}

type PluginRegistry struct {
	BaseURL   string `hcl:"base_url,optional"`
	MirrorDir string `hcl:"mirror_dir,optional"`
}

func (g *GlobalConfig) Merge(other *GlobalConfig) {
	if other.CacheDir != "" {
		g.CacheDir = other.CacheDir
	}
	if other.PluginRegistry != nil {
		if g.PluginRegistry == nil {
			g.PluginRegistry = other.PluginRegistry
		} else {
			if other.PluginRegistry.BaseURL != "" {
				g.PluginRegistry.BaseURL = other.PluginRegistry.BaseURL
			}
			if other.PluginRegistry.MirrorDir != "" {
				g.PluginRegistry.MirrorDir = other.PluginRegistry.MirrorDir
			}
		}
	}
	if other.EnvVarsPattern != DefaultEnvVarsPattern {
		g.EnvVarsPattern = other.EnvVarsPattern
	}
	g.PluginVersions = other.PluginVersions
}
