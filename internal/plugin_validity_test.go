package internal

import (
	"strings"
	"testing"

	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/stretchr/testify/assert"

	"github.com/blackstork-io/fabric/internal/builtin"
	"github.com/blackstork-io/fabric/internal/elastic"
	"github.com/blackstork-io/fabric/internal/github"
	"github.com/blackstork-io/fabric/internal/graphql"
	"github.com/blackstork-io/fabric/internal/hackerone"
	"github.com/blackstork-io/fabric/internal/nistnvd"
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
		elastic.Plugin(ver, nil),
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
		nistnvd.Plugin(ver, nil),
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
			for name, pub := range p.Publishers {
				pub := pub
				t.Run(name, func(t *testing.T) {
					t.Parallel()
					validatePublisher(t, pub)
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
		assert.False(t, ds.Config.IsEmpty(), "data source config should have at least one attribute")
		assert.Empty(t, ds.Config.ValidateSpec(), "data source config validation errors")
	}
	if ds.Args != nil {
		assert.False(t, ds.Args.IsEmpty(), "data source args should have at least one attribute")
		assert.Empty(t, ds.Args.ValidateSpec(), "data source args validation errors")
	}
}

func validateContentProvider(t testing.TB, cp *plugin.ContentProvider) {
	t.Helper()
	assert.NotNil(t, cp, "content provider should not be nil")
	assert.NotEmpty(t, cp.ContentFunc, "content provider should have a content function")
	if cp.Config != nil {
		assert.False(t, cp.Config.IsEmpty(), "content provider config should have at least one attribute")
		assert.Empty(t, cp.Config.ValidateSpec(), "content provider config validation errors")

	}
	if cp.Args != nil {
		assert.False(t, cp.Args.IsEmpty(), "content provider args should have at least one attribute")
		assert.Empty(t, cp.Args.ValidateSpec(), "content provider args validation errors")
	}
}

func validatePublisher(t testing.TB, pub *plugin.Publisher) {
	t.Helper()
	assert.NotNil(t, pub, "publisher should not be nil")
	assert.NotEmpty(t, pub.PublishFunc, "publisher should have a publish function")
	if pub.Config != nil {
		assert.False(t, pub.Config.IsEmpty(), "publisher config should have at least one attribute")
		assert.Empty(t, pub.Config.Validate(), "publisher config validation errors")
	}
}

func validateAttrSpec(t testing.TB, key string, spec *hcldec.AttrSpec) {
	t.Helper()
	assert.Equal(t, key, spec.Name, "attribute name should match")
	assert.NotEmpty(t, spec.Name, "attribute name should not be empty")
	assert.NotEmpty(t, spec.Type, "attribute type should not be empty")
}
