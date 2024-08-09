package plugin

import (
	"fmt"
	"slices"
	"sync"

	"github.com/blackstork-io/fabric/plugin/plugindata"
)

type idStore struct {
	id uint32
	sync.Mutex
}

func (i *idStore) next() uint32 {
	i.Lock()
	defer i.Unlock()
	i.id++
	return i.id
}

type LocationEffect int

const (
	LocationEffectUnspecified LocationEffect = iota
	LocationEffectBefore
	LocationEffectAfter
)

type Location struct {
	Index  uint32
	Effect LocationEffect
}

type ContentResult struct {
	Location *Location
	Content  Content
}

type Content interface {
	setID(id uint32)
	setMeta(meta *ContentMeta)
	AsData() plugindata.Data
	ID() uint32
	AsPluginData() plugindata.Data
	Meta() *ContentMeta
}

type ContentEmpty struct {
	id   uint32
	meta *ContentMeta
}

func (n *ContentEmpty) setID(id uint32) {
	n.id = id
}

func (n *ContentEmpty) setMeta(meta *ContentMeta) {
	n.meta = meta
}

func (n *ContentEmpty) AsData() plugindata.Data {
	return plugindata.Map{
		"type": plugindata.String("empty"),
		"id":   plugindata.Number(n.id),
		"meta": n.meta.AsData(),
	}
}

func (n *ContentEmpty) ID() uint32 {
	return n.id
}

func (n *ContentEmpty) Meta() *ContentMeta {
	return n.meta
}

func (n *ContentEmpty) AsPluginData() plugindata.Data {
	return n.AsData()
}

type ContentSection struct {
	*idStore
	id       uint32
	Children []Content
	meta     *ContentMeta
}

func NewSection(contentID uint32) *ContentSection {
	return &ContentSection{
		idStore: &idStore{
			id: contentID,
		},
		id: contentID,
	}
}

// Add content to the content tree.
func (c *ContentSection) Add(content Content, loc *Location) error {
	return addContent(c, content, loc)
}

func addContent(parent *ContentSection, child Content, loc *Location) error {
	if parent.idStore == nil {
		parent.idStore = &idStore{}
	}
	if section, ok := child.(*ContentSection); ok {
		section.idStore = parent.idStore
	}
	if loc == nil {
		child.setID(parent.next())
		parent.Children = append(parent.Children, child)
		return nil
	}
	if loc.Effect != LocationEffectUnspecified {
		child.setID(parent.next())
	} else {
		child.setID(loc.Index)
	}
	foundIdx := slices.IndexFunc(parent.Children, func(c Content) bool {
		return c.ID() == loc.Index
	})
	if foundIdx > -1 {
		switch loc.Effect {
		case LocationEffectBefore:
			parent.Children = append(parent.Children[:foundIdx], append([]Content{child}, parent.Children[foundIdx:]...)...)
		case LocationEffectAfter:
			parent.Children = append(parent.Children[:foundIdx+1], append([]Content{child}, parent.Children[foundIdx+1:]...)...)
		default:
			parent.Children[foundIdx] = child
		}
		return nil
	}
	for _, c := range parent.Children {
		section, ok := c.(*ContentSection)
		if !ok {
			continue
		}
		err := addContent(section, child, loc)
		if err == ErrContentLocationNotFound {
			continue
		} else if err != nil {
			return err
		}
	}
	return ErrContentLocationNotFound
}

func (c *ContentSection) setID(id uint32) {
	c.id = id
}

func (c *ContentSection) setMeta(meta *ContentMeta) {
	c.meta = meta
	for _, child := range c.Children {
		child.setMeta(meta)
	}
}

func (c *ContentSection) ID() uint32 {
	return c.id
}

func (c *ContentSection) Meta() *ContentMeta {
	return c.meta
}

func (c *ContentSection) AsPluginData() plugindata.Data {
	return c.AsData()
}

// Compact removes empty sections from the content tree.
func (c *ContentSection) Compact() {
	c.Children = slices.DeleteFunc(c.Children, func(c Content) bool {
		_, ok := c.(*ContentEmpty)
		return ok
	})
	for _, child := range c.Children {
		if section, ok := child.(*ContentSection); ok {
			section.Compact()
		}
	}
}

