package internal

import (
	"strings"
	"testing"

	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

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

// TestAllPluginSchemaValidity tests that all plugin schemas are valid
func TestAllPluginSchemaValidity(t *testing.T) {
	ver := "1.2.3"
	plugins := []*plugin.Schema{
		builtin.Plugin(ver),
		elastic.Plugin(ver),
		github.Plugin(ver, nil),
		graphql.Plugin(ver),
		openai.Plugin(ver, nil),
		opencti.Plugin(ver),
		postgresql.Plugin(ver),
		sqlite.Plugin(ver),
		terraform.Plugin(ver),
		hackerone.Plugin(ver, nil),
		virustotal.Plugin(ver, nil),
		stixview.Plugin(ver),
		splunk.Plugin(ver, nil),
	}
	for _, p := range plugins {
		p := p
		t.Run(p.Name, func(t *testing.T) {
			t.Parallel()
			assert.True(t, strings.HasPrefix(p.Name, "blackstork/"), "plugin name should be prefixed with 'blackstork/'")
			assert.Equal(t, ver, p.Version, "plugin version should match")
			assert.Greater(t, len(p.DataSources)+len(p.ContentProviders), 0, "plugin should have at least one data source or content provider")
			for name, ds := range p.DataSources {
				ds := ds
				t.Run(name, func(t *testing.T) {
					t.Parallel()
					validateDataSource(t, ds)
				})
			}
			for name, cp := range p.ContentProviders {
				cp := cp
				t.Run(name, func(t *testing.T) {
					t.Parallel()
					validateContentProvider(t, cp)
				})
			}
			assert.False(t, p.Validate().HasErrors(), "plugin should not have validation errors")
		})
	}
}

func validateDataSource(t testing.TB, ds *plugin.DataSource) {
	t.Helper()
	assert.NotNil(t, ds, "data source should not be nil")
	assert.NotEmpty(t, ds.DataFunc, "data source should have a data function")
	if ds.Config != nil {
		switch spec := ds.Config.(type) {
		case hcldec.ObjectSpec:
			assert.Greater(t, len(spec), 0, "data source config should have at least one attribute")
			for key, val := range spec {
				attr, ok := val.(*hcldec.AttrSpec)
				require.True(t, ok, "data source config attribute should be of type *hcldec.AttrSpec")
				validateAttrSpec(t, key, attr)
			}
		default:
			t.Errorf("unexpected data source config type: %T", ds.Config)
		}
	}
	if ds.Args != nil {
		switch spec := ds.Args.(type) {
		case hcldec.ObjectSpec:
			assert.Greater(t, len(spec), 0, "data source args should have at least one attribute")
			for key, val := range spec {
				attr, ok := val.(*hcldec.AttrSpec)
				require.True(t, ok, "data source args attribute should be of type *hcldec.AttrSpec")
				validateAttrSpec(t, key, attr)
			}
		default:
			t.Errorf("unexpected data source args type: %T", ds.Args)
		}
	}
}

func validateContentProvider(t testing.TB, cp *plugin.ContentProvider) {
	t.Helper()
	assert.NotNil(t, cp, "content provider should not be nil")
	assert.NotEmpty(t, cp.ContentFunc, "content provider should have a content function")
	if cp.Config != nil {
		switch spec := cp.Config.(type) {
		case hcldec.ObjectSpec:
			assert.Greater(t, len(spec), 0, "content provider config should have at least one attribute")
			for key, val := range spec {
				attr, ok := val.(*hcldec.AttrSpec)
				require.True(t, ok, "content provider config attribute should be of type *hcldec.AttrSpec")
				validateAttrSpec(t, key, attr)
			}
		default:
			t.Errorf("unexpected content provider config type: %T", cp.Config)
		}
	}
	if cp.Args != nil {
		switch spec := cp.Args.(type) {
		case hcldec.ObjectSpec:
			assert.Greater(t, len(spec), 0, "content provider args should have at least one attribute")
			for key, val := range spec {
				attr, ok := val.(*hcldec.AttrSpec)
				require.True(t, ok, "content provider args attribute should be of type *hcldec.AttrSpec")
				validateAttrSpec(t, key, attr)
			}
		default:
			t.Errorf("unexpected content provider args type: %T", cp.Args)
		}
	}
}

func validateAttrSpec(t testing.TB, key string, spec *hcldec.AttrSpec) {
	t.Helper()
	assert.Equal(t, key, spec.Name, "attribute name should match")
	assert.NotEmpty(t, spec.Name, "attribute name should not be empty")
	assert.NotEmpty(t, spec.Type, "attribute type should not be empty")
}
