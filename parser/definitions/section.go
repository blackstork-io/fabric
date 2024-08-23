package definitions

import (
	"sync"

	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/zclconf/go-cty/cty"

	"github.com/blackstork-io/fabric/pkg/diagnostics"
	"github.com/blackstork-io/fabric/pkg/encapsulator"
)

type Section struct {
	Block       *hclsyntax.Block
	Once        sync.Once
	Parsed      bool
	ParseResult *ParsedSection
}

type ParsedSection struct {
	Source       *Section
	Meta         *MetaBlock
	Title        *ParsedContent
	Content      []*ParsedContent
	Vars         *ParsedVars
	RequiredVars []string
}

func (s ParsedSection) Name() string {
	return ""
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

func (s *Section) GetHCLBlock() *hclsyntax.Block {
	return s.Block
}

var ctySectionType = encapsulator.NewEncoder[Section]("section", nil)

func (*Section) CtyType() cty.Type {
	return ctySectionType.CtyType()
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
