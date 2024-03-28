package definitions

import (
	"context"
	"sync"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/zclconf/go-cty/cty"

	"github.com/blackstork-io/fabric/parser/blocks/internal/tree"
	"github.com/blackstork-io/fabric/parser/evaluation"
	"github.com/blackstork-io/fabric/pkg/diagnostics"
	"github.com/blackstork-io/fabric/plugin"
)

type Section struct {
	tree.NodeSigil
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

func (s ParsedSection) Render(ctx context.Context, caller evaluation.ContentCaller, dataCtx evaluation.DataContext, result *evaluation.Result) (diags diagnostics.Diag) {
	if s.Meta != nil {
		dataCtx.Set(BlockKindSection, plugin.ConvMapData{
			BlockKindMeta: s.Meta.AsJQData(),
		})
	} else {
		dataCtx.Delete(BlockKindSection)
	}
	if title := s.Title; title != nil {
		diags.Extend(title.Render(ctx, caller, dataCtx.Share(), result))
	}

	for _, content := range s.Content {
		diags.Extend(
			content.Render(ctx, caller, dataCtx.Share(), result),
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

func (s *Section) AsCtyValue() cty.Value {
	return cty.CapsuleVal(s.CtyType(), s)
}

func (*Section) FriendlyName() string {
	return "section"
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
