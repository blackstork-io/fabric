package postgresql

import (
	"context"
	"database/sql"
	"path/filepath"
	"testing"
	"time"

	"github.com/blackstork-io/fabric/plugininterface/v1"
	"github.com/hashicorp/hcl/v2"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"github.com/zclconf/go-cty/cty"
)

// IntegrationTestSuite is a test suite to test integration with real postgres instance
type IntegrationTestSuite struct {
	suite.Suite
	container *postgres.PostgresContainer
	connURL   string
	plugin    Plugin
	ctx       context.Context
}

func TestIntegrationSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration tests")
	}
	suite.Run(t, &IntegrationTestSuite{})
}

func (s *IntegrationTestSuite) SetupSuite() {
	s.ctx = context.Background()
	opts := []testcontainers.ContainerCustomizer{
		testcontainers.WithImage("docker.io/postgres:15.2-alpine"),
		postgres.WithInitScripts(filepath.Join("testdata", "data.sql")),
		postgres.WithDatabase("testusr123"),
		postgres.WithPassword("testpsw123"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5 * time.Second)),
	}
	container, err := postgres.RunContainer(s.ctx, opts...)
	s.Require().NoError(err, "failed to start postgres container")
	s.container = container
	connURL, err := container.ConnectionString(s.ctx, "sslmode=disable")
	s.Require().NoError(err, "failed to get postgres connection string")
	s.connURL = connURL
	db, err := sql.Open("postgres", connURL)
	s.Require().NoError(err, "failed to open postgres database")
	err = db.Ping()
	s.Require().NoError(err, "failed to ping postgres database")
	err = db.Close()
	s.Require().NoError(err, "failed to close postgres database")
}

func (s *IntegrationTestSuite) TearDownSuite() {
	s.Require().NoError(s.container.Terminate(s.ctx), "failed to stop postgres container")
}

func (s *IntegrationTestSuite) TestEmptyDatabaseURL() {
	res := s.plugin.Call(plugininterface.Args{
		Kind: "data",
		Name: "postgresql",
		Config: cty.ObjectVal(map[string]cty.Value{
			"database_url": cty.StringVal(""),
		}),
	})
	s.Require().Equal(res, plugininterface.Result{
		Diags: hcl.Diagnostics{
			{
				Severity: hcl.DiagError,
				Summary:  "Invalid configuration",
				Detail:   "database_url is required",
			},
		},
	})
}
func (s *IntegrationTestSuite) TestNilDatabaseURL() {
	res := s.plugin.Call(plugininterface.Args{
		Kind: "data",
		Name: "postgresql",
		Config: cty.ObjectVal(map[string]cty.Value{
			"database_url": cty.NullVal(cty.String),
		}),
	})
	s.Require().Equal(res, plugininterface.Result{
		Diags: hcl.Diagnostics{
			{
				Severity: hcl.DiagError,
				Summary:  "Invalid configuration",
				Detail:   "database_url is required",
			},
		},
	})
}

func (s *IntegrationTestSuite) TestEmptySQLQuery() {
	res := s.plugin.Call(plugininterface.Args{
		Kind: "data",
		Name: "postgresql",
		Config: cty.ObjectVal(map[string]cty.Value{
			"database_url": cty.StringVal(s.connURL),
		}),
		Args: cty.ObjectVal(map[string]cty.Value{
			"sql_query": cty.StringVal(""),
			"sql_args":  cty.ListValEmpty(cty.DynamicPseudoType),
		}),
	})
	s.Require().Equal(res, plugininterface.Result{
		Diags: hcl.Diagnostics{
			{
				Severity: hcl.DiagError,
				Summary:  "Invalid arguments",
				Detail:   "sql_query is required",
			},
		},
	})
}
func (s *IntegrationTestSuite) TestNilSQLQuery() {
	res := s.plugin.Call(plugininterface.Args{
		Kind: "data",
		Name: "postgresql",
		Config: cty.ObjectVal(map[string]cty.Value{
			"database_url": cty.StringVal(s.connURL),
		}),
		Args: cty.ObjectVal(map[string]cty.Value{
			"sql_query": cty.NullVal(cty.String),
			"sql_args":  cty.ListValEmpty(cty.DynamicPseudoType),
		}),
	})
	s.Require().Equal(res, plugininterface.Result{
		Diags: hcl.Diagnostics{
			{
				Severity: hcl.DiagError,
				Summary:  "Invalid arguments",
				Detail:   "sql_query is required",
			},
		},
	})
}