// AsData returns the content tree as a map.
func (c *ContentSection) AsData() plugindata.Data {
	if c == nil {
		return nil
	}
	children := make(plugindata.List, len(c.Children))
	for i, child := range c.Children {
		children[i] = child.AsData()
	}
	return plugindata.Map{
		"type":     plugindata.String("section"),
		"id":       plugindata.Number(c.id),
		"children": children,
		"meta":     c.meta.AsData(),
	}
}

type ContentElement struct {
	id       uint32
	Markdown string
	meta     *ContentMeta
}

func (c *ContentElement) ID() uint32 {
	return c.id
}

func (c *ContentElement) setID(id uint32) {
	c.id = id
}

func (c *ContentElement) Meta() *ContentMeta {
	return c.meta
}

func (c *ContentElement) setMeta(meta *ContentMeta) {
	c.meta = meta
}

func (c *ContentElement) AsPluginData() plugindata.Data {
	return c.AsData()
}

func (c *ContentElement) AsData() plugindata.Data {
	if c == nil {
		return nil
	}
	return plugindata.Map{
		"type":     plugindata.String("element"),
		"id":       plugindata.Number(c.id),
		"markdown": plugindata.String(c.Markdown),
		"meta":     c.meta.AsData(),
	}
}

type ContentMeta struct {
	Provider string
	Plugin   string
	Version  string
}

func (meta *ContentMeta) AsData() plugindata.Data {
	if meta == nil {
		return nil
	}
	return plugindata.Map{
		"provider": plugindata.String(meta.Provider),
		"plugin":   plugindata.String(meta.Plugin),
		"version":  plugindata.String(meta.Version),
	}
}

func ParseContentData(data plugindata.Map) (Content, error) {
	if data == nil {
		return nil, nil
	}
	typ, ok := data["type"].(plugindata.String)
	if !ok {
		return nil, fmt.Errorf("missing type")
	}
	switch string(typ) {
	case "section":
		return parseContentSection(data)
	case "element":
		return parseContentElement(data)
	case "empty":
		return parseContentEmpty(data)
	default:
		return nil, fmt.Errorf("unknown type: %s", typ)
	}
}

func parseContentSection(data plugindata.Map) (*ContentSection, error) {
	if data == nil {
		return nil, nil
	}
	section := &ContentSection{}
	children, ok := data["children"].(plugindata.List)
	if !ok {
		return nil, fmt.Errorf("missing children")
	}
	section.Children = make([]Content, len(children))
	var err error
	for i, child := range children {
		section.Children[i], err = ParseContentData(child.(plugindata.Map))
		if err != nil {
			return nil, err
		}
	}
	id, ok := data["id"].(plugindata.Number)
	if ok {
		section.id = uint32(id)
	}
	meta, ok := data["meta"].(plugindata.Map)
	if ok {
		section.meta = ParseContentMeta(meta)
	}
	return section, nil
}

func parseContentElement(data plugindata.Map) (*ContentElement, error) {
	if data == nil {
		return nil, nil
	}
	elem := &ContentElement{}
	markdown, ok := data["markdown"].(plugindata.String)
	if !ok {
		return nil, fmt.Errorf("missing markdown")
	}
	elem.Markdown = string(markdown)
	id, ok := data["id"].(plugindata.Number)
	if ok {
		elem.id = uint32(id)
	}
	meta, ok := data["meta"].(plugindata.Map)
	if ok {
		elem.meta = ParseContentMeta(meta)
	}
	return elem, nil
}

func parseContentEmpty(data plugindata.Map) (*ContentEmpty, error) {
	if data == nil {
		return nil, nil
	}
	empty := &ContentEmpty{}
	id, ok := data["id"].(plugindata.Number)
	if !ok {
		return nil, fmt.Errorf("missing id")
	}
	empty.id = uint32(id)
	meta, ok := data["meta"].(plugindata.Map)
	if ok {
		empty.meta = ParseContentMeta(meta)
	}
	return empty, nil
}

func ParseContentMeta(data plugindata.Data) *ContentMeta {
	if data == nil {
		return nil
	}
	meta := data.(plugindata.Map)
	provider, _ := meta["provider"].(plugindata.String)
	plugin, _ := meta["plugin"].(plugindata.String)
	version, _ := meta["version"].(plugindata.String)
	return &ContentMeta{
		Provider: string(provider),
		Plugin:   string(plugin),
		Version:  string(version),
	}
}
