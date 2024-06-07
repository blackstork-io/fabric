// Code generated by mockery v2.42.1. DO NOT EDIT.

package client_mocks

import (
	context "context"

	client "github.com/blackstork-io/fabric/internal/sentinel/client"

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

// GetClientCredentialsToken provides a mock function with given fields: ctx, req
func (_m *Client) GetClientCredentialsToken(ctx context.Context, req *client.GetClientCredentialsTokenReq) (*client.GetClientCredentialsTokenRes, error) {
	ret := _m.Called(ctx, req)

	if len(ret) == 0 {
		panic("no return value specified for GetClientCredentialsToken")
	}

	var r0 *client.GetClientCredentialsTokenRes
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *client.GetClientCredentialsTokenReq) (*client.GetClientCredentialsTokenRes, error)); ok {
		return rf(ctx, req)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *client.GetClientCredentialsTokenReq) *client.GetClientCredentialsTokenRes); ok {
		r0 = rf(ctx, req)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*client.GetClientCredentialsTokenRes)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *client.GetClientCredentialsTokenReq) error); ok {
		r1 = rf(ctx, req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Client_GetClientCredentialsToken_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetClientCredentialsToken'
type Client_GetClientCredentialsToken_Call struct {
	*mock.Call
}

// GetClientCredentialsToken is a helper method to define mock.On call
//   - ctx context.Context
//   - req *client.GetClientCredentialsTokenReq
func (_e *Client_Expecter) GetClientCredentialsToken(ctx interface{}, req interface{}) *Client_GetClientCredentialsToken_Call {
	return &Client_GetClientCredentialsToken_Call{Call: _e.mock.On("GetClientCredentialsToken", ctx, req)}
}

func (_c *Client_GetClientCredentialsToken_Call) Run(run func(ctx context.Context, req *client.GetClientCredentialsTokenReq)) *Client_GetClientCredentialsToken_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*client.GetClientCredentialsTokenReq))
	})
	return _c
}

func (_c *Client_GetClientCredentialsToken_Call) Return(_a0 *client.GetClientCredentialsTokenRes, _a1 error) *Client_GetClientCredentialsToken_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *Client_GetClientCredentialsToken_Call) RunAndReturn(run func(context.Context, *client.GetClientCredentialsTokenReq) (*client.GetClientCredentialsTokenRes, error)) *Client_GetClientCredentialsToken_Call {
	_c.Call.Return(run)
	return _c
}

// ListIncidents provides a mock function with given fields: ctx, req
func (_m *Client) ListIncidents(ctx context.Context, req *client.ListIncidentsReq) (*client.ListIncidentsRes, error) {
	ret := _m.Called(ctx, req)

	if len(ret) == 0 {
		panic("no return value specified for ListIncidents")
	}

	var r0 *client.ListIncidentsRes
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *client.ListIncidentsReq) (*client.ListIncidentsRes, error)); ok {
		return rf(ctx, req)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *client.ListIncidentsReq) *client.ListIncidentsRes); ok {
		r0 = rf(ctx, req)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*client.ListIncidentsRes)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *client.ListIncidentsReq) error); ok {
		r1 = rf(ctx, req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Client_ListIncidents_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ListIncidents'
type Client_ListIncidents_Call struct {
	*mock.Call
}

// ListIncidents is a helper method to define mock.On call
//   - ctx context.Context
//   - req *client.ListIncidentsReq
func (_e *Client_Expecter) ListIncidents(ctx interface{}, req interface{}) *Client_ListIncidents_Call {
	return &Client_ListIncidents_Call{Call: _e.mock.On("ListIncidents", ctx, req)}
}

func (_c *Client_ListIncidents_Call) Run(run func(ctx context.Context, req *client.ListIncidentsReq)) *Client_ListIncidents_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*client.ListIncidentsReq))
	})
	return _c
}

func (_c *Client_ListIncidents_Call) Return(_a0 *client.ListIncidentsRes, _a1 error) *Client_ListIncidents_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *Client_ListIncidents_Call) RunAndReturn(run func(context.Context, *client.ListIncidentsReq) (*client.ListIncidentsRes, error)) *Client_ListIncidents_Call {
	_c.Call.Return(run)
	return _c
}

// UseAuth provides a mock function with given fields: token
func (_m *Client) UseAuth(token string) {
	_m.Called(token)
}

// Client_UseAuth_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UseAuth'
type Client_UseAuth_Call struct {
	*mock.Call
}

// UseAuth is a helper method to define mock.On call
//   - token string
func (_e *Client_Expecter) UseAuth(token interface{}) *Client_UseAuth_Call {
	return &Client_UseAuth_Call{Call: _e.mock.On("UseAuth", token)}
}

func (_c *Client_UseAuth_Call) Run(run func(token string)) *Client_UseAuth_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *Client_UseAuth_Call) Return() *Client_UseAuth_Call {
	_c.Call.Return()
	return _c
}

func (_c *Client_UseAuth_Call) RunAndReturn(run func(string)) *Client_UseAuth_Call {
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
