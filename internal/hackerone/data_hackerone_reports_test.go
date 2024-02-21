package hackerone

import (
	"context"
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/zclconf/go-cty/cty"

	"github.com/blackstork-io/fabric/internal/hackerone/client"
	client_mocks "github.com/blackstork-io/fabric/mocks/internalpkg/hackerone/client"
	"github.com/blackstork-io/fabric/plugin"
)

type ReportsDataSourceTestSuite struct {
	suite.Suite
	schema    *plugin.DataSource
	ctx       context.Context
	cli       *client_mocks.Client
	storedUsr string
	storedTkn string
}

func TestReportsDataSourceTestSuite(t *testing.T) {
	suite.Run(t, new(ReportsDataSourceTestSuite))
}

func (s *ReportsDataSourceTestSuite) SetupSuite() {
	s.schema = makeHackerOneReportsDataSchema(func(user, token string) client.Client {
		s.storedUsr = user
		s.storedTkn = token
		return s.cli
	})
	s.ctx = context.Background()
}

func (s *ReportsDataSourceTestSuite) SetupTest() {
	s.cli = &client_mocks.Client{}
}

func (s *ReportsDataSourceTestSuite) TearDownTest() {
	s.cli.AssertExpectations(s.T())
}

func (s *ReportsDataSourceTestSuite) TestSchema() {
	s.Require().NotNil(s.schema)
	s.NotNil(s.schema.Config)
	s.NotNil(s.schema.Args)
	s.NotNil(s.schema.DataFunc)
}

func (s *ReportsDataSourceTestSuite) TestPageNumber() {
	s.cli.On("GetAllReports", mock.Anything, &client.GetAllReportsReq{
		PageNumber:    client.Int(123),
		FilterProgram: []string{"test_program"},
	}).Return(&client.GetAllReportsRes{
		Data: []any{
			map[string]any{
				"id": "1",
			},
		},
	}, nil)
	res, diags := s.schema.DataFunc(s.ctx, &plugin.RetrieveDataParams{
		Config: cty.ObjectVal(map[string]cty.Value{
			"api_username": cty.StringVal("test_user"),
			"api_token":    cty.StringVal("test_token"),
		}),
		Args: cty.ObjectVal(map[string]cty.Value{
			"size":                            cty.NullVal(cty.Number),
			"page_number":                     cty.NumberIntVal(123),
			"sort":                            cty.NullVal(cty.String),
			"program":                         cty.ListVal([]cty.Value{cty.StringVal("test_program")}),
			"inbox_ids":                       cty.NullVal(cty.List(cty.Number)),
			"reporter":                        cty.NullVal(cty.List(cty.String)),
			"assignee":                        cty.NullVal(cty.List(cty.String)),
			"state":                           cty.NullVal(cty.List(cty.String)),
			"id":                              cty.NullVal(cty.List(cty.Number)),
			"weakness_id":                     cty.NullVal(cty.List(cty.Number)),
			"severity":                        cty.NullVal(cty.List(cty.String)),
			"hacker_published":                cty.NullVal(cty.Bool),
			"created_at__gt":                  cty.NullVal(cty.String),
			"created_at__lt":                  cty.NullVal(cty.String),
			"submitted_at__gt":                cty.NullVal(cty.String),
			"submitted_at__lt":                cty.NullVal(cty.String),
			"triaged_at__gt":                  cty.NullVal(cty.String),
			"triaged_at__lt":                  cty.NullVal(cty.String),
			"triaged_at__null":                cty.NullVal(cty.Bool),
			"closed_at__gt":                   cty.NullVal(cty.String),
			"closed_at__lt":                   cty.NullVal(cty.String),
			"closed_at__null":                 cty.NullVal(cty.Bool),
			"disclosed_at__gt":                cty.NullVal(cty.String),
			"disclosed_at__lt":                cty.NullVal(cty.String),
			"disclosed_at__null":              cty.NullVal(cty.Bool),
			"reporter_agreed_on_going_public": cty.NullVal(cty.Bool),
			"bounty_awarded_at__gt":           cty.NullVal(cty.String),
			"bounty_awarded_at__lt":           cty.NullVal(cty.String),
			"bounty_awarded_at__null":         cty.NullVal(cty.Bool),
			"swag_awarded_at__gt":             cty.NullVal(cty.String),
			"swag_awarded_at__lt":             cty.NullVal(cty.String),
			"swag_awarded_at__null":           cty.NullVal(cty.Bool),
			"last_report_activity_at__gt":     cty.NullVal(cty.String),
			"last_report_activity_at__lt":     cty.NullVal(cty.String),
			"first_program_activity_at__gt":   cty.NullVal(cty.String),
			"first_program_activity_at__lt":   cty.NullVal(cty.String),
			"first_program_activity_at__null": cty.NullVal(cty.Bool),
			"last_program_activity_at__gt":    cty.NullVal(cty.String),
			"last_program_activity_at__lt":    cty.NullVal(cty.String),
			"last_program_activity_at__null":  cty.NullVal(cty.Bool),
			"last_activity_at__gt":            cty.NullVal(cty.String),
			"last_activity_at__lt":            cty.NullVal(cty.String),
			"last_public_activity_at__gt":     cty.NullVal(cty.String),
			"last_public_activity_at__lt":     cty.NullVal(cty.String),
			"keyword":                         cty.NullVal(cty.String),
			"custom_fields":                   cty.NullVal(cty.Map(cty.String)),
		}),
	})
	s.Equal("test_user", s.storedUsr)
	s.Equal("test_token", s.storedTkn)
	s.Len(diags, 0)
	s.Equal(plugin.ListData{
		plugin.MapData{
			"id": plugin.StringData("1"),
		},
	}, res)
}

