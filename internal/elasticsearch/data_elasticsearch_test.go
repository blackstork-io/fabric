package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"os"
	"testing"

	es "github.com/elastic/go-elasticsearch/v8"
	"github.com/hashicorp/hcl/v2"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/elasticsearch"
	"github.com/zclconf/go-cty/cty"

	"github.com/blackstork-io/fabric/plugin"
)

// IntegrationTestSuite is a test suite to test integration with real elasticsearch
type IntegrationTestSuite struct {
	suite.Suite
	container *elasticsearch.ElasticsearchContainer
	client    *es.Client
	plugin    *plugin.Schema
	cfg       cty.Value
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
		testcontainers.WithImage("docker.io/elasticsearch:8.9.0"),
		elasticsearch.WithPassword("password123"),
	}
	container, err := elasticsearch.RunContainer(s.ctx, opts...)

	s.Require().NoError(err, "failed to start elasticsearch container")
	s.container = container
	client, err := es.NewClient(es.Config{
		Addresses: []string{
			container.Settings.Address,
		},
		Username: "elastic",
		Password: container.Settings.Password,
		CACert:   container.Settings.CACert,
	})
	s.Require().NoError(err, "failed to create elasticsearch client")
	s.client = client
	s.cfg = cty.ObjectVal(map[string]cty.Value{
		"base_url":            cty.StringVal(s.container.Settings.Address),
		"cloud_id":            cty.NullVal(cty.String),
		"api_key_str":         cty.NullVal(cty.String),
		"api_key":             cty.NullVal(cty.List(cty.String)),
		"basic_auth_username": cty.StringVal("elastic"),
		"basic_auth_password": cty.StringVal("password123"),
		"bearer_auth":         cty.NullVal(cty.String),
		"ca_certs":            cty.StringVal(string(s.container.Settings.CACert)),
	})
	s.plugin = Plugin("")
}

func (s *IntegrationTestSuite) TearDownSuite() {
	s.Require().NoError(s.container.Terminate(s.ctx), "failed to stop elasticsearch container")
}

type testDataObject struct {
	ID     string `json:"id"`
	Type   string `json:"type"`
	Active bool   `json:"active"`
	Age    int    `json:"age"`
	Name   string `json:"name"`
}

func (s *IntegrationTestSuite) SetupTest() {
	file, err := os.ReadFile("testdata/data.json")
	s.Require().NoError(err, "failed to read data.json")
	dataList := []testDataObject{}
	err = json.Unmarshal(file, &dataList)
	s.Require().NoError(err, "failed to unmarshal data.json")
	res, err := s.client.Indices.Create("test_index")
	s.Require().NoError(err, "failed to create index test_index")
	s.Require().Equal(http.StatusOK, res.StatusCode)
	for _, data := range dataList {
		dataBytes, err := json.Marshal(data)
		s.Require().NoError(err, "failed to marshal data")
		res, err := s.client.Create("test_index", data.ID, bytes.NewReader(dataBytes))
		s.Require().NoError(err, "failed to index data")
		s.Require().False(res.IsError(), "failed to index data")
		res, err = s.client.Index("test_index", bytes.NewReader(dataBytes), s.client.Index.WithDocumentID(data.ID))
		s.Require().NoError(err, "failed to index data")
		s.Require().False(res.IsError(), "failed to index data")
	}
	res, err = s.client.Indices.Refresh()
	s.Require().NoError(err, "failed to refresh indices")
	s.Require().False(res.IsError(), "failed to refresh indices")
}

func (s *IntegrationTestSuite) TearDownTest() {
	res, err := s.client.Indices.Delete([]string{"test_index"})
	s.Require().NoError(err, "failed to delete indices")
	s.Require().False(res.IsError(), "failed to delete indices: %s", res.String())
}

