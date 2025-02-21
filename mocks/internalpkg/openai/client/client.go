// Code generated by mockery v2.52.2. DO NOT EDIT.

package client_mocks

import (
	context "context"

	client "github.com/blackstork-io/fabric/internal/openai/client"

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

// GenerateChatCompletion provides a mock function with given fields: ctx, params
func (_m *Client) GenerateChatCompletion(ctx context.Context, params *client.ChatCompletionParams) (*client.ChatCompletionResult, error) {
	ret := _m.Called(ctx, params)

	if len(ret) == 0 {
		panic("no return value specified for GenerateChatCompletion")
	}

	var r0 *client.ChatCompletionResult
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *client.ChatCompletionParams) (*client.ChatCompletionResult, error)); ok {
		return rf(ctx, params)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *client.ChatCompletionParams) *client.ChatCompletionResult); ok {
		r0 = rf(ctx, params)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*client.ChatCompletionResult)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *client.ChatCompletionParams) error); ok {
		r1 = rf(ctx, params)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Client_GenerateChatCompletion_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GenerateChatCompletion'
type Client_GenerateChatCompletion_Call struct {
	*mock.Call
}

// GenerateChatCompletion is a helper method to define mock.On call
//   - ctx context.Context
//   - params *client.ChatCompletionParams
func (_e *Client_Expecter) GenerateChatCompletion(ctx interface{}, params interface{}) *Client_GenerateChatCompletion_Call {
	return &Client_GenerateChatCompletion_Call{Call: _e.mock.On("GenerateChatCompletion", ctx, params)}
}

func (_c *Client_GenerateChatCompletion_Call) Run(run func(ctx context.Context, params *client.ChatCompletionParams)) *Client_GenerateChatCompletion_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*client.ChatCompletionParams))
	})
	return _c
}

func (_c *Client_GenerateChatCompletion_Call) Return(_a0 *client.ChatCompletionResult, _a1 error) *Client_GenerateChatCompletion_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *Client_GenerateChatCompletion_Call) RunAndReturn(run func(context.Context, *client.ChatCompletionParams) (*client.ChatCompletionResult, error)) *Client_GenerateChatCompletion_Call {
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