func (s *ReportsDataSourceTestSuite) TestProgram() {
	s.cli.On("GetAllReports", mock.Anything, &client.GetAllReportsReq{
		PageNumber:    client.Int(1),
		FilterProgram: []string{"test_program"},
	}).Return(&client.GetAllReportsRes{
		Data: []any{
			map[string]any{
				"id": "1",
			},
		},
	}, nil)
	s.cli.On("GetAllReports", mock.Anything, &client.GetAllReportsReq{
		PageNumber:    client.Int(2),
		FilterProgram: []string{"test_program"},
	}).Return(&client.GetAllReportsRes{
		Data: []any{},
	}, nil)
	res, diags := s.schema.DataFunc(s.ctx, &plugin.RetrieveDataParams{
		Config: cty.ObjectVal(map[string]cty.Value{
			"api_username": cty.StringVal("test_user"),
			"api_token":    cty.StringVal("test_token"),
		}),
		Args: cty.ObjectVal(map[string]cty.Value{
			"size":                            cty.NullVal(cty.Number),
			"page_number":                     cty.NullVal(cty.Number),
			"sort":                            cty.NullVal(cty.String),
			"program":                         cty.ListVal([]cty.Value{cty.StringVal("test_program")}),
			"inbox_ids":                       cty.NullVal(cty.List(cty.Number)),
			"reporter":                        cty.NullVal(cty.List(cty.String)),
			"assignee":                        cty.NullVal(cty.List(cty.String)),
			"state":                           cty.NullVal(cty.List(cty.String)),
			"id":                              cty.NullVal(cty.List(cty.Number)),
			"weakness_id":                     cty.NullVal(cty.List(cty.Number)),
			"severity":                        cty.NullVal(cty.List(cty.String)),
			"hacker_published":                cty.NullVal(cty.Bool),
			"created_at__gt":                  cty.NullVal(cty.String),
			"created_at__lt":                  cty.NullVal(cty.String),
			"submitted_at__gt":                cty.NullVal(cty.String),
			"submitted_at__lt":                cty.NullVal(cty.String),
			"triaged_at__gt":                  cty.NullVal(cty.String),
			"triaged_at__lt":                  cty.NullVal(cty.String),
			"triaged_at__null":                cty.NullVal(cty.Bool),
			"closed_at__gt":                   cty.NullVal(cty.String),
			"closed_at__lt":                   cty.NullVal(cty.String),
			"closed_at__null":                 cty.NullVal(cty.Bool),
			"disclosed_at__gt":                cty.NullVal(cty.String),
			"disclosed_at__lt":                cty.NullVal(cty.String),
			"disclosed_at__null":              cty.NullVal(cty.Bool),
			"reporter_agreed_on_going_public": cty.NullVal(cty.Bool),
			"bounty_awarded_at__gt":           cty.NullVal(cty.String),
			"bounty_awarded_at__lt":           cty.NullVal(cty.String),
			"bounty_awarded_at__null":         cty.NullVal(cty.Bool),
			"swag_awarded_at__gt":             cty.NullVal(cty.String),
			"swag_awarded_at__lt":             cty.NullVal(cty.String),
			"swag_awarded_at__null":           cty.NullVal(cty.Bool),
			"last_report_activity_at__gt":     cty.NullVal(cty.String),
			"last_report_activity_at__lt":     cty.NullVal(cty.String),
			"first_program_activity_at__gt":   cty.NullVal(cty.String),
			"first_program_activity_at__lt":   cty.NullVal(cty.String),
			"first_program_activity_at__null": cty.NullVal(cty.Bool),
			"last_program_activity_at__gt":    cty.NullVal(cty.String),
			"last_program_activity_at__lt":    cty.NullVal(cty.String),
			"last_program_activity_at__null":  cty.NullVal(cty.Bool),
			"last_activity_at__gt":            cty.NullVal(cty.String),
			"last_activity_at__lt":            cty.NullVal(cty.String),
			"last_public_activity_at__gt":     cty.NullVal(cty.String),
			"last_public_activity_at__lt":     cty.NullVal(cty.String),
			"keyword":                         cty.NullVal(cty.String),
			"custom_fields":                   cty.NullVal(cty.Map(cty.String)),
		}),
	})
	s.Equal("test_user", s.storedUsr)
	s.Equal("test_token", s.storedTkn)
	s.Len(diags, 0)
	s.Equal(plugin.ListData{
		plugin.MapData{
			"id": plugin.StringData("1"),
		},
	}, res)
}

