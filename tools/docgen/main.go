package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/spf13/pflag"

	"github.com/blackstork-io/fabric/internal/builtin"
	"github.com/blackstork-io/fabric/internal/elastic"
	"github.com/blackstork-io/fabric/internal/github"
	"github.com/blackstork-io/fabric/internal/graphql"
	"github.com/blackstork-io/fabric/internal/hackerone"
	"github.com/blackstork-io/fabric/internal/openai"
	"github.com/blackstork-io/fabric/internal/opencti"
	"github.com/blackstork-io/fabric/internal/postgresql"
	"github.com/blackstork-io/fabric/internal/splunk"
	"github.com/blackstork-io/fabric/internal/sqlite"
	"github.com/blackstork-io/fabric/internal/stixview"
	"github.com/blackstork-io/fabric/internal/terraform"
	"github.com/blackstork-io/fabric/internal/virustotal"
	"github.com/blackstork-io/fabric/plugin"
)

var (
	version   string
	outputDir string
)

//go:embed content-provider.md.gotempl
var contentProviderTemplValue string

var contentProviderTempl *template.Template

//go:embed data-source.md.gotempl
var dataSourceTemplValue string

var dataSourceTempl *template.Template

//go:embed plugin.md.gotempl
var pluginTemplValue string

var pluginTempl *template.Template

type PluginResourceMeta struct {
	Name         string   `json:"name"`
	Type         string   `json:"type"`
	ConfigParams []string `json:"config_params,omitempty"`
	Arguments    []string `json:"arguments,omitempty"`
}

func generateDataSourceDocs(log *slog.Logger, p *plugin.Schema, outputDir string) {
	log.Info("Found data sources inside the plugin", "count", len(p.DataSources))

	dataSourcesDir := filepath.Join(outputDir, "data-sources")

	// Create a directory for plugin's data sources if it doesn't exist
	err := os.MkdirAll(dataSourcesDir, 0766)
	if err != nil {
		log.Error("Can't create a directory", "path", dataSourcesDir)
		panic(err)
	}

	for name, dataSource := range p.DataSources {
		log.Info("Found a data source", "name", name)
		docFilename := fmt.Sprintf("%s.md", name)
		docPath := filepath.Join(dataSourcesDir, docFilename)
		err := renderDataSourceDoc(p, name, dataSource, docPath)
		if err != nil {
			log.Error("Error while rendering a data source doc", "plugin", p.Name, "dataSource", name)
			panic(err)
		}
	}
}

func generateContentProviderDocs(log *slog.Logger, p *plugin.Schema, outputDir string) {
	log.Info("Found content providers inside the plugin", "count", len(p.ContentProviders))

	contentProvidersDir := filepath.Join(outputDir, "content-providers")

	// Create a directory for plugin's content providers if it doesn't exist
	err := os.MkdirAll(contentProvidersDir, 0766)
	if err != nil {
		log.Error("Can't create a directory", "path", contentProvidersDir)
		panic(err)
	}

	for name, contentProvider := range p.ContentProviders {
		log.Info("Found a content provider", "name", name)

		docFilename := fmt.Sprintf("%s.md", name)
		docPath := filepath.Join(contentProvidersDir, docFilename)
		err := renderContentProviderDoc(p, name, contentProvider, docPath)
		if err != nil {
			log.Error("Error while rendering a content provider doc", "plugin", p.Name, "contentProvider", name)
			panic(err)
		}
	}
}

func marshalDataSource(name string, ds *plugin.DataSource) PluginResourceMeta {

	var configParams []string
	configSpec, ok := ds.Config.(hcldec.ObjectSpec)
	if ok && configSpec != nil {
		for k := range configSpec {
			configParams = append(configParams, k)
		}
	}

	var arguments []string
	argsSpec, ok := ds.Args.(hcldec.ObjectSpec)
	if ok && argsSpec != nil {
		for k := range argsSpec {
			arguments = append(arguments, k)
		}
	}
	return PluginResourceMeta{
		Name:         name,
		Type:         "data-source",
		ConfigParams: configParams,
		Arguments:    arguments,
	}
}

func marshalContentProvider(name string, p *plugin.ContentProvider) PluginResourceMeta {

	var configParams []string
	configSpec, ok := p.Config.(hcldec.ObjectSpec)
	if ok && configSpec != nil {
		for k := range configSpec {
			configParams = append(configParams, k)
		}
	}

	var arguments []string
	argsSpec, ok := p.Args.(hcldec.ObjectSpec)
	if ok && argsSpec != nil {
		for k := range argsSpec {
			arguments = append(arguments, k)
		}
	}
	return PluginResourceMeta{
		Name:         name,
		Type:         "content-provider",
		ConfigParams: configParams,
		Arguments:    arguments,
	}
}

