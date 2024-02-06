package builtin

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPluginSchema(t *testing.T) {
	schema := Plugin("1.2.3")
	assert.Equal(t, "builtin", schema.Name)
	assert.Equal(t, "1.2.3", schema.Version)
	assert.NotNil(t, schema.DataSources["csv"])
	assert.NotNil(t, schema.DataSources["txt"])
	assert.NotNil(t, schema.DataSources["json"])
	assert.NotNil(t, schema.DataSources["inline"])
	assert.NotNil(t, schema.ContentProviders["text"])
	assert.NotNil(t, schema.ContentProviders["image"])
	assert.NotNil(t, schema.ContentProviders["list"])
	assert.NotNil(t, schema.ContentProviders["table"])
	assert.NotNil(t, schema.ContentProviders["frontmatter"])
}
