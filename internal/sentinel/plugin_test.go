package sentinel

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPlugin_Schema(t *testing.T) {
	schema := Plugin("1.2.3", nil)
	assert.Equal(t, "blackstork/microsoft_sentinel", schema.Name)
	assert.Equal(t, "1.2.3", schema.Version)
	assert.NotNil(t, schema.DataSources["microsoft_sentinel_incidents"])
}