func (s *ReportsDataSourceTestSuite) TestInboxIDs() {
	s.cli.On("GetAllReports", mock.Anything, &client.GetAllReportsReq{
		PageNumber:     client.Int(1),
		FilterInboxIDs: []int{1, 2, 3},
	}).Return(&client.GetAllReportsRes{
		Data: []any{
			map[string]any{
				"id": "1",
			},
		},
	}, nil)
	s.cli.On("GetAllReports", mock.Anything, &client.GetAllReportsReq{
		PageNumber:     client.Int(2),
		FilterInboxIDs: []int{1, 2, 3},
	}).Return(&client.GetAllReportsRes{
		Data: []any{},
	}, nil)
	res, diags := s.schema.DataFunc(s.ctx, &plugin.RetrieveDataParams{
		Config: cty.ObjectVal(map[string]cty.Value{
			"api_username": cty.StringVal("test_user"),
			"api_token":    cty.StringVal("test_token"),
		}),
		Args: cty.ObjectVal(map[string]cty.Value{
			"size":                            cty.NullVal(cty.Number),
			"page_number":                     cty.NullVal(cty.Number),
			"sort":                            cty.NullVal(cty.String),
			"program":                         cty.NullVal(cty.List(cty.String)),
			"inbox_ids":                       cty.ListVal([]cty.Value{cty.NumberIntVal(1), cty.NumberIntVal(2), cty.NumberIntVal(3)}),
			"reporter":                        cty.NullVal(cty.List(cty.String)),
			"assignee":                        cty.NullVal(cty.List(cty.String)),
			"state":                           cty.NullVal(cty.List(cty.String)),
			"id":                              cty.NullVal(cty.List(cty.Number)),
			"weakness_id":                     cty.NullVal(cty.List(cty.Number)),
			"severity":                        cty.NullVal(cty.List(cty.String)),
			"hacker_published":                cty.NullVal(cty.Bool),
			"created_at__gt":                  cty.NullVal(cty.String),
			"created_at__lt":                  cty.NullVal(cty.String),
			"submitted_at__gt":                cty.NullVal(cty.String),
			"submitted_at__lt":                cty.NullVal(cty.String),
			"triaged_at__gt":                  cty.NullVal(cty.String),
			"triaged_at__lt":                  cty.NullVal(cty.String),
			"triaged_at__null":                cty.NullVal(cty.Bool),
			"closed_at__gt":                   cty.NullVal(cty.String),
			"closed_at__lt":                   cty.NullVal(cty.String),
			"closed_at__null":                 cty.NullVal(cty.Bool),
			"disclosed_at__gt":                cty.NullVal(cty.String),
			"disclosed_at__lt":                cty.NullVal(cty.String),
			"disclosed_at__null":              cty.NullVal(cty.Bool),
			"reporter_agreed_on_going_public": cty.NullVal(cty.Bool),
			"bounty_awarded_at__gt":           cty.NullVal(cty.String),
			"bounty_awarded_at__lt":           cty.NullVal(cty.String),
			"bounty_awarded_at__null":         cty.NullVal(cty.Bool),
			"swag_awarded_at__gt":             cty.NullVal(cty.String),
			"swag_awarded_at__lt":             cty.NullVal(cty.String),
			"swag_awarded_at__null":           cty.NullVal(cty.Bool),
			"last_report_activity_at__gt":     cty.NullVal(cty.String),
			"last_report_activity_at__lt":     cty.NullVal(cty.String),
			"first_program_activity_at__gt":   cty.NullVal(cty.String),
			"first_program_activity_at__lt":   cty.NullVal(cty.String),
			"first_program_activity_at__null": cty.NullVal(cty.Bool),
			"last_program_activity_at__gt":    cty.NullVal(cty.String),
			"last_program_activity_at__lt":    cty.NullVal(cty.String),
			"last_program_activity_at__null":  cty.NullVal(cty.Bool),
			"last_activity_at__gt":            cty.NullVal(cty.String),
			"last_activity_at__lt":            cty.NullVal(cty.String),
			"last_public_activity_at__gt":     cty.NullVal(cty.String),
			"last_public_activity_at__lt":     cty.NullVal(cty.String),
			"keyword":                         cty.NullVal(cty.String),
			"custom_fields":                   cty.NullVal(cty.Map(cty.String)),
		}),
	})
	s.Equal("test_user", s.storedUsr)
	s.Equal("test_token", s.storedTkn)
	s.Len(diags, 0)
	s.Equal(plugin.ListData{
		plugin.MapData{
			"id": plugin.StringData("1"),
		},
	}, res)
}

