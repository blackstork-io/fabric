package definitions

type ParsedDocument struct {
	Source       *Document
	Meta         *MetaBlock
	Vars         *ParsedVars
	RequiredVars []string
	Content      []*ParsedContent
	Data         []*ParsedPlugin
	Publish      []*ParsedPlugin
}
