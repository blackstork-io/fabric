package definitions

type ParsedDocument struct {
	Meta    *MetaBlock
	Content []*ParsedContent
	Data    []*ParsedPlugin
	Publish []*ParsedPlugin
}