func (s *ReportsDataSourceTestSuite) TestInvalid() {
	res, diags := s.schema.DataFunc(s.ctx, &plugin.RetrieveDataParams{
		Config: cty.ObjectVal(map[string]cty.Value{
			"api_username": cty.StringVal("test_user"),
			"api_token":    cty.StringVal("test_token"),
		}),
		Args: cty.ObjectVal(map[string]cty.Value{
			"size":                            cty.NullVal(cty.Number),
			"page_number":                     cty.NullVal(cty.Number),
			"sort":                            cty.NullVal(cty.String),
			"program":                         cty.NullVal(cty.List(cty.String)),
			"inbox_ids":                       cty.NullVal(cty.List(cty.Number)),
			"reporter":                        cty.NullVal(cty.List(cty.String)),
			"assignee":                        cty.NullVal(cty.List(cty.String)),
			"state":                           cty.NullVal(cty.List(cty.String)),
			"id":                              cty.NullVal(cty.List(cty.Number)),
			"weakness_id":                     cty.NullVal(cty.List(cty.Number)),
			"severity":                        cty.NullVal(cty.List(cty.String)),
			"hacker_published":                cty.NullVal(cty.Bool),
			"created_at__gt":                  cty.NullVal(cty.String),
			"created_at__lt":                  cty.NullVal(cty.String),
			"submitted_at__gt":                cty.NullVal(cty.String),
			"submitted_at__lt":                cty.NullVal(cty.String),
			"triaged_at__gt":                  cty.NullVal(cty.String),
			"triaged_at__lt":                  cty.NullVal(cty.String),
			"triaged_at__null":                cty.NullVal(cty.Bool),
			"closed_at__gt":                   cty.NullVal(cty.String),
			"closed_at__lt":                   cty.NullVal(cty.String),
			"closed_at__null":                 cty.NullVal(cty.Bool),
			"disclosed_at__gt":                cty.NullVal(cty.String),
			"disclosed_at__lt":                cty.NullVal(cty.String),
			"disclosed_at__null":              cty.NullVal(cty.Bool),
			"reporter_agreed_on_going_public": cty.NullVal(cty.Bool),
			"bounty_awarded_at__gt":           cty.NullVal(cty.String),
			"bounty_awarded_at__lt":           cty.NullVal(cty.String),
			"bounty_awarded_at__null":         cty.NullVal(cty.Bool),
			"swag_awarded_at__gt":             cty.NullVal(cty.String),
			"swag_awarded_at__lt":             cty.NullVal(cty.String),
			"swag_awarded_at__null":           cty.NullVal(cty.Bool),
			"last_report_activity_at__gt":     cty.NullVal(cty.String),
			"last_report_activity_at__lt":     cty.NullVal(cty.String),
			"first_program_activity_at__gt":   cty.NullVal(cty.String),
			"first_program_activity_at__lt":   cty.NullVal(cty.String),
			"first_program_activity_at__null": cty.NullVal(cty.Bool),
			"last_program_activity_at__gt":    cty.NullVal(cty.String),
			"last_program_activity_at__lt":    cty.NullVal(cty.String),
			"last_program_activity_at__null":  cty.NullVal(cty.Bool),
			"last_activity_at__gt":            cty.NullVal(cty.String),
			"last_activity_at__lt":            cty.NullVal(cty.String),
			"last_public_activity_at__gt":     cty.NullVal(cty.String),
			"last_public_activity_at__lt":     cty.NullVal(cty.String),
			"keyword":                         cty.NullVal(cty.String),
			"custom_fields":                   cty.NullVal(cty.Map(cty.String)),
		}),
	})
	s.Nil(res)
	s.Equal(hcl.Diagnostics{{
		Severity: hcl.DiagError,
		Summary:  "Failed to parse arguments",
		Detail:   "at least one of program or inbox_ids must be provided",
	}}, diags)
}

