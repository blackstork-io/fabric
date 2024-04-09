package pluginapiv1

import "github.com/blackstork-io/fabric/plugin"

func decodeSchema(src *Schema) (*plugin.Schema, error) {
	if src == nil {
		return nil, nil
	}
	dataSources, err := decodeDataSourceSchemaMap(src.GetDataSources())
	if err != nil {
		return nil, err
	}
	contentProviders, err := decodeContentProviderSchemaMap(src.GetContentProviders())
	if err != nil {
		return nil, err
	}
	return &plugin.Schema{
		Name:             src.GetName(),
		Version:          src.GetVersion(),
		DataSources:      dataSources,
		ContentProviders: contentProviders,
		Doc:              src.GetDoc(),
		Tags:             src.GetTags(),
	}, nil
}

func decodeDataSourceSchemaMap(src map[string]*DataSourceSchema) (plugin.DataSources, error) {
	dst := make(plugin.DataSources, len(src))
	var err error
	for k, v := range src {
		dst[k], err = decodeDataSourceSchema(v)
		if err != nil {
			return nil, err
		}
	}
	return dst, nil
}

func decodeDataSourceSchema(src *DataSourceSchema) (*plugin.DataSource, error) {
	if src == nil {
		return nil, nil
	}
	args, err := decodeRootSpec(src.GetArgs())
	if err != nil {
		return nil, err
	}
	config, err := decodeRootSpec(src.GetConfig())
	if err != nil {
		return nil, err
	}

	return &plugin.DataSource{
		Args:   args,
		Config: config,
		Doc:    src.GetDoc(),
		Tags:   src.GetTags(),
	}, nil
}

func decodeContentProviderSchemaMap(src map[string]*ContentProviderSchema) (plugin.ContentProviders, error) {
	dst := make(plugin.ContentProviders, len(src))
	var err error
	for k, v := range src {
		dst[k], err = decodeContentProviderSchema(v)
		if err != nil {
			return nil, err
		}
	}
	return dst, nil
}

func decodeContentProviderSchema(src *ContentProviderSchema) (*plugin.ContentProvider, error) {
	if src == nil {
		return nil, nil
	}
	args, err := decodeRootSpec(src.GetArgs())
	if err != nil {
		return nil, err
	}
	config, err := decodeRootSpec(src.GetConfig())
	if err != nil {
		return nil, err
	}
	return &plugin.ContentProvider{
		Args:            args,
		Config:          config,
		InvocationOrder: decodeInvocationOrder(src.GetInvocationOrder()),
		Doc:             src.GetDoc(),
		Tags:            src.GetTags(),
	}, nil
}

func decodeInvocationOrder(src InvocationOrder) plugin.InvocationOrder {
	switch src {
	case InvocationOrder_INVOCATION_ORDER_BEGIN:
		return plugin.InvocationOrderBegin
	case InvocationOrder_INVOCATION_ORDER_END:
		return plugin.InvocationOrderEnd
	default:
		return plugin.InvocationOrderUnspecified
	}
}
