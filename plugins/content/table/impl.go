package table

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/blackstork-io/fabric/pkg/jsontools"
	"github.com/blackstork-io/fabric/plugins/content"
)

// Actual implementation of the plugin

type Impl struct{}

var _ content.Plugin = (*Impl)(nil)

const PluginName = "content.table"

func (Impl) Execute(attrsRaw, dictRaw any) (resp string, err error) {
	var attrs struct {
		Text    string   `json:"text"`
		Columns []string `json:"columns"`
	}
	var dict any
	err = jsontools.UnmarshalBytes(attrsRaw, &attrs)
	if err != nil {
		return
	}
	err = jsontools.UnmarshalBytes(dictRaw, &dict)
	if err != nil {
		return
	}

	tmpl, err := template.New(PluginName).Parse(attrs.Text)
	if err != nil {
		err = fmt.Errorf("failed to parse the template: %w; template: `%s`", err, attrs.Text)
		return
	}

	var buf bytes.Buffer
	buf.WriteString(PluginName)
	buf.WriteByte(':')

	err = tmpl.Execute(&buf, dict)
	if err != nil {
		err = fmt.Errorf("failed to execute the template: %w; template: `%s`; dict: `%s`", err, attrs.Text, jsontools.Dump(dict))
		return
	}
	buf.WriteByte('.')

	if len(attrs.Columns) == 0 {
		return buf.String(), nil
	}
	buf.WriteString(attrs.Columns[0])
	for _, col := range attrs.Columns[1:] {
		buf.WriteByte(',')
		buf.WriteString(col)
	}
	return buf.String(), nil
}
