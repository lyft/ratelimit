// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/envoyproxy/ratelimit/src/redis/driver (interfaces: Client)

// Package mock_driver is a generated GoMock package.
package mock_driver

import (
	driver "github.com/envoyproxy/ratelimit/src/redis/driver"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockClient is a mock of Client interface
type MockClient struct {
	ctrl     *gomock.Controller
	recorder *MockClientMockRecorder
}

// MockClientMockRecorder is the mock recorder for MockClient
type MockClientMockRecorder struct {
	mock *MockClient
}

// NewMockClient creates a new mock instance
func NewMockClient(ctrl *gomock.Controller) *MockClient {
	mock := &MockClient{ctrl: ctrl}
	mock.recorder = &MockClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockClient) EXPECT() *MockClientMockRecorder {
	return m.recorder
}

// DoCmd mocks base method
func (m *MockClient) DoCmd(rcv interface{}, cmd, key string, args ...interface{}) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{rcv, cmd, key}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "DoCmd", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// DoCmd indicates an expected call of DoCmd
func (mr *MockClientMockRecorder) DoCmd(rcv, cmd, key interface{}, args ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{rcv, cmd, key}, args...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DoCmd", reflect.TypeOf((*MockClient)(nil).DoCmd), varargs...)
}

// PipeAppend mocks base method
func (m *MockClient) PipeAppend(pipeline driver.Pipeline, rcv interface{}, cmd, key string, args ...interface{}) driver.Pipeline {
	m.ctrl.T.Helper()
	varargs := []interface{}{pipeline, rcv, cmd, key}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "PipeAppend", varargs...)
	ret0, _ := ret[0].(driver.Pipeline)
	return ret0
}

// PipeAppend indicates an expected call of PipeAppend
func (mr *MockClientMockRecorder) PipeAppend(pipeline, rcv, cmd, key interface{}, args ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{pipeline, rcv, cmd, key}, args...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PipeAppend", reflect.TypeOf((*MockClient)(nil).PipeAppend), varargs...)
}

// PipeDo mocks base method
func (m *MockClient) PipeDo(pipeline driver.Pipeline) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PipeDo", pipeline)
	ret0, _ := ret[0].(error)
	return ret0
}

// PipeDo indicates an expected call of PipeDo
func (mr *MockClientMockRecorder) PipeDo(pipeline interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PipeDo", reflect.TypeOf((*MockClient)(nil).PipeDo), pipeline)
}

// Close mocks base method
func (m *MockClient) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close
func (mr *MockClientMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockClient)(nil).Close))
}

// NumActiveConns mocks base method
func (m *MockClient) NumActiveConns() int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NumActiveConns")
	ret0, _ := ret[0].(int)
	return ret0
}

// NumActiveConns indicates an expected call of NumActiveConns
func (mr *MockClientMockRecorder) NumActiveConns() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NumActiveConns", reflect.TypeOf((*MockClient)(nil).NumActiveConns))
}

// ImplicitPipeliningEnabled mocks base method
func (m *MockClient) ImplicitPipeliningEnabled() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ImplicitPipeliningEnabled")
	ret0, _ := ret[0].(bool)
	return ret0
}

// ImplicitPipeliningEnabled indicates an expected call of ImplicitPipeliningEnabled
func (mr *MockClientMockRecorder) ImplicitPipeliningEnabled() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ImplicitPipeliningEnabled", reflect.TypeOf((*MockClient)(nil).ImplicitPipeliningEnabled))
}