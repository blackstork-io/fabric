package eval

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"sync"

	"github.com/hashicorp/hcl/v2"

	"github.com/blackstork-io/fabric/parser/definitions"
	"github.com/blackstork-io/fabric/pkg/diagnostics"
	"github.com/blackstork-io/fabric/plugin"
	"github.com/blackstork-io/fabric/plugin/plugindata"
)

type asyncDataEvaluator struct {
	ctx    context.Context
	blocks []*PluginDataAction
	logger *slog.Logger
}

func makeAsyncDataEvaluator(ctx context.Context, doc *Document, logger *slog.Logger) *asyncDataEvaluator {
	return &asyncDataEvaluator{
		ctx:    ctx,
		blocks: doc.DataBlocks,
		logger: logger,
	}
}

func makeAsyncDataEvaluatorWithPath(ctx context.Context, doc *Document, path []string, logger *slog.Logger) *asyncDataEvaluator {

	var dataSourceName string
	if len(path) > 0 {
		dataSourceName = path[0]
	}

	var blockName string
	if len(path) > 1 {
		blockName = path[1]
	}

	matchingBlocks := []*PluginDataAction{}

	for i := range doc.DataBlocks {
		block := doc.DataBlocks[i]

		if dataSourceName != "" && block.PluginName != dataSourceName {
			continue
		}

		if blockName != "" && block.BlockName != blockName {
			continue
		}

		matchingBlocks = append(matchingBlocks, block)
	}

	return &asyncDataEvaluator{
		ctx:    ctx,
		blocks: matchingBlocks,
		logger: logger,
	}
}


type asyncDataEvalResult struct {
	pluginName string
	blockName  string
	data       plugindata.Data
	diags      diagnostics.Diag
}

func (doc *asyncDataEvaluator) Execute() (plugindata.Data, diagnostics.Diag) {
	doc.logger.DebugContext(doc.ctx, "Fetching data for the document template")

	resultch := make(chan *asyncDataEvalResult, len(doc.blocks))
	for _, block := range doc.blocks {
		go func(block *PluginDataAction, resultch chan<- *asyncDataEvalResult) {
			doc.logger.DebugContext(doc.ctx, "Fetching data for block", "plugin", block.PluginName, "block", block.BlockName)
			data, diags := block.FetchData(doc.ctx)
			resultch <- &asyncDataEvalResult{
				pluginName: block.PluginName,
				blockName:  block.BlockName,
				data:       data,
				diags:      diags,
			}
		}(block, resultch)
	}

	result := make(plugindata.Map)
	diags := diagnostics.Diag{}

	for i := 0; i < len(doc.blocks); i++ {
		res := <-resultch
		for _, diag := range res.diags {
			diags.Append(diag)
		}

		if diags.HasErrors() {
			return nil, res.diags
		}

		var dsMap plugindata.Map

		if found, ok := result[res.pluginName]; ok {
			dsMap, ok = found.(plugindata.Map)
			if !ok {
				return nil, diagnostics.Diag{{
					Severity: hcl.DiagError,
					Summary:  "Conflicting data",
					Detail:   fmt.Sprintf("Different type data block with the same name already exists at plugin '%s' and block '%s'", res.pluginName, res.blockName),
				}}
			}
		} else {
			dsMap = make(plugindata.Map)
			result[res.pluginName] = dsMap
		}

		dsMap[res.blockName] = res.data
	}

	return result, diags
}

type asyncContent struct {
	ctx       context.Context
	doc       *plugin.ContentSection
	parent    *plugin.ContentSection
	contentID uint32
	content   *Content
	section   *Section
	wg        *sync.WaitGroup
	dependsOn []*asyncContent
}

func (ac *asyncContent) Wait() {
	ac.wg.Wait()
}

func (ac *asyncContent) Execute(dataCtx plugindata.Map) diagnostics.Diag {
	defer ac.wg.Done()
	for _, dep := range ac.dependsOn {
		dep.Wait()
	}
	dataCtx = dataCtx.Clone()
	if ac.section != nil {
		diags := ac.section.PrepareData(ac.ctx, dataCtx, ac.doc, ac.parent)
		if diags.HasErrors() {
			return diags
		}
	}
	daigs := ac.content.RenderContent(ac.ctx, dataCtx, ac.doc, ac.parent, ac.contentID)
	return daigs
}

func makeAsyncContent(
	ctx context.Context,
	doc *plugin.ContentSection,
	parent *plugin.ContentSection,
	content *Content,
	dependsOn []*asyncContent,
	section *Section,
) *asyncContent {
	tmp := new(plugin.ContentEmpty)
	_ = parent.Add(tmp, nil)

	contentID := tmp.ID()
	wg := new(sync.WaitGroup)
	wg.Add(1)
	return &asyncContent{
		ctx:       ctx,
		doc:       doc,
		parent:    parent,
		contentID: contentID,
		content:   content,
		wg:        wg,
		dependsOn: dependsOn,
		section:   section,
	}
}

