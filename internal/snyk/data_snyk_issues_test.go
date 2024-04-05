package snyk

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/zclconf/go-cty/cty"

	"github.com/blackstork-io/fabric/internal/snyk/client"
	client_mocks "github.com/blackstork-io/fabric/mocks/internalpkg/snyk/client"
	"github.com/blackstork-io/fabric/plugin"
)

type IssuesDataSourceTestSuite struct {
	suite.Suite
	schema       *plugin.DataSource
	ctx          context.Context
	cli          *client_mocks.Client
	storedAPIKey string
}

func TestIssuesDataSourceTestSuite(t *testing.T) {
	suite.Run(t, new(IssuesDataSourceTestSuite))
}

func (s *IssuesDataSourceTestSuite) SetupSuite() {
	s.schema = makeSnykIssuesDataSource(func(apiKey string) client.Client {
		s.storedAPIKey = apiKey
		return s.cli
	})
	s.ctx = context.Background()
}

func (s *IssuesDataSourceTestSuite) SetupTest() {
	s.cli = &client_mocks.Client{}
}

func (s *IssuesDataSourceTestSuite) TearDownTest() {
	s.cli.AssertExpectations(s.T())
}

func (s *IssuesDataSourceTestSuite) TestSchema() {
	s.Require().NotNil(s.schema)
	s.NotNil(s.schema.Config)
	s.NotNil(s.schema.Args)
	s.NotNil(s.schema.DataFunc)
}

func (s *IssuesDataSourceTestSuite) TestPaging() {
	s.cli.On("ListIssues", mock.Anything, &client.ListIssuesReq{
		GroupID: client.String("test_group_id"),
		Limit:   pageSize,
	}).Return(&client.ListIssuesRes{
		Data: []any{
			map[string]any{
				"id": "1",
			},
		},
		Links: &client.Links{
			Next: client.String("2"),
		},
	}, nil)
	s.cli.On("ListIssues", mock.Anything, &client.ListIssuesReq{
		GroupID:       client.String("test_group_id"),
		Limit:         pageSize,
		StartingAfter: client.String("2"),
	}).Return(&client.ListIssuesRes{
		Data: []any{},
	}, nil)
	res, diags := s.schema.DataFunc(s.ctx, &plugin.RetrieveDataParams{
		Config: cty.ObjectVal(map[string]cty.Value{
			"api_key": cty.StringVal("test_api_key"),
		}),
		Args: cty.ObjectVal(map[string]cty.Value{
			"group_id":                 cty.StringVal("test_group_id"),
			"org_id":                   cty.NullVal(cty.String),
			"scan_item_id":             cty.NullVal(cty.String),
			"scan_item_type":           cty.NullVal(cty.String),
			"type":                     cty.NullVal(cty.String),
			"updated_before":           cty.NullVal(cty.String),
			"updated_after":            cty.NullVal(cty.String),
			"created_before":           cty.NullVal(cty.String),
			"created_after":            cty.NullVal(cty.String),
			"effective_severity_level": cty.NullVal(cty.List(cty.String)),
			"status":                   cty.NullVal(cty.List(cty.String)),
			"ignored":                  cty.NullVal(cty.Bool),
			"limit":                    cty.NullVal(cty.Number),
		}),
	})
	s.Equal("test_api_key", s.storedAPIKey)
	s.Len(diags, 0)
	s.Equal(plugin.ListData{
		plugin.MapData{
			"id": plugin.StringData("1"),
		},
	}, res)
}