func generateMetadataFile(plugins []*plugin.Schema, outputDir string) {
	pluginDetails := make([]any, len(plugins))

	for i, p := range plugins {

		var resources []PluginResourceMeta

		for name, dataSource := range p.DataSources {
			resources = append(resources, marshalDataSource(name, dataSource))
		}
		for name, contentProvider := range p.ContentProviders {
			resources = append(resources, marshalContentProvider(name, contentProvider))
		}

		pluginDetails[i] = map[string]any{
			"name":      p.Name,
			"version":   p.Version,
			"shortname": shortname(p.Name),
			"resources": resources,
		}
	}

	jsonData, err := json.MarshalIndent(pluginDetails, "", "  ")
	if err != nil {
		slog.Error("Failed to marshal the plugin details into JSON")
		return
	}

	pluginDetailsPath := filepath.Join(outputDir, "plugins.json")
	file, err := os.Create(pluginDetailsPath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	file.Write(jsonData)

	slog.Info("Plugin details file generated", "path", pluginDetailsPath)
}

func main() {
	// parse flags
	flags := pflag.NewFlagSet("docgen", pflag.ExitOnError)
	flags.StringVar(&version, "version", "v0.0.0-dev", "version of the build")
	flags.StringVar(&outputDir, "output", "./dist/docs", "output directory")
	flags.Parse(os.Args[1:])
	// load all plugins
	plugins := []*plugin.Schema{
		builtin.Plugin(version),
		elastic.Plugin(version),
		github.Plugin(version, nil),
		graphql.Plugin(version),
		openai.Plugin(version, nil),
		opencti.Plugin(version),
		postgresql.Plugin(version),
		sqlite.Plugin(version),
		terraform.Plugin(version),
		hackerone.Plugin(version, nil),
		virustotal.Plugin(version, nil),
		splunk.Plugin(version, nil),
		stixview.Plugin(version),
	}
	// generate markdown for each plugin
	for _, p := range plugins {

		log := slog.With("plugin", p.Name)

		pluginShortname := shortname(p.Name)

		// Use a shortname as a plugin directory name
		pluginOutputDir := filepath.Join(outputDir, pluginShortname)

		// Create a plugin directory if it doesn't exist
		err := os.MkdirAll(pluginOutputDir, 0766)
		if err != nil {
			log.Error("Can't create a plugin directory", "path", pluginOutputDir)
			panic(err)
		}

		pluginDocPath := filepath.Join(pluginOutputDir, "_index.md")
		err = renderPluginDoc(p, pluginDocPath)
		if err != nil {
			log.Error("Error while rendering a plugin doc", "plugin", p.Name)
			panic(err)
		}

		log.Info("Plugin doc rendered", "path", pluginDocPath)

		if len(p.DataSources) != 0 {
			generateDataSourceDocs(log, p, pluginOutputDir)
		}

		if len(p.ContentProviders) != 0 {
			generateContentProviderDocs(log, p, pluginOutputDir)
		}
	}
	generateMetadataFile(plugins, outputDir)
}

func renderPluginDoc(pluginSchema *plugin.Schema, fp string) error {
	f, err := os.Create(fp)
	if err != nil {
		return err
	}
	defer f.Close()

	return pluginTempl.Execute(f, pluginSchema)
}

func renderContentProviderDoc(pluginSchema *plugin.Schema, contentProviderName string, contentProvider *plugin.ContentProvider, fp string) error {
	f, err := os.Create(fp)
	if err != nil {
		return err
	}
	defer f.Close()

	var templContext = map[string]any{
		"plugin":           pluginSchema,
		"plugin_shortname": shortname(pluginSchema.Name),
		"name":             contentProviderName,
		"content_provider": contentProvider,
	}
	return contentProviderTempl.Execute(f, templContext)
}

func renderDataSourceDoc(pluginSchema *plugin.Schema, dataSourceName string, dataSource *plugin.DataSource, fp string) error {
	f, err := os.Create(fp)
	if err != nil {
		return err
	}
	defer f.Close()

	var templContext = map[string]any{
		"plugin":           pluginSchema,
		"plugin_shortname": shortname(pluginSchema.Name),
		"name":             dataSourceName,
		"data_source":      dataSource,
	}
	return dataSourceTempl.Execute(f, templContext)
}

func shortname(name string) string {
	parts := strings.SplitN(name, "/", 2)
	if len(parts) == 2 {
		return parts[1]
	}
	return name
}

func templateAttrTypeFunc(val hcldec.Spec) string {
	switch v := val.(type) {
	case *hcldec.AttrSpec:
		return v.Type.FriendlyName()
	default:
		return "unknown"
	}
}

func init() {
	contentProviderTempl = template.Must(template.New("content-provider").Funcs(template.FuncMap{
		"attrType": templateAttrTypeFunc,
	}).Parse(contentProviderTemplValue))
	dataSourceTempl = template.Must(template.New("data-source").Funcs(template.FuncMap{
		"attrType": templateAttrTypeFunc,
	}).Parse(dataSourceTemplValue))
	pluginTempl = template.Must(template.New("plugin").Funcs(template.FuncMap{
		"attrType":  templateAttrTypeFunc,
		"shortname": shortname,
	}).Parse(pluginTemplValue))
}
