package main

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/spf13/pflag"

	"github.com/blackstork-io/fabric/internal/builtin"
	"github.com/blackstork-io/fabric/internal/elasticsearch"
	"github.com/blackstork-io/fabric/internal/github"
	"github.com/blackstork-io/fabric/internal/graphql"
	"github.com/blackstork-io/fabric/internal/openai"
	"github.com/blackstork-io/fabric/internal/postgresql"
	"github.com/blackstork-io/fabric/internal/sqlite"
	"github.com/blackstork-io/fabric/internal/terraform"
	"github.com/blackstork-io/fabric/plugin"
)

var (
	version   string
	outputDir string
)

//go:embed markdown.gotempl
var markdownTempl string

var templ *template.Template

func main() {
	// parse flags
	flags := pflag.NewFlagSet("docgen", pflag.ExitOnError)
	flags.StringVar(&version, "version", "v0.0.0-dev", "version of the build")
	flags.StringVar(&outputDir, "output", "./dist/docs", "output directory")
	flags.Parse(os.Args[1:])
	// ensure output directory exists
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		panic(err)
	}
	// load all plugins
	plugins := []*plugin.Schema{
		builtin.Plugin(version),
		elasticsearch.Plugin(version),
		github.Plugin(version, nil),
		graphql.Plugin(version),
		openai.Plugin(version, nil),
		postgresql.Plugin(version),
		sqlite.Plugin(version),
		terraform.Plugin(version),
	}
	// generate markdown for each plugin
	for _, p := range plugins {
		fp := filepath.Join(outputDir, fmt.Sprintf("%s.md", shortname(p.Name)))
		fmt.Printf("Generating '%s': '%s'\n", p.Name, fp)
		if err := generate(p, fp); err != nil {
			panic(err)
		}
	}
}

func generate(schema *plugin.Schema, fp string) error {
	f, err := os.Create(fp)
	if err != nil {
		return err
	}
	defer f.Close()
	return templ.Execute(f, schema)
}

func shortname(name string) string {
	parts := strings.SplitN(name, "/", 2)
	if len(parts) == 2 {
		return parts[1]
	}
	return name
}

func init() {
	templ = template.Must(template.New("markdown").Funcs(template.FuncMap{
		"shortname": shortname,
		"attrType": func(val hcldec.Spec) string {
			switch v := val.(type) {
			case *hcldec.AttrSpec:
				return v.Type.FriendlyName()
			default:
				return "unknown"
			}
		},
	}).Parse(markdownTempl))
}
