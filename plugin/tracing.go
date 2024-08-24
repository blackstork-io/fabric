package plugin

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	"github.com/blackstork-io/fabric/pkg/diagnostics"
	"github.com/blackstork-io/fabric/plugin/plugindata"
)

// WithTracing wraps a plugin schema with tracing instrumentation.
func WithTracing(plugin *Schema, tracer trace.Tracer) *Schema {
	plugin.ContentProviders = makeContentProvidersTracing(plugin.Name, plugin.ContentProviders, tracer)
	plugin.DataSources = makeDataSourcesTracing(plugin.Name, plugin.DataSources, tracer)
	plugin.Publishers = makePublishersTracing(plugin.Name, plugin.Publishers, tracer)
	return plugin
}

func makeContentProvidersTracing(plugin string, providers ContentProviders, tracer trace.Tracer) ContentProviders {
	result := make(ContentProviders)
	for name, provider := range providers {
		provider.ContentFunc = makeContentProviderTracing(plugin, name, provider, tracer)
		result[name] = provider
	}
	return result
}

func makeContentProviderTracing(plugin, name string, provider *ContentProvider, tracer trace.Tracer) ProvideContentFunc {
	next := provider.ContentFunc
	return func(ctx context.Context, params *ProvideContentParams) (_ *ContentResult, diags diagnostics.Diag) {
		ctx, span := tracer.Start(ctx, "ContentProvider.Execute", trace.WithAttributes(
			attribute.String("plugin", plugin),
			attribute.String("provider", name),
		))
		defer func() {
			if diags.HasErrors() {
				span.RecordError(diags)
				span.SetStatus(codes.Error, diags.Error())
			}
			span.End()
		}()
		return next(ctx, params)
	}
}

func makeDataSourcesTracing(plugin string, sources DataSources, tracer trace.Tracer) DataSources {
	result := make(DataSources)
	for name, source := range sources {
		source.DataFunc = makeDataSourceTracing(plugin, name, source, tracer)
		result[name] = source
	}
	return result
}

func makeDataSourceTracing(plugin, name string, source *DataSource, tracer trace.Tracer) RetrieveDataFunc {
	next := source.DataFunc
	return func(ctx context.Context, params *RetrieveDataParams) (_ plugindata.Data, diags diagnostics.Diag) {
		ctx, span := tracer.Start(ctx, "DataSource.Execute", trace.WithAttributes(
			attribute.String("plugin", plugin),
			attribute.String("datasource", name),
		))
		defer func() {
			if diags.HasErrors() {
				span.RecordError(diags)
				span.SetStatus(codes.Error, diags.Error())
			}
			span.End()
		}()
		return next(ctx, params)
	}
}

func makePublishersTracing(plugin string, publishers Publishers, tracer trace.Tracer) Publishers {
	result := make(Publishers)
	for name, publisher := range publishers {
		publisher.PublishFunc = makePublisherTracing(plugin, name, publisher, tracer)
		result[name] = publisher
	}
	return result
}

func makePublisherTracing(plugin, name string, publisher *Publisher, tracer trace.Tracer) PublishFunc {
	next := publisher.PublishFunc
	return func(ctx context.Context, params *PublishParams) (diags diagnostics.Diag) {
		ctx, span := tracer.Start(ctx, "Publisher.Execute", trace.WithAttributes(
			attribute.String("plugin", plugin),
			attribute.String("publisher", name),
		))
		defer func() {
			if diags.HasErrors() {
				span.RecordError(diags)
				span.SetStatus(codes.Error, diags.Error())
			}
			span.End()
		}()
		return next(ctx, params)
	}
}