func (s *ReportsDataSourceTestSuite) TestInvalidConfig() {
	res, diags := s.schema.DataFunc(s.ctx, &plugin.RetrieveDataParams{
		Config: cty.ObjectVal(map[string]cty.Value{
			"api_username": cty.NullVal(cty.String),
			"api_token":    cty.NullVal(cty.String),
		}),
		Args: cty.ObjectVal(map[string]cty.Value{
			"size":                            cty.NullVal(cty.Number),
			"page_number":                     cty.NullVal(cty.Number),
			"sort":                            cty.NullVal(cty.String),
			"program":                         cty.NullVal(cty.List(cty.String)),
			"inbox_ids":                       cty.NullVal(cty.List(cty.Number)),
			"reporter":                        cty.NullVal(cty.List(cty.String)),
			"assignee":                        cty.NullVal(cty.List(cty.String)),
			"state":                           cty.NullVal(cty.List(cty.String)),
			"id":                              cty.NullVal(cty.List(cty.Number)),
			"weakness_id":                     cty.NullVal(cty.List(cty.Number)),
			"severity":                        cty.NullVal(cty.List(cty.String)),
			"hacker_published":                cty.NullVal(cty.Bool),
			"created_at__gt":                  cty.NullVal(cty.String),
			"created_at__lt":                  cty.NullVal(cty.String),
			"submitted_at__gt":                cty.NullVal(cty.String),
			"submitted_at__lt":                cty.NullVal(cty.String),
			"triaged_at__gt":                  cty.NullVal(cty.String),
			"triaged_at__lt":                  cty.NullVal(cty.String),
			"triaged_at__null":                cty.NullVal(cty.Bool),
			"closed_at__gt":                   cty.NullVal(cty.String),
			"closed_at__lt":                   cty.NullVal(cty.String),
			"closed_at__null":                 cty.NullVal(cty.Bool),
			"disclosed_at__gt":                cty.NullVal(cty.String),
			"disclosed_at__lt":                cty.NullVal(cty.String),
			"disclosed_at__null":              cty.NullVal(cty.Bool),
			"reporter_agreed_on_going_public": cty.NullVal(cty.Bool),
			"bounty_awarded_at__gt":           cty.NullVal(cty.String),
			"bounty_awarded_at__lt":           cty.NullVal(cty.String),
			"bounty_awarded_at__null":         cty.NullVal(cty.Bool),
			"swag_awarded_at__gt":             cty.NullVal(cty.String),
			"swag_awarded_at__lt":             cty.NullVal(cty.String),
			"swag_awarded_at__null":           cty.NullVal(cty.Bool),
			"last_report_activity_at__gt":     cty.NullVal(cty.String),
			"last_report_activity_at__lt":     cty.NullVal(cty.String),
			"first_program_activity_at__gt":   cty.NullVal(cty.String),
			"first_program_activity_at__lt":   cty.NullVal(cty.String),
			"first_program_activity_at__null": cty.NullVal(cty.Bool),
			"last_program_activity_at__gt":    cty.NullVal(cty.String),
			"last_program_activity_at__lt":    cty.NullVal(cty.String),
			"last_program_activity_at__null":  cty.NullVal(cty.Bool),
			"last_activity_at__gt":            cty.NullVal(cty.String),
			"last_activity_at__lt":            cty.NullVal(cty.String),
			"last_public_activity_at__gt":     cty.NullVal(cty.String),
			"last_public_activity_at__lt":     cty.NullVal(cty.String),
			"keyword":                         cty.NullVal(cty.String),
			"custom_fields":                   cty.NullVal(cty.Map(cty.String)),
		}),
	})
	s.Nil(res)
	s.Equal(hcl.Diagnostics{{
		Severity: hcl.DiagError,
		Summary:  "Failed to create client",
		Detail:   "api_username is required in configuration",
	}}, diags)
}
