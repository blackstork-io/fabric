package pluginapiv1

import (
	"github.com/blackstork-io/fabric/pkg/diagnostics"
	"github.com/blackstork-io/fabric/pkg/utils"
	"github.com/blackstork-io/fabric/plugin"
)

func encodeSchema(src *plugin.Schema) (*Schema, diagnostics.Diag) {
	if src == nil {
		return nil, nil
	}
	var diags diagnostics.Diag

	return &Schema{
		Name:             src.Name,
		Version:          src.Version,
		DataSources:      utils.MapMapDiags(&diags, src.DataSources, encodeDataSourceSchema),
		ContentProviders: utils.MapMapDiags(&diags, src.ContentProviders, encodeContentProviderSchema),
		Publishers:       utils.MapMapDiags(&diags, src.Publishers, encodePublisherShema),
		Doc:              src.Doc,
		Tags:             src.Tags,
	}, diags
}

func encodeDataSourceSchema(src *plugin.DataSource) (_ *DataSourceSchema, diags diagnostics.Diag) {
	if src == nil {
		return nil, nil
	}
	schema := &DataSourceSchema{
		Doc:  src.Doc,
		Tags: src.Tags,
	}
	var diag diagnostics.Diag
	if src.Args != nil {
		schema.Args, diag = encodeRootSpec(src.Args)
		diags.Extend(diag)
	}
	if src.Config != nil {
		schema.Config, diag = encodeRootSpec(src.Config)
		diags.Extend(diag)
	}
	return schema, diags
}

func encodeContentProviderSchema(src *plugin.ContentProvider) (_ *ContentProviderSchema, diags diagnostics.Diag) {
	if src == nil {
		return nil, nil
	}
	schema := &ContentProviderSchema{
		InvocationOrder: encodeInvocationOrder(src.InvocationOrder),
		Doc:             src.Doc,
		Tags:            src.Tags,
	}
	var diag diagnostics.Diag
	if src.Args != nil {
		schema.Args, diag = encodeRootSpec(src.Args)
		diags.Extend(diag)
	}
	if src.Config != nil {
		schema.Config, diag = encodeRootSpec(src.Config)
		diags.Extend(diag)
	}
	return schema, diags
}

func encodeInvocationOrder(src plugin.InvocationOrder) InvocationOrder {
	switch src {
	case plugin.InvocationOrderBegin:
		return InvocationOrder_INVOCATION_ORDER_BEGIN
	case plugin.InvocationOrderEnd:
		return InvocationOrder_INVOCATION_ORDER_END
	default:
		return InvocationOrder_INVOCATION_ORDER_UNSPECIFIED
	}
}

func encodeOutputFormat(src plugin.OutputFormat) OutputFormat {
	switch src {
	case plugin.OutputFormatMD:
		return OutputFormat_OUTPUT_FORMAT_MD
	case plugin.OutputFormatHTML:
		return OutputFormat_OUTPUT_FORMAT_HTML
	case plugin.OutputFormatPDF:
		return OutputFormat_OUTPUT_FORMAT_PDF
	default:
		return OutputFormat_OUTPUT_FORMAT_UNSPECIFIED
	}
}

func encodePublisherShema(src *plugin.Publisher) (_ *PublisherSchema, diags diagnostics.Diag) {
	if src == nil {
		return nil, nil
	}
	schema := &PublisherSchema{
		Doc:            src.Doc,
		Tags:           src.Tags,
		AllowedFormats: utils.FnMap(src.AllowedFormats, encodeOutputFormat),
	}

	var diag diagnostics.Diag
	if src.Args != nil {
		schema.Args, diag = encodeRootSpec(src.Args)
		diags.Extend(diag)
	}
	if src.Config != nil {
		schema.Config, diag = encodeRootSpec(src.Config)
		diags.Extend(diag)
	}
	return schema, diags
}
