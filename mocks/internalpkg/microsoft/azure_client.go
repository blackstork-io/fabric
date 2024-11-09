// Code generated by mockery v2.42.1. DO NOT EDIT.

package microsoft_mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	plugindata "github.com/blackstork-io/fabric/plugin/plugindata"

	url "net/url"
)

// AzureClient is an autogenerated mock type for the AzureClient type
type AzureClient struct {
	mock.Mock
}

type AzureClient_Expecter struct {
	mock *mock.Mock
}

func (_m *AzureClient) EXPECT() *AzureClient_Expecter {
	return &AzureClient_Expecter{mock: &_m.Mock}
}

// QueryObjects provides a mock function with given fields: ctx, endpoint, queryParams, size
func (_m *AzureClient) QueryObjects(ctx context.Context, endpoint string, queryParams url.Values, size int) (plugindata.List, error) {
	ret := _m.Called(ctx, endpoint, queryParams, size)

	if len(ret) == 0 {
		panic("no return value specified for QueryObjects")
	}

	var r0 plugindata.List
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, url.Values, int) (plugindata.List, error)); ok {
		return rf(ctx, endpoint, queryParams, size)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, url.Values, int) plugindata.List); ok {
		r0 = rf(ctx, endpoint, queryParams, size)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(plugindata.List)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, url.Values, int) error); ok {
		r1 = rf(ctx, endpoint, queryParams, size)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// AzureClient_QueryObjects_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'QueryObjects'
type AzureClient_QueryObjects_Call struct {
	*mock.Call
}

// QueryObjects is a helper method to define mock.On call
//   - ctx context.Context
//   - endpoint string
//   - queryParams url.Values
//   - size int
func (_e *AzureClient_Expecter) QueryObjects(ctx interface{}, endpoint interface{}, queryParams interface{}, size interface{}) *AzureClient_QueryObjects_Call {
	return &AzureClient_QueryObjects_Call{Call: _e.mock.On("QueryObjects", ctx, endpoint, queryParams, size)}
}

func (_c *AzureClient_QueryObjects_Call) Run(run func(ctx context.Context, endpoint string, queryParams url.Values, size int)) *AzureClient_QueryObjects_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(url.Values), args[3].(int))
	})
	return _c
}

func (_c *AzureClient_QueryObjects_Call) Return(objects plugindata.List, err error) *AzureClient_QueryObjects_Call {
	_c.Call.Return(objects, err)
	return _c
}

func (_c *AzureClient_QueryObjects_Call) RunAndReturn(run func(context.Context, string, url.Values, int) (plugindata.List, error)) *AzureClient_QueryObjects_Call {
	_c.Call.Return(run)
	return _c
}

// NewAzureClient creates a new instance of AzureClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewAzureClient(t interface {
	mock.TestingT
	Cleanup(func())
}) *AzureClient {
	mock := &AzureClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
