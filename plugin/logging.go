package plugin

import (
	"context"
	"log/slog"
	"strconv"

	"github.com/zclconf/go-cty/cty"

	"github.com/blackstork-io/fabric/pkg/diagnostics"
	"github.com/blackstork-io/fabric/plugin/dataspec"
)

// WithLogging wraps the plugin with logging instrumentation.
func WithLogging(plugin *Schema, logger *slog.Logger) *Schema {
	logger = logger.With("component", "plugin")
	plugin.ContentProviders = makeContentProvidersLogging(plugin.Name, plugin.ContentProviders, logger)
	plugin.DataSources = makeDataSourcesLogging(plugin.Name, plugin.DataSources, logger)
	plugin.Publishers = makePublishersLogging(plugin.Name, plugin.Publishers, logger)
	return plugin
}

func makeContentProvidersLogging(plugin string, providers ContentProviders, logger *slog.Logger) ContentProviders {
	result := make(ContentProviders)
	for name, provider := range providers {
		provider.ContentFunc = makeContentProviderLogging(plugin, name, *provider, logger)
		result[name] = provider
	}
	return result
}

func makeDataSourcesLogging(plugin string, sources DataSources, logger *slog.Logger) DataSources {
	result := make(DataSources)
	for name, source := range sources {
		source.DataFunc = makeDataSourceLogging(plugin, name, *source, logger)
		result[name] = source
	}
	return result
}

func makePublishersLogging(plugin string, publishers Publishers, logger *slog.Logger) Publishers {
	result := make(Publishers)
	for name, publisher := range publishers {
		publisher.PublishFunc = makePublisherLogging(plugin, name, *publisher, logger)
		result[name] = publisher
	}
	return result
}

func makePublisherLogging(plugin, name string, publisher Publisher, logger *slog.Logger) PublishFunc {
	next := publisher.PublishFunc
	return func(ctx context.Context, params *PublishParams) diagnostics.Diag {
		logger.DebugContext(ctx, "Executing publisher", "params", slog.GroupValue(
			slog.String("plugin", plugin),
			slog.String("publisher", name),
			slog.String("format", params.Format.String()),
			slog.Any("config", logDataSpecValue(publisher.Config, params.Config)),
			slog.Any("args", logDataSpecValue(publisher.Args, params.Args)),
			slog.String("document_name", params.DocumentName),
		))
		return next(ctx, params)
	}
}

func makeContentProviderLogging(plugin, name string, provider ContentProvider, logger *slog.Logger) ProvideContentFunc {
	next := provider.ContentFunc
	return func(ctx context.Context, params *ProvideContentParams) (*ContentResult, diagnostics.Diag) {
		logger.DebugContext(ctx, "Executing content provider", "params", slog.GroupValue(
			slog.String("plugin", plugin),
			slog.String("provider", name),
			slog.Any("config", logDataSpecValue(provider.Config, params.Config)),
			slog.Any("args", logDataSpecValue(provider.Args, params.Args)),
			slog.Uint64("content_id", uint64(params.ContentID)),
		))
		return next(ctx, params)
	}
}

func makeDataSourceLogging(plugin, name string, source DataSource, logger *slog.Logger) RetrieveDataFunc {
	next := source.DataFunc
	return func(ctx context.Context, params *RetrieveDataParams) (Data, diagnostics.Diag) {
		logger.DebugContext(ctx, "Executing datasource", "params", slog.GroupValue(
			slog.String("plugin", plugin),
			slog.String("datasource", name),
			slog.Any("config", logDataSpecValue(source.Config, params.Config)),
			slog.Any("args", logDataSpecValue(source.Args, params.Args)),
		))
		return next(ctx, params)
	}
}

func logDataSpecValue(spec dataspec.Spec, value cty.Value) slog.Value {
	switch spec := spec.(type) {
	case dataspec.ObjectSpec:
		return logDataSpecObjectValue(spec, value)
	case *dataspec.ObjectSpec:
		return logDataSpecObjectValue(*spec, value)
	case *dataspec.AttrSpec:
		return logDataSpecAttrValue(spec, value)
	case *dataspec.BlockSpec:
		return logDataSpecBlockValue(spec, value)
	default:
		return slog.Value{}
	}
}

