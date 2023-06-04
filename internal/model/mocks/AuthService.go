// Code generated by MockGen. DO NOT EDIT.
// Source: ./internal/model/auth_service.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockIAuthService is a mock of IAuthService interface.
type MockIAuthService struct {
	ctrl     *gomock.Controller
	recorder *MockIAuthServiceMockRecorder
}

// MockIAuthServiceMockRecorder is the mock recorder for MockIAuthService.
type MockIAuthServiceMockRecorder struct {
	mock *MockIAuthService
}

// NewMockIAuthService creates a new mock instance.
func NewMockIAuthService(ctrl *gomock.Controller) *MockIAuthService {
	mock := &MockIAuthService{ctrl: ctrl}
	mock.recorder = &MockIAuthServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIAuthService) EXPECT() *MockIAuthServiceMockRecorder {
	return m.recorder
}

// LogIn mocks base method.
func (m *MockIAuthService) LogIn(arg0 context.Context, arg1, arg2 string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LogIn", arg0, arg1, arg2)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LogIn indicates an expected call of LogIn.
func (mr *MockIAuthServiceMockRecorder) LogIn(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LogIn", reflect.TypeOf((*MockIAuthService)(nil).LogIn), arg0, arg1, arg2)
}

// SignUp mocks base method.
func (m *MockIAuthService) SignUp(arg0 context.Context, arg1, arg2 string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SignUp", arg0, arg1, arg2)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SignUp indicates an expected call of SignUp.
func (mr *MockIAuthServiceMockRecorder) SignUp(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SignUp", reflect.TypeOf((*MockIAuthService)(nil).SignUp), arg0, arg1, arg2)
}