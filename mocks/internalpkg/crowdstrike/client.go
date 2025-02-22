// Code generated by mockery v2.52.2. DO NOT EDIT.

package crowdstrike_mocks

import (
	crowdstrike "github.com/blackstork-io/fabric/internal/crowdstrike"
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

// CspmRegistration provides a mock function with no fields
func (_m *Client) CspmRegistration() crowdstrike.CspmRegistrationClient {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for CspmRegistration")
	}

	var r0 crowdstrike.CspmRegistrationClient
	if rf, ok := ret.Get(0).(func() crowdstrike.CspmRegistrationClient); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(crowdstrike.CspmRegistrationClient)
		}
	}

	return r0
}

// Client_CspmRegistration_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CspmRegistration'
type Client_CspmRegistration_Call struct {
	*mock.Call
}

// CspmRegistration is a helper method to define mock.On call
func (_e *Client_Expecter) CspmRegistration() *Client_CspmRegistration_Call {
	return &Client_CspmRegistration_Call{Call: _e.mock.On("CspmRegistration")}
}

func (_c *Client_CspmRegistration_Call) Run(run func()) *Client_CspmRegistration_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *Client_CspmRegistration_Call) Return(_a0 crowdstrike.CspmRegistrationClient) *Client_CspmRegistration_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Client_CspmRegistration_Call) RunAndReturn(run func() crowdstrike.CspmRegistrationClient) *Client_CspmRegistration_Call {
	_c.Call.Return(run)
	return _c
}

// Detects provides a mock function with no fields
func (_m *Client) Detects() crowdstrike.DetectsClient {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Detects")
	}

	var r0 crowdstrike.DetectsClient
	if rf, ok := ret.Get(0).(func() crowdstrike.DetectsClient); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(crowdstrike.DetectsClient)
		}
	}

	return r0
}

// Client_Detects_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Detects'
type Client_Detects_Call struct {
	*mock.Call
}

// Detects is a helper method to define mock.On call
func (_e *Client_Expecter) Detects() *Client_Detects_Call {
	return &Client_Detects_Call{Call: _e.mock.On("Detects")}
}

func (_c *Client_Detects_Call) Run(run func()) *Client_Detects_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *Client_Detects_Call) Return(_a0 crowdstrike.DetectsClient) *Client_Detects_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Client_Detects_Call) RunAndReturn(run func() crowdstrike.DetectsClient) *Client_Detects_Call {
	_c.Call.Return(run)
	return _c
}

// Discover provides a mock function with no fields
func (_m *Client) Discover() crowdstrike.DiscoverClient {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Discover")
	}

	var r0 crowdstrike.DiscoverClient
	if rf, ok := ret.Get(0).(func() crowdstrike.DiscoverClient); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(crowdstrike.DiscoverClient)
		}
	}

	return r0
}

// Client_Discover_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Discover'
type Client_Discover_Call struct {
	*mock.Call
}

// Discover is a helper method to define mock.On call
func (_e *Client_Expecter) Discover() *Client_Discover_Call {
	return &Client_Discover_Call{Call: _e.mock.On("Discover")}
}

func (_c *Client_Discover_Call) Run(run func()) *Client_Discover_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *Client_Discover_Call) Return(_a0 crowdstrike.DiscoverClient) *Client_Discover_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Client_Discover_Call) RunAndReturn(run func() crowdstrike.DiscoverClient) *Client_Discover_Call {
	_c.Call.Return(run)
	return _c
}

// Intel provides a mock function with no fields
func (_m *Client) Intel() crowdstrike.IntelClient {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Intel")
	}

	var r0 crowdstrike.IntelClient
	if rf, ok := ret.Get(0).(func() crowdstrike.IntelClient); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(crowdstrike.IntelClient)
		}
	}

	return r0
}

// Client_Intel_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Intel'
type Client_Intel_Call struct {
	*mock.Call
}

// Intel is a helper method to define mock.On call
func (_e *Client_Expecter) Intel() *Client_Intel_Call {
	return &Client_Intel_Call{Call: _e.mock.On("Intel")}
}

func (_c *Client_Intel_Call) Run(run func()) *Client_Intel_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *Client_Intel_Call) Return(_a0 crowdstrike.IntelClient) *Client_Intel_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Client_Intel_Call) RunAndReturn(run func() crowdstrike.IntelClient) *Client_Intel_Call {
	_c.Call.Return(run)
	return _c
}

// SpotlightVulnerabilities provides a mock function with no fields
func (_m *Client) SpotlightVulnerabilities() crowdstrike.SpotVulnerabilitiesClient {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for SpotlightVulnerabilities")
	}

	var r0 crowdstrike.SpotVulnerabilitiesClient
	if rf, ok := ret.Get(0).(func() crowdstrike.SpotVulnerabilitiesClient); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(crowdstrike.SpotVulnerabilitiesClient)
		}
	}

	return r0
}

// Client_SpotlightVulnerabilities_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SpotlightVulnerabilities'
type Client_SpotlightVulnerabilities_Call struct {
	*mock.Call
}

// SpotlightVulnerabilities is a helper method to define mock.On call
func (_e *Client_Expecter) SpotlightVulnerabilities() *Client_SpotlightVulnerabilities_Call {
	return &Client_SpotlightVulnerabilities_Call{Call: _e.mock.On("SpotlightVulnerabilities")}
}

func (_c *Client_SpotlightVulnerabilities_Call) Run(run func()) *Client_SpotlightVulnerabilities_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *Client_SpotlightVulnerabilities_Call) Return(_a0 crowdstrike.SpotVulnerabilitiesClient) *Client_SpotlightVulnerabilities_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Client_SpotlightVulnerabilities_Call) RunAndReturn(run func() crowdstrike.SpotVulnerabilitiesClient) *Client_SpotlightVulnerabilities_Call {
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
