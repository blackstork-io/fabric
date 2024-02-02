package opencti

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/blackstork-io/fabric/plugininterface/v1"
	"github.com/hashicorp/hcl/v2"
	"github.com/stretchr/testify/suite"
	"github.com/zclconf/go-cty/cty"
)

type PluginTestSuite struct {
	suite.Suite
	ctx    context.Context
	cancel context.CancelFunc
}

func (s *PluginTestSuite) SetupTest() {
	s.ctx, s.cancel = context.WithCancel(context.Background())
}

func (s *PluginTestSuite) TearDownTest() {
	s.cancel()
}
func TestPluginTestSuite(t *testing.T) {
	suite.Run(t, new(PluginTestSuite))
}

func (s *PluginTestSuite) TestGetPlugins() {
	plugins := Plugin{}.GetPlugins()
	s.Require().Len(plugins, 1, "expected 1 plugin")
	got := plugins[0]
	s.Equal("opencti", got.Name)
	s.Equal("data", got.Kind)
	s.Equal("blackstork", got.Namespace)
	s.Equal(Version.String(), got.Version.Cast().String())
	s.NotNil(got.ConfigSpec)
	s.NotNil(got.InvocationSpec)
}

func (s *PluginTestSuite) TestBasicValid() {
	want := plugininterface.Result{
		Result: jsonAny(`
			{
				"data": {
					"stixCoreRelationships": {
						"edges": []
					}
				}
			}
		`),
	}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.Equal("application/json", r.Header.Get("Content-Type"))
		s.Equal("application/json", r.Header.Get("Accept"))
		body, err := io.ReadAll(r.Body)
		s.NoError(err)
		s.Equal(`{"query":"query issue { stixCoreRelationships { edges { node { x_opencti_stix_ids } } } }"}`, string(body))
		s.Equal("POST", r.Method)
		w.Write([]byte(`
			{
				"data": {
					"stixCoreRelationships": {
						"edges": []
					}
				}
			}
		`))
	}))
	defer srv.Close()
	p := Plugin{}
	result := p.Call(plugininterface.Args{
		Kind: "data",
		Name: "opencti",
		Config: cty.ObjectVal(map[string]cty.Value{
			"graphql_url": cty.StringVal(srv.URL),
			"auth_token":  cty.NullVal(cty.String),
		}),
		Args: cty.ObjectVal(map[string]cty.Value{
			"graphql_query": cty.StringVal("query issue { stixCoreRelationships { edges { node { x_opencti_stix_ids } } } }"),
		}),
	})
	s.Equal(want, result)
}

func (s *PluginTestSuite) TestFailRequest() {
	want := plugininterface.Result{
		Diags: hcl.Diagnostics{
			{
				Severity: hcl.DiagError,
				Summary:  "Failed to execute query",
				Detail:   "unexpected status code: 404",
			},
		},
	}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer srv.Close()
	p := Plugin{}
	result := p.Call(plugininterface.Args{
		Kind: "data",
		Name: "opencti",
		Config: cty.ObjectVal(map[string]cty.Value{
			"graphql_url": cty.StringVal(srv.URL),
			"auth_token":  cty.NullVal(cty.String),
		}),
		Args: cty.ObjectVal(map[string]cty.Value{
			"graphql_query": cty.StringVal("query issue { stixCoreRelationships { edges { node { x_opencti_stix_ids } } } }"),
		}),
	})
	s.Equal(want, result)
}

func (s *PluginTestSuite) TestInvalidQuery() {
	want := plugininterface.Result{
		Diags: hcl.Diagnostics{
			{
				Severity: hcl.DiagError,
				Summary:  "Invalid GraphQL query",
				Detail:   "external: field: stixCoreRelationshipsInvalid not defined on type: Query, locations: [], path: [query]",
			},
		},
	}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.Equal("application/json", r.Header.Get("Content-Type"))
		s.Equal("application/json", r.Header.Get("Accept"))
		body, err := io.ReadAll(r.Body)
		s.NoError(err)
		s.Equal(`{"query":"query issue { stixCoreRelationshipsInvalid { edges { node { x_opencti_stix_ids } } } }"}`, string(body))
		s.Equal("POST", r.Method)
		w.Write([]byte(`
			{
				"data": {
					"stixCoreRelationships": {
						"edges": []
					}
				}
			}
		`))
	}))
	defer srv.Close()
	p := Plugin{}
	result := p.Call(plugininterface.Args{
		Kind: "data",
		Name: "opencti",
		Config: cty.ObjectVal(map[string]cty.Value{
			"graphql_url": cty.StringVal(srv.URL),
			"auth_token":  cty.NullVal(cty.String),
		}),
		Args: cty.ObjectVal(map[string]cty.Value{
			"graphql_query": cty.StringVal("query issue { stixCoreRelationshipsInvalid { edges { node { x_opencti_stix_ids } } } }"),
		}),
	})
	s.Equal(want, result)
}

