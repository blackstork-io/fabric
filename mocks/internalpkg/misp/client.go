// Code generated by mockery v2.52.2. DO NOT EDIT.

package misp_mocks

import (
	context "context"

	client "github.com/blackstork-io/fabric/internal/misp/client"

	mock "github.com/stretchr/testify/mock"
)

// Client is an autogenerated mock type for the Client type
type Client struct {
	mock.Mock
}

type Client_Expecter struct {
	mock *mock.Mock
}

func (_m *Client) EXPECT() *Client_Expecter {
	return &Client_Expecter{mock: &_m.Mock}
}

// AddEventReport provides a mock function with given fields: ctx, req
func (_m *Client) AddEventReport(ctx context.Context, req client.AddEventReportRequest) (client.AddEventReportResponse, error) {
	ret := _m.Called(ctx, req)

	if len(ret) == 0 {
		panic("no return value specified for AddEventReport")
	}

	var r0 client.AddEventReportResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, client.AddEventReportRequest) (client.AddEventReportResponse, error)); ok {
		return rf(ctx, req)
	}
	if rf, ok := ret.Get(0).(func(context.Context, client.AddEventReportRequest) client.AddEventReportResponse); ok {
		r0 = rf(ctx, req)
	} else {
		r0 = ret.Get(0).(client.AddEventReportResponse)
	}

	if rf, ok := ret.Get(1).(func(context.Context, client.AddEventReportRequest) error); ok {
		r1 = rf(ctx, req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Client_AddEventReport_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'AddEventReport'
type Client_AddEventReport_Call struct {
	*mock.Call
}

// AddEventReport is a helper method to define mock.On call
//   - ctx context.Context
//   - req client.AddEventReportRequest
func (_e *Client_Expecter) AddEventReport(ctx interface{}, req interface{}) *Client_AddEventReport_Call {
	return &Client_AddEventReport_Call{Call: _e.mock.On("AddEventReport", ctx, req)}
}

func (_c *Client_AddEventReport_Call) Run(run func(ctx context.Context, req client.AddEventReportRequest)) *Client_AddEventReport_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(client.AddEventReportRequest))
	})
	return _c
}

func (_c *Client_AddEventReport_Call) Return(resp client.AddEventReportResponse, err error) *Client_AddEventReport_Call {
	_c.Call.Return(resp, err)
	return _c
}

func (_c *Client_AddEventReport_Call) RunAndReturn(run func(context.Context, client.AddEventReportRequest) (client.AddEventReportResponse, error)) *Client_AddEventReport_Call {
	_c.Call.Return(run)
	return _c
}

// RestSearchEvents provides a mock function with given fields: ctx, req
func (_m *Client) RestSearchEvents(ctx context.Context, req client.RestSearchEventsRequest) (client.RestSearchEventsResponse, error) {
	ret := _m.Called(ctx, req)

	if len(ret) == 0 {
		panic("no return value specified for RestSearchEvents")
	}

	var r0 client.RestSearchEventsResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, client.RestSearchEventsRequest) (client.RestSearchEventsResponse, error)); ok {
		return rf(ctx, req)
	}
	if rf, ok := ret.Get(0).(func(context.Context, client.RestSearchEventsRequest) client.RestSearchEventsResponse); ok {
		r0 = rf(ctx, req)
	} else {
		r0 = ret.Get(0).(client.RestSearchEventsResponse)
	}

	if rf, ok := ret.Get(1).(func(context.Context, client.RestSearchEventsRequest) error); ok {
		r1 = rf(ctx, req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Client_RestSearchEvents_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'RestSearchEvents'
type Client_RestSearchEvents_Call struct {
	*mock.Call
}

// RestSearchEvents is a helper method to define mock.On call
//   - ctx context.Context
//   - req client.RestSearchEventsRequest
func (_e *Client_Expecter) RestSearchEvents(ctx interface{}, req interface{}) *Client_RestSearchEvents_Call {
	return &Client_RestSearchEvents_Call{Call: _e.mock.On("RestSearchEvents", ctx, req)}
}

func (_c *Client_RestSearchEvents_Call) Run(run func(ctx context.Context, req client.RestSearchEventsRequest)) *Client_RestSearchEvents_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(client.RestSearchEventsRequest))
	})
	return _c
}

func (_c *Client_RestSearchEvents_Call) Return(events client.RestSearchEventsResponse, err error) *Client_RestSearchEvents_Call {
	_c.Call.Return(events, err)
	return _c
}

func (_c *Client_RestSearchEvents_Call) RunAndReturn(run func(context.Context, client.RestSearchEventsRequest) (client.RestSearchEventsResponse, error)) *Client_RestSearchEvents_Call {
	_c.Call.Return(run)
	return _c
}

// NewClient creates a new instance of Client. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewClient(t interface {
	mock.TestingT
	Cleanup(func())
}) *Client {
	mock := &Client{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