func (s *IntegrationTestSuite) TestSearchDefaults() {
	args := cty.ObjectVal(map[string]cty.Value{
		"id":           cty.NullVal(cty.String),
		"index":        cty.StringVal("test_index"),
		"query":        cty.NullVal(cty.DynamicPseudoType),
		"query_string": cty.NullVal(cty.String),
		"fields":       cty.NullVal(cty.String),
		"size":         cty.NullVal(cty.Number),
	})
	data, diags := s.plugin.RetrieveData(s.ctx, "elasticsearch", &plugin.RetrieveDataParams{
		Config: s.cfg,
		Args:   args,
	})
	s.Require().Nil(diags)
	m := data.(plugin.MapData)
	raw, err := json.MarshalIndent(m["hits"], "", "  ")
	s.Require().NoError(err, "failed to marshal data: %v", err)
	s.JSONEq(`{
		"hits": [
		  {
			"_id": "54f7a815-eac5-4f7c-a339-5fefd0f54967",
			"_index": "test_index",
			"_score": 1,
			"_source": {
			  "active": false,
			  "age": 39,
			  "id": "54f7a815-eac5-4f7c-a339-5fefd0f54967",
			  "name": "Davidson",
			  "type": "foo"
			}
		  },
		  {
			"_id": "0c68e63d-daaa-4a62-92e6-e855bd144fb6",
			"_index": "test_index",
			"_score": 1,
			"_source": {
			  "active": false,
			  "age": 20,
			  "id": "0c68e63d-daaa-4a62-92e6-e855bd144fb6",
			  "name": "Thompson",
			  "type": "bar"
			}
		  },
		  {
			"_id": "a117a5e6-23d0-4daa-be3c-a70900ca4163",
			"_index": "test_index",
			"_score": 1,
			"_source": {
			  "active": true,
			  "age": 21,
			  "id": "a117a5e6-23d0-4daa-be3c-a70900ca4163",
			  "name": "Armstrong",
			  "type": "foo"
			}
		  }
		],
		"max_score": 1,
		"total": {
		  "relation": "eq",
		  "value": 3
		}
	  }`, string(raw))
}

func (s *IntegrationTestSuite) TestSearchFields() {
	args := cty.ObjectVal(map[string]cty.Value{
		"id":           cty.NullVal(cty.String),
		"index":        cty.StringVal("test_index"),
		"query":        cty.NullVal(cty.DynamicPseudoType),
		"query_string": cty.NullVal(cty.String),
		"fields":       cty.ListVal([]cty.Value{cty.StringVal("name"), cty.StringVal("age")}),
		"size":         cty.NullVal(cty.Number),
	})
	data, diags := s.plugin.RetrieveData(s.ctx, "elasticsearch", &plugin.RetrieveDataParams{
		Config: s.cfg,
		Args:   args,
	})
	s.Require().Nil(diags)
	m := data.(plugin.MapData)
	raw, err := json.MarshalIndent(m["hits"], "", "  ")
	s.Require().NoError(err, "failed to marshal data: %v", err)
	s.JSONEq(`{
		"hits": [
		  {
			"_id": "54f7a815-eac5-4f7c-a339-5fefd0f54967",
			"_index": "test_index",
			"_score": 1,
			"_source": {
			  "age": 39,
			  "name": "Davidson"
			}
		  },
		  {
			"_id": "0c68e63d-daaa-4a62-92e6-e855bd144fb6",
			"_index": "test_index",
			"_score": 1,
			"_source": {
			  "age": 20,
			  "name": "Thompson"
			}
		  },
		  {
			"_id": "a117a5e6-23d0-4daa-be3c-a70900ca4163",
			"_index": "test_index",
			"_score": 1,
			"_source": {
			  "age": 21,
			  "name": "Armstrong"
			}
		  }
		],
		"max_score": 1,
		"total": {
		  "relation": "eq",
		  "value": 3
		}
	  }`, string(raw))
}