func (s *PluginTestSuite) TestWithAuth() {
	want := plugininterface.Result{
		Result: jsonAny(`
			{
				"data": {
					"stixCoreRelationships": {
						"edges": []
					}
				}
			}
		`),
	}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.Equal("Bearer token-123", r.Header.Get("Authorization"))
		s.Equal("application/json", r.Header.Get("Content-Type"))
		s.Equal("application/json", r.Header.Get("Accept"))
		body, err := io.ReadAll(r.Body)
		s.NoError(err)
		s.Equal(`{"query":"query issue { stixCoreRelationships { edges { node { x_opencti_stix_ids } } } }"}`, string(body))
		s.Equal("POST", r.Method)
		w.Write([]byte(`
			{
				"data": {
					"stixCoreRelationships": {
						"edges": []
					}
				}
			}
		`))
	}))
	defer srv.Close()
	p := Plugin{}
	result := p.Call(plugininterface.Args{
		Kind: "data",
		Name: "opencti",
		Config: cty.ObjectVal(map[string]cty.Value{
			"graphql_url": cty.StringVal(srv.URL),
			"auth_token":  cty.StringVal("token-123"),
		}),
		Args: cty.ObjectVal(map[string]cty.Value{
			"graphql_query": cty.StringVal("query issue { stixCoreRelationships { edges { node { x_opencti_stix_ids } } } }"),
		}),
	})
	s.Equal(want, result)
}

func (s *PluginTestSuite) TestNullURL() {
	p := Plugin{}
	result := p.Call(plugininterface.Args{
		Kind: "data",
		Name: "opencti",
		Config: cty.ObjectVal(map[string]cty.Value{
			"graphql_url": cty.NullVal(cty.String),
			"auth_token":  cty.NullVal(cty.String),
		}),
		Args: cty.ObjectVal(map[string]cty.Value{
			"graphql_query": cty.StringVal("query{user{id,name}}"),
		}),
	})
	s.Len(result.Diags, 1)
	s.Equal("Failed to parse config", result.Diags[0].Summary)
}

func (s *PluginTestSuite) TestEmptyURL() {
	p := Plugin{}
	result := p.Call(plugininterface.Args{
		Kind: "data",
		Name: "opencti",
		Config: cty.ObjectVal(map[string]cty.Value{
			"graphql_url": cty.StringVal(""),
			"auth_token":  cty.NullVal(cty.String),
		}),
		Args: cty.ObjectVal(map[string]cty.Value{
			"graphql_query": cty.StringVal("query{user{id,name}}"),
		}),
	})
	s.Len(result.Diags, 1)
	s.Equal("Failed to parse config", result.Diags[0].Summary)
}

func (s *PluginTestSuite) TestEmptyQuery() {
	p := Plugin{}
	result := p.Call(plugininterface.Args{
		Kind: "data",
		Name: "opencti",
		Config: cty.ObjectVal(map[string]cty.Value{
			"graphql_url": cty.StringVal("http://localhost"),
			"auth_token":  cty.NullVal(cty.String),
		}),
		Args: cty.ObjectVal(map[string]cty.Value{
			"graphql_query": cty.StringVal(""),
		}),
	})
	s.Len(result.Diags, 1)
	s.Equal("Failed to parse arguments", result.Diags[0].Summary)
}

func (s *PluginTestSuite) TestNullQuery() {
	p := Plugin{}
	result := p.Call(plugininterface.Args{
		Kind: "data",
		Name: "opencti",
		Config: cty.ObjectVal(map[string]cty.Value{
			"graphql_url": cty.StringVal("http://localhost"),
			"auth_token":  cty.NullVal(cty.String),
		}),
		Args: cty.ObjectVal(map[string]cty.Value{
			"graphql_query": cty.NullVal(cty.String),
		}),
	})
	s.Len(result.Diags, 1)
	s.Equal("Failed to parse arguments", result.Diags[0].Summary)
}

func jsonAny(s string) any {
	var v any
	err := json.Unmarshal([]byte(s), &v)
	if err != nil {
		panic(err)
	}
	return v
}
