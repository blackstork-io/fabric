package pluginapiv1

import "github.com/blackstork-io/fabric/plugin"

func encodeSchema(src *plugin.Schema) (*Schema, error) {
	if src == nil {
		return nil, nil
	}
	dataSources, err := encodeDataSourceSchemaMap(src.DataSources)
	if err != nil {
		return nil, err
	}
	contentProviders, err := encodeContentProviderSchemaMap(src.ContentProviders)
	if err != nil {
		return nil, err
	}
	publishers, err := encodePublisherSchemaMap(src.Publishers)
	if err != nil {
		return nil, err
	}
	return &Schema{
		Name:             src.Name,
		Version:          src.Version,
		DataSources:      dataSources,
		ContentProviders: contentProviders,
		Publishers:       publishers,
		Doc:              src.Doc,
		Tags:             src.Tags,
	}, nil
}

func encodeDataSourceSchemaMap(src plugin.DataSources) (map[string]*DataSourceSchema, error) {
	dst := make(map[string]*DataSourceSchema, len(src))
	var err error
	for k, v := range src {
		dst[k], err = encodeDataSourceSchema(v)
		if err != nil {
			return nil, err
		}
	}
	return dst, nil
}

func encodeDataSourceSchema(src *plugin.DataSource) (*DataSourceSchema, error) {
	if src == nil {
		return nil, nil
	}
	args, err := encodeBlockSpec(src.Args.BlockSpec())
	if err != nil {
		return nil, err
	}
	config, err := encodeBlockSpec(src.Config.BlockSpec())
	if err != nil {
		return nil, err
	}
	return &DataSourceSchema{
		Args:   args,
		Config: config,
		Doc:    src.Doc,
		Tags:   src.Tags,
	}, nil
}

func encodeContentProviderSchemaMap(src plugin.ContentProviders) (map[string]*ContentProviderSchema, error) {
	dst := make(map[string]*ContentProviderSchema, len(src))
	var err error
	for k, v := range src {
		dst[k], err = encodeContentProviderSchema(v)
		if err != nil {
			return nil, err
		}
	}
	return dst, nil
}

func encodeContentProviderSchema(src *plugin.ContentProvider) (*ContentProviderSchema, error) {
	if src == nil {
		return nil, nil
	}
	args, err := encodeBlockSpec(src.Args.BlockSpec())
	if err != nil {
		return nil, err
	}
	config, err := encodeBlockSpec(src.Config.BlockSpec())
	if err != nil {
		return nil, err
	}
	return &ContentProviderSchema{
		Args:            args,
		Config:          config,
		InvocationOrder: encodeInvocationOrder(src.InvocationOrder),
		Doc:             src.Doc,
		Tags:            src.Tags,
	}, nil
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

func encodePublisherSchemaMap(src plugin.Publishers) (map[string]*PublisherSchema, error) {
	dst := make(map[string]*PublisherSchema, len(src))
	var err error
	for k, v := range src {
		dst[k], err = encodePublisherShema(v)
		if err != nil {
			return nil, err
		}
	}
	return dst, nil
}

func encodePublisherShema(src *plugin.Publisher) (*PublisherSchema, error) {
	if src == nil {
		return nil, nil
	}
	args, err := encodeBlockSpec(src.Args.BlockSpec())
	if err != nil {
		return nil, err
	}
	config, err := encodeBlockSpec(src.Config.BlockSpec())
	if err != nil {
		return nil, err
	}
	formats := make([]OutputFormat, len(src.AllowedFormats))
	for i, f := range src.AllowedFormats {
		formats[i] = encodeOutputFormat(f)
	}
	return &PublisherSchema{
		Args:           args,
		Config:         config,
		Doc:            src.Doc,
		Tags:           src.Tags,
		AllowedFormats: formats,
	}, nil
}
