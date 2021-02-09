// Code generated by MockGen. DO NOT EDIT.
// Source: ./src/utils/time.go

// Package mock_utils is a generated GoMock package.
package mock_utils

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockTimeSource is a mock of TimeSource interface
type MockTimeSource struct {
	ctrl     *gomock.Controller
	recorder *MockTimeSourceMockRecorder
}

// MockTimeSourceMockRecorder is the mock recorder for MockTimeSource
type MockTimeSourceMockRecorder struct {
	mock *MockTimeSource
}

// NewMockTimeSource creates a new mock instance
func NewMockTimeSource(ctrl *gomock.Controller) *MockTimeSource {
	mock := &MockTimeSource{ctrl: ctrl}
	mock.recorder = &MockTimeSourceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockTimeSource) EXPECT() *MockTimeSourceMockRecorder {
	return m.recorder
}

// UnixNow mocks base method
func (m *MockTimeSource) UnixNow() int64 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UnixNow")
	ret0, _ := ret[0].(int64)
	return ret0
}

// UnixNow indicates an expected call of UnixNow
func (mr *MockTimeSourceMockRecorder) UnixNow() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UnixNow", reflect.TypeOf((*MockTimeSource)(nil).UnixNow))
}

// UnixNanoNow mocks base method
func (m *MockTimeSource) UnixNanoNow() int64 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UnixNanoNow")
	ret0, _ := ret[0].(int64)
	return ret0
}

// UnixNanoNow indicates an expected call of UnixNanoNow
func (mr *MockTimeSourceMockRecorder) UnixNanoNow() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UnixNanoNow", reflect.TypeOf((*MockTimeSource)(nil).UnixNanoNow))
}
