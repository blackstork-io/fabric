package github

import (
	"bytes"
	"context"
	"log/slog"

	gh "github.com/google/go-github/v58/github"
	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"
	"go.opentelemetry.io/otel/attribute"
	nooptrace "go.opentelemetry.io/otel/trace/noop"

	"github.com/blackstork-io/fabric/pkg/diagnostics"
	"github.com/blackstork-io/fabric/plugin"
	"github.com/blackstork-io/fabric/plugin/dataspec"
	"github.com/blackstork-io/fabric/plugin/dataspec/constraint"
	"github.com/blackstork-io/fabric/print"
	"github.com/blackstork-io/fabric/print/htmlprint"
	"github.com/blackstork-io/fabric/print/mdprint"
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
				Name:       "make_public",
				Type:       cty.Bool,
				DefaultVal: cty.False,
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

func publishGithubGist(loader ClientLoaderFn) plugin.PublishFunc {
	// TODO: confirm if to be passed from the caller
	logger := slog.Default()
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
		fileName := params.DocumentName + "." + params.Format.String()
		filenameAttr := params.Args.GetAttr("filename")
		if !filenameAttr.IsNull() && filenameAttr.AsString() != "" {
			fileName = filenameAttr.AsString()
		}
		payload := &gh.Gist{
			Public: gh.Bool(params.Args.GetAttr("make_public").True()),
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
			slog.InfoContext(ctx, "Created gist", "url", *gist.HTMLURL)
		} else {

			gist, _, err := client.Gists().Get(ctx, gistId.AsString())
			if err != nil {
				return diagnostics.Diag{{
					Severity: hcl.DiagError,
					Summary:  "Failed to retreive gist",
					Detail:   err.Error(),
				}}
			}
			// changing filename or output format will create a new file instead of updating the existing one.
			// following logic will remove the old files and add new files.
			for _, file := range gist.Files {
				_, exists := payload.Files[gh.GistFilename(*file.Filename)]
				if !exists {
					payload.Files[gh.GistFilename(*file.Filename)] = gh.GistFile{}
				}
			}
			gist, _, err = client.Gists().Edit(ctx, gistId.AsString(), payload)
			if err != nil {
				return diagnostics.Diag{{
					Severity: hcl.DiagError,
					Summary:  "Failed to update gist",
					Detail:   err.Error(),
				}}
			}
			slog.InfoContext(ctx, "Updated gist", "url", *gist.HTMLURL)
		}
		return nil
	}
}
