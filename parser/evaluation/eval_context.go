package evaluation

import (
	"log/slog"
	"os"
	"strings"
	"sync"

	"github.com/hashicorp/hcl/v2"
	"github.com/joho/godotenv"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"
)

func buildEnvVarMap() cty.Value {
	envVars := os.Environ()
	envFromFile, err := godotenv.Read()
	if err != nil && !os.IsNotExist(err) {
		slog.Error("Error reading .env file", "err", err)
	}
	envMap := make(map[string]cty.Value, len(envFromFile)+len(envVars))
	if envFromFile != nil {
		for k, v := range envFromFile {
			envMap[k] = cty.StringVal(v)
		}
	}
	for k, v := range envFromFile {
		envMap[k] = cty.StringVal(v)
	}
	for _, e := range envVars {
		pair := strings.SplitN(e, "=", 2)
		if len(pair) != 2 {
			continue
		}
		envMap[pair[0]] = cty.StringVal(pair[1])
	}
	return cty.MapVal(envMap)
}

var makeFromEnvVariableFunc = function.New(&function.Spec{
	Params: []function.Parameter{
		{
			Name: "name",
			Type: cty.String,
		},
	},
	Type: function.StaticReturnType(cty.String),
	Impl: func(args []cty.Value, retType cty.Type) (cty.Value, error) {
		name := args[0].AsString()
		return cty.StringVal(os.Getenv(name)), nil
	},
})

func buildEvalContext() *hcl.EvalContext {
	return &hcl.EvalContext{
		Functions: map[string]function.Function{
			"from_env_variable": makeFromEnvVariableFunc,
		},
		Variables: map[string]cty.Value{
			"env": buildEnvVarMap(),
		},
	}
}

var EvalContext = sync.OnceValue(buildEvalContext)
