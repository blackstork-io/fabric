package definitions

import (
	"context"
	"slices"
	"sync"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/zclconf/go-cty/cty"

	"github.com/blackstork-io/fabric/parser/evaluation"
	"github.com/blackstork-io/fabric/pkg/diagnostics"
	"github.com/blackstork-io/fabric/plugin"
)

type Section struct {
	Block       *hclsyntax.Block
	Once        sync.Once
	Parsed      bool
	ParseResult *ParsedSection
}

type ParsedSection struct {
	Meta    *MetaBlock
	Title   *ParsedContent
	Content []Renderable
}

func (s ParsedSection) Name() string {
	return ""
}

func (s ParsedSection) Render(ctx context.Context, caller evaluation.ContentCaller, dataCtx evaluation.DataContext, result *evaluation.Result, contentID uint32) (diags diagnostics.Diag) {
	var meta plugin.ConvertableData
	if s.Meta != nil {
		meta = s.Meta.AsJQData()
	}
	section := new(plugin.ContentSection)
	result.Add(section, &plugin.Location{
		Index: contentID,
	})
	if title := s.Title; title != nil {
		empty := new(plugin.ContentEmpty)
		section.Add(empty, nil)
		dataCtx.Set(BlockKindSection, plugin.ConvMapData{
			BlockKindContent: section,
			BlockKindMeta:    meta,
		})
		diags.Extend(title.Render(ctx, caller, dataCtx.Share(), section, empty.ID()))
	}
	posMap := make(map[int]uint32)
	for i := range s.Content {
		empty := new(plugin.ContentEmpty)
		section.Add(empty, nil)
		posMap[i] = empty.ID()
	}
	execList := make([]int, 0, len(s.Content))
	for i := range s.Content {
		execList = append(execList, i)
	}
	slices.SortStableFunc(execList, func(a, b int) int {
		ao, _ := caller.ContentInvocationOrder(ctx, s.Content[a].Name())
		bo, _ := caller.ContentInvocationOrder(ctx, s.Content[b].Name())
		return ao.Weight() - bo.Weight()
	})

	for _, idx := range execList {
		content := s.Content[idx]
		dataCtx.Set(BlockKindSection, plugin.ConvMapData{
			BlockKindContent: section,
			BlockKindMeta:    meta,
		})
		diags.Extend(
			content.Render(ctx, caller, dataCtx.Share(), section, posMap[idx]),
		)
	}
	return
}

func (s *Section) IsRef() bool {
	return len(s.Block.Labels) > 0 && s.Block.Labels[0] == PluginTypeRef
}

func (s *Section) nameIdx() int {
	if s.IsRef() {
		return 1
	}
	return 0
}

func (s *Section) Name() string {
	nameIdx := s.nameIdx()
	if len(s.Block.Labels) > nameIdx {
		return s.Block.Labels[nameIdx]
	}
	return ""
}

var _ FabricBlock = (*Section)(nil)

func (s *Section) GetHCLBlock() *hcl.Block {
	return s.Block.AsHCLBlock()
}

var ctySectionType = capsuleTypeFor[Section]()

func (*Section) CtyType() cty.Type {
	return ctySectionType
}

func DefineSection(block *hclsyntax.Block, atTopLevel bool) (section *Section, diags diagnostics.Diag) {
	sect := Section{
		Block: block,
	}

	nameRequired := atTopLevel

	labels := "<ref> "
	if nameRequired {
		labels += "block_name"
	} else {
		labels += "<block_name>"
	}

	diags.Append(validateBlockName(block, sect.nameIdx(), nameRequired))
	diags.Append(validateLabelsLength(block, 2, labels))
	if diags.HasErrors() {
		return
	}

	section = &sect
	return
}
