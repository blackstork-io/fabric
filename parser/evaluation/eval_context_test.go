package evaluation

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zclconf/go-cty/cty"
)

func Test_NewEvalContext(t *testing.T) {
	ctx := NewEvalContext()
	assert.NotNil(t, ctx)
	assert.NotNil(t, ctx.Functions)
	assert.NotNil(t, ctx.Functions["from_env_variable"])
}

func Test_makeFromEnvVariableFunc(t *testing.T) {
	t.Setenv("TEST_KEY", "test_value")
	fn := makeFromEnvVariableFunc()
	val, err := fn.Call([]cty.Value{cty.StringVal("TEST_KEY")})
	assert.Nil(t, err)
	assert.Equal(t, cty.StringVal("test_value"), val)
	val, err = fn.Call([]cty.Value{cty.StringVal("NON_EXISTENT_KEY")})
	assert.Nil(t, err)
	assert.Equal(t, cty.StringVal(""), val)
}
