package evaluation

import (
	"log/slog"
	"os"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/joho/godotenv"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"
	"github.com/zclconf/go-cty/cty/function/stdlib"
)

func buildEnvVarMap() cty.Value {
	envVars := os.Environ()
	envFromFile, err := godotenv.Read()
	if err != nil && !os.IsNotExist(err) {
		slog.Error("Error reading .env file", "err", err)
	}
	envMap := make(map[string]cty.Value, len(envFromFile)+len(envVars))
	for k, v := range envFromFile {
		envMap[k] = cty.StringVal(v)
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

var fromFileFunc = function.New(&function.Spec{
	Description: "Reads the content of a file and returns it as a string",
	Params: []function.Parameter{{
		Name:        "path",
		Description: "The path to the file to read",
		Type:        cty.String,
	}},
	Type: function.StaticReturnType(cty.String),
	Impl: func(args []cty.Value, retType cty.Type) (cty.Value, error) {
		res, err := os.ReadFile(args[0].AsString())
		if err != nil {
			return cty.NullVal(cty.String), err
		}
		return cty.StringVal(string(res)), nil
	},
})

func EvalContext() *hcl.EvalContext {
	return &hcl.EvalContext{
		Variables: map[string]cty.Value{
			"env": buildEnvVarMap(),
		},
		Functions: map[string]function.Function{
			"from_file": fromFileFunc,
			"join":      stdlib.JoinFunc,
		},
	}
}