func (s *IntegrationTestSuite) TestSelectEmptyTable() {
	res := s.plugin.Call(plugininterface.Args{
		Kind: "data",
		Name: "postgresql",
		Config: cty.ObjectVal(map[string]cty.Value{
			"database_url": cty.StringVal(s.connURL),
		}),
		Args: cty.ObjectVal(map[string]cty.Value{
			"sql_query": cty.StringVal(`SELECT * FROM testdata_empty`),
			"sql_args":  cty.ListValEmpty(cty.DynamicPseudoType),
		}),
	})
	s.Require().Equal(res, plugininterface.Result{
		Result: []map[string]any{},
	})
}

func (s *IntegrationTestSuite) TestSelect() {
	res := s.plugin.Call(plugininterface.Args{
		Kind: "data",
		Name: "postgresql",
		Config: cty.ObjectVal(map[string]cty.Value{
			"database_url": cty.StringVal(s.connURL),
		}),
		Args: cty.ObjectVal(map[string]cty.Value{
			"sql_query": cty.StringVal(`SELECT * FROM testdata`),
			"sql_args":  cty.ListValEmpty(cty.DynamicPseudoType),
		}),
	})
	s.Require().Equal(res, plugininterface.Result{
		Result: []map[string]any{
			{
				"id":       int64(1),
				"text_val": "text_1",
				"int_val":  int64(1),
				"bool_val": true,
			},
			{
				"id":       int64(2),
				"text_val": "text_2",
				"int_val":  int64(2),
				"bool_val": false,
			},
		},
	})
}

func (s *IntegrationTestSuite) TestSelectSomeFields() {
	res := s.plugin.Call(plugininterface.Args{
		Kind: "data",
		Name: "postgresql",
		Config: cty.ObjectVal(map[string]cty.Value{
			"database_url": cty.StringVal(s.connURL),
		}),
		Args: cty.ObjectVal(map[string]cty.Value{
			"sql_query": cty.StringVal(`SELECT id, text_val AS text FROM testdata`),
			"sql_args":  cty.ListValEmpty(cty.DynamicPseudoType),
		}),
	})
	s.Require().Equal(res, plugininterface.Result{
		Result: []map[string]any{
			{
				"id":   int64(1),
				"text": "text_1",
			},
			{
				"id":   int64(2),
				"text": "text_2",
			},
		},
	})
}

func (s *IntegrationTestSuite) TestSelectWithArgs() {
	res := s.plugin.Call(plugininterface.Args{
		Kind: "data",
		Name: "postgresql",
		Config: cty.ObjectVal(map[string]cty.Value{
			"database_url": cty.StringVal(s.connURL),
		}),
		Args: cty.ObjectVal(map[string]cty.Value{
			"sql_query": cty.StringVal(`SELECT * FROM testdata WHERE bool_val = $1;`),
			"sql_args": cty.ListVal([]cty.Value{
				cty.BoolVal(false),
			}),
		}),
	})
	s.Require().Equal(res, plugininterface.Result{
		Result: []map[string]any{
			{
				"id":       int64(2),
				"text_val": "text_2",
				"int_val":  int64(2),
				"bool_val": false,
			},
		},
	})
}

func (s *IntegrationTestSuite) TestSelectWithMissingArgs() {
	res := s.plugin.Call(plugininterface.Args{
		Kind: "data",
		Name: "postgresql",
		Config: cty.ObjectVal(map[string]cty.Value{
			"database_url": cty.StringVal(s.connURL),
		}),
		Args: cty.ObjectVal(map[string]cty.Value{
			"sql_query": cty.StringVal(`SELECT * FROM testdata WHERE bool_val = $1;`),
			"sql_args":  cty.NullVal(cty.List(cty.DynamicPseudoType)),
		}),
	})
	s.Require().Equal(res, plugininterface.Result{
		Diags: hcl.Diagnostics{
			{
				Severity: hcl.DiagError,
				Summary:  "Failed to query database",
				Detail:   "pq: there is no parameter $1",
			},
		},
	})
}