func logDataSpecObjectValue(spec dataspec.ObjectSpec, value cty.Value) slog.Value {
	if value.IsNull() || (!value.Type().IsObjectType() && !value.Type().IsMapType()) {
		return slog.Value{}
	}
	m := value.AsValueMap()
	attrs := []slog.Attr{}
	for _, attr := range spec {
		v, ok := m[attr.KeyForObjectSpec()]
		if !ok {
			continue
		}
		attrs = append(attrs, slog.Attr{
			Key:   attr.KeyForObjectSpec(),
			Value: logDataSpecValue(attr, v),
		})
	}
	return slog.GroupValue(attrs...)
}

func logDataSpecBlockValue(spec *dataspec.BlockSpec, value cty.Value) slog.Value {
	return logDataSpecValue(spec.Nested, value)
}

func logDataSpecAttrValue(spec *dataspec.AttrSpec, value cty.Value) slog.Value {
	if value.IsNull() {
		return slog.Value{}
	}
	if spec.Secret {
		return slog.StringValue("<secret>")
	}
	switch {
	case value.Type() == cty.String:
		return slog.StringValue(value.AsString())
	case value.Type() == cty.Number:
		f, _ := value.AsBigFloat().Float64()
		return slog.Float64Value(f)
	case value.Type() == cty.Bool:
		return slog.BoolValue(value.True())
	case value.Type().IsListType() || value.Type().IsTupleType() || value.Type().IsSetType():
		return logCtyListValue(value.AsValueSlice())
	case value.Type().IsMapType() || value.Type().IsObjectType():
		return logCtyMapValue(value.AsValueMap())
	default:
		return slog.Value{}
	}
}

func logCtyListValue(values []cty.Value) slog.Value {
	attrs := []slog.Attr{}
	for i, v := range values {
		k := strconv.Itoa(i)
		switch {
		case v.Type() == cty.String:
			attrs = append(attrs, slog.String(k, v.AsString()))
		case v.Type() == cty.Number:
			f, _ := v.AsBigFloat().Float64()
			attrs = append(attrs, slog.Float64(k, f))
		case v.Type() == cty.Bool:
			attrs = append(attrs, slog.Bool(k, v.True()))
		case v.Type().IsListType() || v.Type().IsTupleType() || v.Type().IsSetType():
			attrs = append(attrs, slog.Any(k, logCtyListValue(v.AsValueSlice())))
		case v.Type().IsMapType() || v.Type().IsObjectType():
			attrs = append(attrs, slog.Any(k, logCtyMapValue(v.AsValueMap())))
		default:
			attrs = append(attrs, slog.String(k, "<unknown>"))
		}
	}
	return slog.GroupValue(attrs...)
}

func logCtyMapValue(m map[string]cty.Value) slog.Value {
	attrs := []slog.Attr{}
	for k, v := range m {
		switch {
		case v.Type() == cty.String:
			attrs = append(attrs, slog.String(k, v.AsString()))
		case v.Type() == cty.Number:
			f, _ := v.AsBigFloat().Float64()
			attrs = append(attrs, slog.Float64(k, f))
		case v.Type() == cty.Bool:
			attrs = append(attrs, slog.Bool(k, v.True()))
		case v.Type().IsListType() || v.Type().IsTupleType() || v.Type().IsSetType():
			attrs = append(attrs, slog.Any(k, logCtyListValue(v.AsValueSlice())))
		case v.Type().IsMapType() || v.Type().IsObjectType():
			attrs = append(attrs, slog.Any(k, logCtyMapValue(v.AsValueMap())))
		default:
			attrs = append(attrs, slog.String(k, "<unknown>"))
		}
	}
	return slog.GroupValue(attrs...)
}