func (s *IntegrationTestSuite) TestSearchQueryString() {
	args := cty.ObjectVal(map[string]cty.Value{
		"id":           cty.NullVal(cty.String),
		"index":        cty.StringVal("test_index"),
		"query":        cty.NullVal(cty.DynamicPseudoType),
		"query_string": cty.StringVal("type:foo"),
		"fields":       cty.NullVal(cty.String),
		"size":         cty.NullVal(cty.Number),
	})
	data, diags := s.plugin.RetrieveData(s.ctx, "elasticsearch", &plugin.RetrieveDataParams{
		Config: s.cfg,
		Args:   args,
	})
	s.Require().Nil(diags)
	m := data.(plugin.MapData)
	raw, err := json.MarshalIndent(m["hits"], "", "  ")
	s.Require().NoError(err, "failed to marshal data: %v", err)
	s.JSONEq(`{
		"hits": [
		  {
			"_id": "54f7a815-eac5-4f7c-a339-5fefd0f54967",
			"_index": "test_index",
			"_score": 0.44183272,
			"_source": {
			  "active": false,
			  "age": 39,
			  "id": "54f7a815-eac5-4f7c-a339-5fefd0f54967",
			  "name": "Davidson",
			  "type": "foo"
			}
		  },
		  {
			"_id": "a117a5e6-23d0-4daa-be3c-a70900ca4163",
			"_index": "test_index",
			"_score": 0.44183272,
			"_source": {
			  "active": true,
			  "age": 21,
			  "id": "a117a5e6-23d0-4daa-be3c-a70900ca4163",
			  "name": "Armstrong",
			  "type": "foo"
			}
		  }
		],
		"max_score": 0.44183272,
		"total": {
		  "relation": "eq",
		  "value": 2
		}
	  }`, string(raw))
}

func (s *IntegrationTestSuite) TestSearchQuery() {
	args := cty.ObjectVal(map[string]cty.Value{
		"id":    cty.NullVal(cty.String),
		"index": cty.StringVal("test_index"),
		"query": cty.MapVal(map[string]cty.Value{
			"match_all": cty.MapValEmpty(cty.DynamicPseudoType),
		}),
		"query_string": cty.NullVal(cty.String),
		"fields":       cty.NullVal(cty.String),
		"size":         cty.NullVal(cty.Number),
	})
	data, diags := s.plugin.RetrieveData(s.ctx, "elasticsearch", &plugin.RetrieveDataParams{
		Config: s.cfg,
		Args:   args,
	})
	s.Require().Nil(diags)
	m := data.(plugin.MapData)
	raw, err := json.MarshalIndent(m["hits"], "", "  ")
	s.Require().NoError(err, "failed to marshal data: %v", err)
	s.JSONEq(`{
		"hits": [
			{
				"_id": "54f7a815-eac5-4f7c-a339-5fefd0f54967",
				"_index": "test_index",
				"_score": 1,
				"_source": {
					"active": false,
					"age": 39,
					"id": "54f7a815-eac5-4f7c-a339-5fefd0f54967",
					"name": "Davidson",
					"type": "foo"
				}
			},
			{
				"_id": "0c68e63d-daaa-4a62-92e6-e855bd144fb6",
				"_index": "test_index",
				"_score": 1,
				"_source": {
					"active": false,
					"age": 20,
					"id": "0c68e63d-daaa-4a62-92e6-e855bd144fb6",
					"name": "Thompson",
					"type": "bar"
				}
			},
			{
				"_id": "a117a5e6-23d0-4daa-be3c-a70900ca4163",
				"_index": "test_index",
				"_score": 1,
				"_source": {
					"active": true,
					"age": 21,
					"id": "a117a5e6-23d0-4daa-be3c-a70900ca4163",
					"name": "Armstrong",
					"type": "foo"
				}
			}
		],
		"max_score": 1,
		"total": {
			"relation": "eq",
			"value": 3
		}
	}`, string(raw))
}

func (s *IntegrationTestSuite) TestSearchSize() {
	args := cty.ObjectVal(map[string]cty.Value{
		"id":           cty.NullVal(cty.String),
		"index":        cty.StringVal("test_index"),
		"query":        cty.NullVal(cty.DynamicPseudoType),
		"query_string": cty.NullVal(cty.String),
		"fields":       cty.NullVal(cty.String),
		"size":         cty.NumberIntVal(1),
	})
	data, diags := s.plugin.RetrieveData(s.ctx, "elasticsearch", &plugin.RetrieveDataParams{
		Config: s.cfg,
		Args:   args,
	})
	s.Require().Nil(diags)
	m := data.(plugin.MapData)
	raw, err := json.MarshalIndent(m["hits"], "", "  ")
	s.Require().NoError(err, "failed to marshal data: %v", err)
	s.JSONEq(`{
		"hits": [
		  {
			"_id": "54f7a815-eac5-4f7c-a339-5fefd0f54967",
			"_index": "test_index",
			"_score": 1,
			"_source": {
			  "active": false,
			  "age": 39,
			  "id": "54f7a815-eac5-4f7c-a339-5fefd0f54967",
			  "name": "Davidson",
			  "type": "foo"
			}
		  }
		],
		"max_score": 1,
		"total": {
		  "relation": "eq",
		  "value": 3
		}
	  }`, string(raw))
}

