package definitions

type ParsedDocument struct {
	Meta    *MetaBlock
	Vars    *ParsedVars
	Content []*ParsedContent
	Data    []*ParsedPlugin
	Publish []*ParsedPlugin
}
