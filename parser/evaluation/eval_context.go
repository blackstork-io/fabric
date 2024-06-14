package evaluation

import (
	"log/slog"
	"os"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/joho/godotenv"
	"github.com/zclconf/go-cty/cty"
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

func EvalContext() *hcl.EvalContext {
	return &hcl.EvalContext{
		Variables: map[string]cty.Value{
			"env": buildEnvVarMap(),
		},
	}
}