func (s *IssuesDataSourceTestSuite) TestFull() {
	s.cli.On("ListIssues", mock.Anything, &client.ListIssuesReq{
		OrgID:                  client.String("test_org_id"),
		Limit:                  pageSize,
		ScanItemID:             client.String("test_scan_item_id"),
		ScanItemType:           client.String("test_scan_item_type"),
		Type:                   client.String("test_type"),
		UpdatedBefore:          client.String("test_updated_before"),
		UpdatedAfter:           client.String("test_updated_after"),
		CreatedBefore:          client.String("test_created_before"),
		CreatedAfter:           client.String("test_created_after"),
		EffectiveSeverityLevel: client.StringList([]string{"test_effective_severity_level_1", "test_effective_severity_level_2"}),
		Status:                 client.StringList([]string{"test_status_1", "test_status_2"}),
		Ignored:                client.Bool(true),
	}).Return(&client.ListIssuesRes{
		Data: []any{
			map[string]any{
				"id": "1",
			},
			map[string]any{
				"id": "2",
			},
			map[string]any{
				"id": "3",
			},
			map[string]any{
				"id": "4",
			},
		},
	}, nil)
	res, diags := s.schema.DataFunc(s.ctx, &plugin.RetrieveDataParams{
		Config: cty.ObjectVal(map[string]cty.Value{
			"api_key": cty.StringVal("test_api_key"),
		}),
		Args: cty.ObjectVal(map[string]cty.Value{
			"group_id":                 cty.NullVal(cty.String),
			"org_id":                   cty.StringVal("test_org_id"),
			"limit":                    cty.NumberIntVal(2),
			"starting_after":           cty.StringVal("test_starting_after"),
			"scan_item_id":             cty.StringVal("test_scan_item_id"),
			"scan_item_type":           cty.StringVal("test_scan_item_type"),
			"type":                     cty.StringVal("test_type"),
			"updated_before":           cty.StringVal("test_updated_before"),
			"updated_after":            cty.StringVal("test_updated_after"),
			"created_before":           cty.StringVal("test_created_before"),
			"created_after":            cty.StringVal("test_created_after"),
			"effective_severity_level": cty.ListVal([]cty.Value{cty.StringVal("test_effective_severity_level_1"), cty.StringVal("test_effective_severity_level_2")}),
			"status":                   cty.ListVal([]cty.Value{cty.StringVal("test_status_1"), cty.StringVal("test_status_2")}),
			"ignored":                  cty.BoolVal(true),
		}),
	})
	s.Equal("test_api_key", s.storedAPIKey)
	s.Len(diags, 0)
	s.Equal(plugin.ListData{
		plugin.MapData{
			"id": plugin.StringData("1"),
		},
		plugin.MapData{
			"id": plugin.StringData("2"),
		},
	}, res)
}

func (s *IssuesDataSourceTestSuite) TestConstraintNoGroupAndOrgID() {
	_, diags := s.schema.DataFunc(s.ctx, &plugin.RetrieveDataParams{
		Config: cty.ObjectVal(map[string]cty.Value{
			"api_key": cty.StringVal("test_api"),
		}),
		Args: cty.ObjectVal(map[string]cty.Value{
			"group_id":                 cty.NullVal(cty.String),
			"org_id":                   cty.NullVal(cty.String),
			"scan_item_id":             cty.NullVal(cty.String),
			"scan_item_type":           cty.NullVal(cty.String),
			"type":                     cty.NullVal(cty.String),
			"updated_before":           cty.NullVal(cty.String),
			"updated_after":            cty.NullVal(cty.String),
			"created_before":           cty.NullVal(cty.String),
			"created_after":            cty.NullVal(cty.String),
			"effective_severity_level": cty.NullVal(cty.List(cty.String)),
			"status":                   cty.NullVal(cty.List(cty.String)),
			"ignored":                  cty.NullVal(cty.Bool),
			"limit":                    cty.NullVal(cty.Number),
		}),
	})
	s.Require().Len(diags, 1)
	s.Equal("Failed to create Snyk request", diags[0].Summary)
	s.Equal("either group_id or org_id must be set", diags[0].Detail)
}

func (s *IssuesDataSourceTestSuite) TestConstraintBothGroupAndOrgID() {
	_, diags := s.schema.DataFunc(s.ctx, &plugin.RetrieveDataParams{
		Config: cty.ObjectVal(map[string]cty.Value{
			"api_key": cty.StringVal("test_api"),
		}),
		Args: cty.ObjectVal(map[string]cty.Value{
			"group_id":                 cty.StringVal("test_group_id"),
			"org_id":                   cty.StringVal("test_org_id"),
			"scan_item_id":             cty.NullVal(cty.String),
			"scan_item_type":           cty.NullVal(cty.String),
			"type":                     cty.NullVal(cty.String),
			"updated_before":           cty.NullVal(cty.String),
			"updated_after":            cty.NullVal(cty.String),
			"created_before":           cty.NullVal(cty.String),
			"created_after":            cty.NullVal(cty.String),
			"effective_severity_level": cty.NullVal(cty.List(cty.String)),
			"status":                   cty.NullVal(cty.List(cty.String)),
			"ignored":                  cty.NullVal(cty.Bool),
			"limit":                    cty.NullVal(cty.Number),
		}),
	})
	s.Require().Len(diags, 1)
	s.Equal("Failed to create Snyk request", diags[0].Summary)
	s.Equal("only one of group_id or org_id is allowed", diags[0].Detail)
}
