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
	return &Schema{
		Name:             src.Name,
		Version:          src.Version,
		DataSources:      dataSources,
		ContentProviders: contentProviders,
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
	args, err := encodeHclSpec(src.Args)
	if err != nil {
		return nil, err
	}
	config, err := encodeHclSpec(src.Config)
	if err != nil {
		return nil, err
	}
	return &DataSourceSchema{
		Args:   args,
		Config: config,
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
	args, err := encodeHclSpec(src.Args)
	if err != nil {
		return nil, err
	}
	config, err := encodeHclSpec(src.Config)
	if err != nil {
		return nil, err
	}
	return &ContentProviderSchema{
		Args:   args,
		Config: config,
	}, nil
}
