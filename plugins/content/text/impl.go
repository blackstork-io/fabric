package text

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

const PLUGIN_NAME = "content.text"

func (Impl) Execute(attrsRaw, dictRaw any) (resp string, err error) {
	var attrs struct {
		Text string `json:"text"`
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

	tmpl, err := template.New(PLUGIN_NAME).Parse(attrs.Text)
	if err != nil {
		err = fmt.Errorf("failed to parse the template: %w; template: `%s`", err, attrs.Text)
		return
	}

	var buf bytes.Buffer
	buf.WriteString(PLUGIN_NAME)
	buf.WriteByte(':')

	err = tmpl.Execute(&buf, dict)
	if err != nil {
		err = fmt.Errorf("failed to execute the template: %w; template: `%s`; dict: `%s`", err, attrs.Text, jsontools.Dump(dict))
		return
	}
	resp = buf.String()
	return
}
