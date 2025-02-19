// Code generated by mockery v2.52.2. DO NOT EDIT.

package client_mocks

import (
	context "context"

	client "github.com/blackstork-io/fabric/internal/virustotal/client"

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

// GetGroupAPIUsage provides a mock function with given fields: ctx, req
func (_m *Client) GetGroupAPIUsage(ctx context.Context, req *client.GetGroupAPIUsageReq) (*client.GetGroupAPIUsageRes, error) {
	ret := _m.Called(ctx, req)

	if len(ret) == 0 {
		panic("no return value specified for GetGroupAPIUsage")
	}

	var r0 *client.GetGroupAPIUsageRes
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *client.GetGroupAPIUsageReq) (*client.GetGroupAPIUsageRes, error)); ok {
		return rf(ctx, req)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *client.GetGroupAPIUsageReq) *client.GetGroupAPIUsageRes); ok {
		r0 = rf(ctx, req)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*client.GetGroupAPIUsageRes)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *client.GetGroupAPIUsageReq) error); ok {
		r1 = rf(ctx, req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Client_GetGroupAPIUsage_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetGroupAPIUsage'
type Client_GetGroupAPIUsage_Call struct {
	*mock.Call
}

// GetGroupAPIUsage is a helper method to define mock.On call
//   - ctx context.Context
//   - req *client.GetGroupAPIUsageReq
func (_e *Client_Expecter) GetGroupAPIUsage(ctx interface{}, req interface{}) *Client_GetGroupAPIUsage_Call {
	return &Client_GetGroupAPIUsage_Call{Call: _e.mock.On("GetGroupAPIUsage", ctx, req)}
}

func (_c *Client_GetGroupAPIUsage_Call) Run(run func(ctx context.Context, req *client.GetGroupAPIUsageReq)) *Client_GetGroupAPIUsage_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*client.GetGroupAPIUsageReq))
	})
	return _c
}

func (_c *Client_GetGroupAPIUsage_Call) Return(_a0 *client.GetGroupAPIUsageRes, _a1 error) *Client_GetGroupAPIUsage_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *Client_GetGroupAPIUsage_Call) RunAndReturn(run func(context.Context, *client.GetGroupAPIUsageReq) (*client.GetGroupAPIUsageRes, error)) *Client_GetGroupAPIUsage_Call {
	_c.Call.Return(run)
	return _c
}

// GetUserAPIUsage provides a mock function with given fields: ctx, req
func (_m *Client) GetUserAPIUsage(ctx context.Context, req *client.GetUserAPIUsageReq) (*client.GetUserAPIUsageRes, error) {
	ret := _m.Called(ctx, req)

	if len(ret) == 0 {
		panic("no return value specified for GetUserAPIUsage")
	}

	var r0 *client.GetUserAPIUsageRes
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *client.GetUserAPIUsageReq) (*client.GetUserAPIUsageRes, error)); ok {
		return rf(ctx, req)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *client.GetUserAPIUsageReq) *client.GetUserAPIUsageRes); ok {
		r0 = rf(ctx, req)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*client.GetUserAPIUsageRes)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *client.GetUserAPIUsageReq) error); ok {
		r1 = rf(ctx, req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Client_GetUserAPIUsage_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetUserAPIUsage'
type Client_GetUserAPIUsage_Call struct {
	*mock.Call
}

// GetUserAPIUsage is a helper method to define mock.On call
//   - ctx context.Context
//   - req *client.GetUserAPIUsageReq
func (_e *Client_Expecter) GetUserAPIUsage(ctx interface{}, req interface{}) *Client_GetUserAPIUsage_Call {
	return &Client_GetUserAPIUsage_Call{Call: _e.mock.On("GetUserAPIUsage", ctx, req)}
}

func (_c *Client_GetUserAPIUsage_Call) Run(run func(ctx context.Context, req *client.GetUserAPIUsageReq)) *Client_GetUserAPIUsage_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*client.GetUserAPIUsageReq))
	})
	return _c
}

func (_c *Client_GetUserAPIUsage_Call) Return(_a0 *client.GetUserAPIUsageRes, _a1 error) *Client_GetUserAPIUsage_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *Client_GetUserAPIUsage_Call) RunAndReturn(run func(context.Context, *client.GetUserAPIUsageReq) (*client.GetUserAPIUsageRes, error)) *Client_GetUserAPIUsage_Call {
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