type asyncContentEvaluator struct {
	invokeMap map[plugin.InvocationOrder][]*asyncContent
	namedMap  map[string]*asyncContent
	rootNode  *plugin.ContentSection
}

func (ace *asyncContentEvaluator) executeGroup(dataCtx plugindata.Map, invokeOrder plugin.InvocationOrder) diagnostics.Diag {
	list, ok := ace.invokeMap[invokeOrder]
	if !ok || len(list) == 0 {
		return nil
	}

	diagch := make(chan diagnostics.Diag, len(list))

	for _, ac := range ace.invokeMap[invokeOrder] {
		go func(ac *asyncContent, dataCtx plugindata.Map, diagch chan<- diagnostics.Diag) {
			diagch <- ac.Execute(dataCtx)
		}(ac, dataCtx, diagch)
	}

	diags := diagnostics.Diag{}

	for i := 0; i < len(list); i++ {
		diags.Extend(<-diagch)
	}

	return diags
}

func (ace *asyncContentEvaluator) Execute(dataCtx plugindata.Map) (*plugin.ContentSection, diagnostics.Diag) {
	order := []plugin.InvocationOrder{
		plugin.InvocationOrderBegin,
		plugin.InvocationOrderUnspecified,
		plugin.InvocationOrderEnd,
	}

	diags := diagnostics.Diag{}

	for _, o := range order {
		diags.Extend(ace.executeGroup(dataCtx, o))
		if diags.HasErrors() {
			return nil, diags
		}
	}

	ace.rootNode.Compact()

	return ace.rootNode, diags
}

func makeAsyncContentEvaluator(ctx context.Context, content []*Content, dataCtx plugindata.Map) (*asyncContentEvaluator, diagnostics.Diag) {
	namedMap := make(map[string]*asyncContent)
	invokeMap := make(map[plugin.InvocationOrder][]*asyncContent)
	rootNode := plugin.NewSection(0)

	diags := diagnostics.Diag{}
	for _, c := range content {
		diag := assignAsyncContent(ctx, dataCtx, c, rootNode, rootNode, namedMap, invokeMap, nil)
		if diags.Extend(diag) {
			return nil, diag
		}
	}

	return &asyncContentEvaluator{
		invokeMap: invokeMap,
		namedMap:  namedMap,
		rootNode:  rootNode,
	}, diags
}

func assignAsyncContent(
	ctx context.Context,
	dataCtx plugindata.Map,
	c *Content,
	rootNode *plugin.ContentSection,
	parent *plugin.ContentSection,
	namedMap map[string]*asyncContent,
	invokeMap map[plugin.InvocationOrder][]*asyncContent,
	section *Section,
) diagnostics.Diag {
	diags := diagnostics.Diag{}
	switch {
	case c.Plugin != nil:
		var dependsOn []*asyncContent
		if len(c.Plugin.DependsOn) > 0 {
			for _, depName := range c.Plugin.DependsOn {
				dep, ok := namedMap[depName]
				if !ok {
					return diagnostics.Diag{{
						Severity: hcl.DiagError,
						Summary:  "Dependency not found",
						Detail:   fmt.Sprintf("Content block '%s' depends on '%s' but it's not found", c.Plugin.BlockName, depName),
					}}
				}

				dependsOn = append(dependsOn, dep)
			}
		}

		ac := makeAsyncContent(ctx, rootNode, parent, c, dependsOn, section)
		name := strings.Join([]string{
			definitions.BlockKindContent,
			c.Plugin.PluginName,
			c.Plugin.BlockName,
		}, ".")
		namedMap[name] = ac
		invokeMap[c.Plugin.Provider.InvocationOrder] = append(invokeMap[c.Plugin.Provider.InvocationOrder], ac)
	case c.Section != nil:
		tmp := plugin.NewSection(0)
		_ = parent.Add(tmp, nil)
		include, children, diag := c.Section.Unwrap(ctx, dataCtx)
		if diags.Extend(diag) || !include {
			return diag
		}
		for _, child := range children {
			diag := assignAsyncContent(ctx, dataCtx, child, rootNode, tmp, namedMap, invokeMap, c.Section)
			if diags.Extend(diag) {
				return diag
			}
		}
	default:
		return diagnostics.Diag{{
			Severity: hcl.DiagError,
			Summary:  "Invalid content",
			Detail:   "Content block must be either a plugin or a section",
		}}
	}

	return diags
}
