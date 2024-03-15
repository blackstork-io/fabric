// Code generated by mockery v2.42.1. DO NOT EDIT.

package client_mocks

import (
	context "context"

	client "github.com/blackstork-io/fabric/internal/hackerone/client"

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

// GetAllReports provides a mock function with given fields: ctx, req
func (_m *Client) GetAllReports(ctx context.Context, req *client.GetAllReportsReq) (*client.GetAllReportsRes, error) {
	ret := _m.Called(ctx, req)

	if len(ret) == 0 {
		panic("no return value specified for GetAllReports")
	}

	var r0 *client.GetAllReportsRes
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *client.GetAllReportsReq) (*client.GetAllReportsRes, error)); ok {
		return rf(ctx, req)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *client.GetAllReportsReq) *client.GetAllReportsRes); ok {
		r0 = rf(ctx, req)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*client.GetAllReportsRes)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *client.GetAllReportsReq) error); ok {
		r1 = rf(ctx, req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Client_GetAllReports_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetAllReports'
type Client_GetAllReports_Call struct {
	*mock.Call
}

// GetAllReports is a helper method to define mock.On call
//   - ctx context.Context
//   - req *client.GetAllReportsReq
func (_e *Client_Expecter) GetAllReports(ctx interface{}, req interface{}) *Client_GetAllReports_Call {
	return &Client_GetAllReports_Call{Call: _e.mock.On("GetAllReports", ctx, req)}
}

func (_c *Client_GetAllReports_Call) Run(run func(ctx context.Context, req *client.GetAllReportsReq)) *Client_GetAllReports_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*client.GetAllReportsReq))
	})
	return _c
}

func (_c *Client_GetAllReports_Call) Return(_a0 *client.GetAllReportsRes, _a1 error) *Client_GetAllReports_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *Client_GetAllReports_Call) RunAndReturn(run func(context.Context, *client.GetAllReportsReq) (*client.GetAllReportsRes, error)) *Client_GetAllReports_Call {
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
