// Code generated by MockGen. DO NOT EDIT.
// Source: post_validator.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	db "post-api/models/db"
	reflect "reflect"
)

// MockPostValidator is a mock of PostValidator interface
type MockPostValidator struct {
	ctrl     *gomock.Controller
	recorder *MockPostValidatorMockRecorder
}

// MockPostValidatorMockRecorder is the mock recorder for MockPostValidator
type MockPostValidatorMockRecorder struct {
	mock *MockPostValidator
}

// NewMockPostValidator creates a new mock instance
func NewMockPostValidator(ctrl *gomock.Controller) *MockPostValidator {
	mock := &MockPostValidator{ctrl: ctrl}
	mock.recorder = &MockPostValidatorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockPostValidator) EXPECT() *MockPostValidatorMockRecorder {
	return m.recorder
}

// ValidateAndGetReadTime mocks base method
func (m *MockPostValidator) ValidateAndGetReadTime(draft *db.Draft, ctx context.Context) (string, int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ValidateAndGetReadTime", draft, ctx)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(int)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// ValidateAndGetReadTime indicates an expected call of ValidateAndGetReadTime
func (mr *MockPostValidatorMockRecorder) ValidateAndGetReadTime(draft, ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidateAndGetReadTime", reflect.TypeOf((*MockPostValidator)(nil).ValidateAndGetReadTime), draft, ctx)
}