func (s *IntegrationTestSuite) TestGetByID() {
	args := cty.ObjectVal(map[string]cty.Value{
		"id":           cty.StringVal("0c68e63d-daaa-4a62-92e6-e855bd144fb6"),
		"index":        cty.StringVal("test_index"),
		"query":        cty.NullVal(cty.DynamicPseudoType),
		"query_string": cty.NullVal(cty.String),
		"fields":       cty.NullVal(cty.String),
	})
	data, diags := s.plugin.RetrieveData(s.ctx, "elasticsearch", &plugin.RetrieveDataParams{
		Config: s.cfg,
		Args:   args,
	})
	s.Require().Nil(diags)
	m := data.(plugin.MapData)
	raw, err := json.MarshalIndent(m, "", "  ")
	s.Require().NoError(err, "failed to marshal data: %v", err)
	s.JSONEq(`{
		"_id": "0c68e63d-daaa-4a62-92e6-e855bd144fb6",
		"_index": "test_index",
		"_primary_term": 1,
		"_seq_no": 3,
		"_source": {
		  "active": false,
		  "age": 20,
		  "id": "0c68e63d-daaa-4a62-92e6-e855bd144fb6",
		  "name": "Thompson",
		  "type": "bar"
		},
		"_version": 2,
		"found": true
	}`, string(raw))
}

func (s *IntegrationTestSuite) TestGetByIDFields() {
	args := cty.ObjectVal(map[string]cty.Value{
		"id":           cty.StringVal("0c68e63d-daaa-4a62-92e6-e855bd144fb6"),
		"index":        cty.StringVal("test_index"),
		"query":        cty.NullVal(cty.DynamicPseudoType),
		"query_string": cty.NullVal(cty.String),
		"fields":       cty.ListVal([]cty.Value{cty.StringVal("name"), cty.StringVal("age")}),
	})
	data, diags := s.plugin.RetrieveData(s.ctx, "elasticsearch", &plugin.RetrieveDataParams{
		Config: s.cfg,
		Args:   args,
	})
	s.Require().Nil(diags)
	m := data.(plugin.MapData)
	raw, err := json.MarshalIndent(m, "", "  ")
	s.Require().NoError(err, "failed to marshal data: %v", err)
	s.JSONEq(`{
		"_id": "0c68e63d-daaa-4a62-92e6-e855bd144fb6",
		"_index": "test_index",
		"_primary_term": 1,
		"_seq_no": 3,
		"_source": {
		  "age": 20,
		  "name": "Thompson"
		},
		"_version": 2,
		"found": true
	}`, string(raw))
}

func (s *IntegrationTestSuite) TestGetByIDNotFound() {
	args := cty.ObjectVal(map[string]cty.Value{
		"id":           cty.StringVal("00000000-0000-0000-0000-000000000000"),
		"index":        cty.StringVal("test_index"),
		"query":        cty.NullVal(cty.DynamicPseudoType),
		"query_string": cty.NullVal(cty.String),
		"fields":       cty.NullVal(cty.String),
	})
	data, diags := s.plugin.RetrieveData(s.ctx, "elasticsearch", &plugin.RetrieveDataParams{
		Config: s.cfg,
		Args:   args,
	})
	s.Nil(data)
	s.Equal(hcl.Diagnostics{{
		Severity: hcl.DiagError,
		Summary:  "Failed to get data",
		Detail:   "failed to get document: [404 Not Found] {\"_index\":\"test_index\",\"_id\":\"00000000-0000-0000-0000-000000000000\",\"found\":false}",
	}}, diags)
}
