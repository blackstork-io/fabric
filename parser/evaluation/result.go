package evaluation

import "github.com/blackstork-io/fabric/plugin"

type Result plugin.ListData

func (d *Result) AsJQData() plugin.Data {
	return plugin.ListData(*d)
}

func (d Result) AsGoType() (result []string) {
	result = make([]string, len(d))
	for i, s := range d {
		// The only way to modify resultsList is through append, so this is always correct
		result[i] = string(s.(plugin.StringData))
	}
	return
}

func (d *Result) Append(s string) {
	*d = append(*d, plugin.StringData(s))
}
