package htmlprint

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"strings"

	"github.com/pelletier/go-toml/v2"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"gopkg.in/yaml.v3"

	"github.com/blackstork-io/fabric/plugin"
	"github.com/blackstork-io/fabric/printer/mdprint"
)

//go:embed document.gotempl
var documentTemplStr string

var templ = template.Must(template.New("document").Parse(documentTemplStr))

const (
	fmTitleKey       = "title"
	fmDescriptionKey = "description"
	fmCSSCodeKey     = "css_code"
	fmJSCodeKey      = "js_code"
	fmCSSSourcesKey  = "css_sources"
	fmJSSourcesKey   = "js_sources"
)

type Data struct {
	Title       string
	Description string
	Content     template.HTML
	CSS         template.CSS
	JS          template.JS
	JSSources   []template.URL
	CSSSources  []template.URL
}

type Printer struct {
	md *mdprint.Printer
}

func New() *Printer {
	return &Printer{
		md: mdprint.New(),
	}
}

func (p *Printer) Print(w io.Writer, el plugin.Content) error {
	data := Data{
		Title: "Untitled",
	}
	if title, ok := p.firstTitle(el); ok {
		data.Title = title
	}
	err := p.evalFrontmatter(&data, el)
	if err != nil {
		return err
	}

	buf := bytes.NewBuffer(nil)
	if err := p.md.Print(buf, el); err != nil {
		return err
	}
	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			html.WithXHTML(),
		),
	)
	buff := bytes.NewBuffer(nil)
	if err := md.Convert(buf.Bytes(), buff); err != nil {
		return err
	}
	data.Content = template.HTML(buff.String()) //nolint: gosec
	return templ.Execute(w, data)
}

func (p *Printer) evalFrontmatter(data *Data, el plugin.Content) error {
	fm, ok := p.extractFrontmatter(el)
	if ok {
		parsed, err := p.parseFrontmatter(fm)
		if err != nil {
			return err
		}
		if title, ok := parsed[fmTitleKey].(string); ok {
			data.Title = title
		}
		if attr, ok := parsed[fmDescriptionKey].(string); ok {
			data.Description = attr
		}
		if css, ok := parsed[fmCSSCodeKey].(string); ok {
			data.CSS = template.CSS(css)
		}
		if js, ok := parsed[fmJSCodeKey].(string); ok {
			data.JS = template.JS(js) //nolint: gosec
		}
		if sources, ok := parsed[fmCSSSourcesKey].([]any); ok {
			for _, source := range sources {
				if source, ok := source.(string); ok {
					data.CSSSources = append(data.CSSSources, template.URL(source)) //nolint: gosec
				}
			}
		}
		if sources, ok := parsed[fmJSSourcesKey].([]any); ok {
			for _, source := range sources {
				if source, ok := source.(string); ok {
					data.JSSources = append(data.JSSources, template.URL(source)) //nolint: gosec
				}
			}
		}
	}
	return nil
}

func (p *Printer) firstTitle(el plugin.Content) (string, bool) {
	switch el := el.(type) {
	case *plugin.ContentSection:
		for _, c := range el.Children {
			if title, ok := p.firstTitle(c); ok {
				return title, true
			}
		}
	case *plugin.ContentElement:
		meta := el.Meta()
		if meta.Plugin == "blackstork/builtin" && meta.Provider == "title" {
			return strings.TrimSpace(
				strings.TrimPrefix(el.Markdown, "#"),
			), true
		}
	}
	return "", false
}

func (p *Printer) extractFrontmatter(el plugin.Content) (*plugin.ContentElement, bool) {
	section, ok := el.(*plugin.ContentSection)
	if !ok {
		return nil, false
	}
	for i, c := range section.Children {
		el, ok := c.(*plugin.ContentElement)
		if !ok {
			continue
		}
		meta := c.Meta()
		if meta != nil && meta.Plugin == "blackstork/builtin" && meta.Provider == "frontmatter" {
			section.Children = append(section.Children[:i], section.Children[i+1:]...)
			return el, true
		}
	}
	return nil, false
}

func (p *Printer) parseFrontmatter(fm *plugin.ContentElement) (result map[string]any, err error) {
	str := fm.Markdown
	switch {
	case strings.HasPrefix(str, "{"):
		err = json.Unmarshal([]byte(str), &result)
	case strings.HasPrefix(str, "---"):
		str = strings.Trim(str, "-")
		err = yaml.Unmarshal([]byte(str), &result)
	case strings.HasPrefix(str, "+++"):
		str = strings.Trim(str, "+")
		err = toml.Unmarshal([]byte(str), &result)
	default:
		err = fmt.Errorf("invalid frontmatter format")
	}
	return
}
