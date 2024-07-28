package github

import (
	"bytes"
	"context"
	"io"
	"log/slog"

	gh "github.com/google/go-github/v58/github"
	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"
	"go.opentelemetry.io/otel/attribute"

	"github.com/blackstork-io/fabric/pkg/diagnostics"
	"github.com/blackstork-io/fabric/plugin"
	"github.com/blackstork-io/fabric/plugin/dataspec"
	"github.com/blackstork-io/fabric/plugin/dataspec/constraint"
	"github.com/blackstork-io/fabric/print"
	"github.com/blackstork-io/fabric/print/htmlprint"
	"github.com/blackstork-io/fabric/print/mdprint"
	nooptrace "go.opentelemetry.io/otel/trace/noop"
)

func makeGithubGistPublisher(loader ClientLoaderFn) *plugin.Publisher {
	return &plugin.Publisher{
		Doc:  "Publishes content to github gist",
		Tags: []string{},
		Config: dataspec.ObjectSpec{
			&dataspec.AttrSpec{
				Name:        "github_token",
				Type:        cty.String,
				Constraints: constraint.RequiredNonNull,
				Secret:      true,
			},
		},
		Args: dataspec.ObjectSpec{
			&dataspec.AttrSpec{
				Name: "description",
				Type: cty.String,
			},
			&dataspec.AttrSpec{
				Name: "filename",
				Type: cty.String,
			},
			&dataspec.AttrSpec{
				Name: "make_public",
				Type: cty.Bool,
			},
			&dataspec.AttrSpec{
				Name: "gist_id",
				Type: cty.String,
			},
		},
		AllowedFormats: []plugin.OutputFormat{plugin.OutputFormatMD, plugin.OutputFormatHTML},
		PublishFunc:    publishGithubGist(loader),
	}
}

func parseContent(data plugin.MapData) (document *plugin.ContentSection) {
	documentMap, ok := data["document"]
	if !ok {
		return
	}
	contentMap, ok := documentMap.(plugin.MapData)["content"]
	if !ok {
		return
	}
	content, err := plugin.ParseContentData(contentMap.(plugin.MapData))
	if err != nil {
		return
	}
	document = content.(*plugin.ContentSection)
	return
}

func parseName(data plugin.MapData) (name string) {
	documentMap, ok := data["document"]
	if !ok {
		return
	}
	metaMap, ok := documentMap.(plugin.MapData)["meta"]
	if !ok {
		return
	}
	documentName := metaMap.(plugin.MapData)["name"]
	docName, ok := documentName.(plugin.StringData)
	if !ok {
		return
	}
	name = string(docName)
	return
}

func publishGithubGist(loader ClientLoaderFn) plugin.PublishFunc {
	// TODO: confirm if to be passed from the caller
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	tracer := nooptrace.Tracer{}
	return func(ctx context.Context, params *plugin.PublishParams) diagnostics.Diag {
		document := parseContent(params.DataContext)
		if document == nil {
			return diagnostics.Diag{{
				Severity: hcl.DiagError,
				Summary:  "Failed to parse document",
				Detail:   "document is required",
			}}
		}
		datactx := params.DataContext
		datactx["format"] = plugin.StringData(params.Format.String())
		var printer print.Printer
		switch params.Format {
		case plugin.OutputFormatMD:
			printer = mdprint.New()
		case plugin.OutputFormatHTML:
			printer = htmlprint.New()
		default:
			return diagnostics.Diag{{
				Severity: hcl.DiagError,
				Summary:  "Unsupported format",
				Detail:   "Only md and html formats are supported",
			}}
		}
		printer = print.WithLogging(printer, logger, slog.String("format", params.Format.String()))
		printer = print.WithTracing(printer, tracer, attribute.String("format", params.Format.String()))

		buff := bytes.NewBuffer(nil)
		err := printer.Print(ctx, buff, document)
		if err != nil {
			return diagnostics.Diag{{
				Severity: hcl.DiagError,
				Summary:  "Failed to write to a file",
				Detail:   err.Error(),
			}}
		}

		client := loader(params.Config.GetAttr("github_token").AsString())

		fileName := parseName(params.DataContext) + "." + params.Format.String()
		filenameAttr := params.Args.GetAttr("filename")
		if !filenameAttr.IsNull() && filenameAttr.AsString() != "" {
			fileName = filenameAttr.AsString()
		}
		payload := &gh.Gist{
			Public: gh.Bool(false),
			Files: map[gh.GistFilename]gh.GistFile{
				gh.GistFilename(fileName): {
					Content:  gh.String(buff.String()),
					Filename: gh.String(fileName),
				},
			},
		}
		// overrides params if defined
		descriptionAttr := params.Args.GetAttr("description")
		if !descriptionAttr.IsNull() && descriptionAttr.AsString() != "" {
			payload.Description = gh.String(descriptionAttr.AsString())
		}
		makePublicAttr := params.Args.GetAttr("make_public")
		if !makePublicAttr.IsNull() {
			payload.Public = gh.Bool(makePublicAttr.True())
		}
		slog.InfoContext(ctx, "Publish to github gist", "filename", fileName)
		gistId := params.Args.GetAttr("gist_id")
		if gistId.IsNull() || gistId.AsString() == "" {
			gist, _, err := client.Gists().Create(ctx, payload)
			if err != nil {
				return diagnostics.Diag{{
					Severity: hcl.DiagError,
					Summary:  "Failed to create gist",
					Detail:   err.Error(),
				}}
			}
			slog.InfoContext(ctx, "created gist", "url", *gist.HTMLURL)
		} else {
			gist, _, err := client.Gists().Edit(ctx, gistId.AsString(), payload)
			if err != nil {
				return diagnostics.Diag{{
					Severity: hcl.DiagError,
					Summary:  "Failed to update gist",
					Detail:   err.Error(),
				}}
			}
			slog.InfoContext(ctx, "updated gist", "url", *gist.HTMLURL)
		}
		return nil
	}
}
