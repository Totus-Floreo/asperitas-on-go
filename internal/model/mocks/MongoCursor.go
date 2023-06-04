// Code generated by MockGen. DO NOT EDIT.
// Source: mongo_cursor.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockICursor is a mock of ICursor interface.
type MockICursor struct {
	ctrl     *gomock.Controller
	recorder *MockICursorMockRecorder
}

// MockICursorMockRecorder is the mock recorder for MockICursor.
type MockICursorMockRecorder struct {
	mock *MockICursor
}

// NewMockICursor creates a new mock instance.
func NewMockICursor(ctrl *gomock.Controller) *MockICursor {
	mock := &MockICursor{ctrl: ctrl}
	mock.recorder = &MockICursorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockICursor) EXPECT() *MockICursorMockRecorder {
	return m.recorder
}

// All mocks base method.
func (m *MockICursor) All(arg0 context.Context, arg1 interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "All", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// All indicates an expected call of All.
func (mr *MockICursorMockRecorder) All(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "All", reflect.TypeOf((*MockICursor)(nil).All), arg0, arg1)
}

// Close mocks base method.
func (m *MockICursor) Close(arg0 context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockICursorMockRecorder) Close(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockICursor)(nil).Close), arg0)
}

// Decode mocks base method.
func (m *MockICursor) Decode(arg0 interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Decode", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Decode indicates an expected call of Decode.
func (mr *MockICursorMockRecorder) Decode(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Decode", reflect.TypeOf((*MockICursor)(nil).Decode), arg0)
}

// Err mocks base method.
func (m *MockICursor) Err() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Err")
	ret0, _ := ret[0].(error)
	return ret0
}

// Err indicates an expected call of Err.
func (mr *MockICursorMockRecorder) Err() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Err", reflect.TypeOf((*MockICursor)(nil).Err))
}

// TryNext mocks base method.
func (m *MockICursor) TryNext(arg0 context.Context) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TryNext", arg0)
	ret0, _ := ret[0].(bool)
	return ret0
}

// TryNext indicates an expected call of TryNext.
func (mr *MockICursorMockRecorder) TryNext(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TryNext", reflect.TypeOf((*MockICursor)(nil).TryNext), arg0)
}