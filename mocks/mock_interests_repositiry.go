// Code generated by MockGen. DO NOT EDIT.
// Source: interests_repository.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	db "post-api/models/db"
	reflect "reflect"
)

// MockInterestsRepository is a mock of InterestsRepository interface
type MockInterestsRepository struct {
	ctrl     *gomock.Controller
	recorder *MockInterestsRepositoryMockRecorder
}

// MockInterestsRepositoryMockRecorder is the mock recorder for MockInterestsRepository
type MockInterestsRepositoryMockRecorder struct {
	mock *MockInterestsRepository
}

// NewMockInterestsRepository creates a new mock instance
func NewMockInterestsRepository(ctrl *gomock.Controller) *MockInterestsRepository {
	mock := &MockInterestsRepository{ctrl: ctrl}
	mock.recorder = &MockInterestsRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockInterestsRepository) EXPECT() *MockInterestsRepositoryMockRecorder {
	return m.recorder
}

// GetInterests mocks base method
func (m *MockInterestsRepository) GetInterests(ctx context.Context, searchKeyword string, selectedTags []string) ([]db.Interest, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetInterests", ctx, searchKeyword, selectedTags)
	ret0, _ := ret[0].([]db.Interest)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetInterests indicates an expected call of GetInterests
func (mr *MockInterestsRepositoryMockRecorder) GetInterests(ctx, searchKeyword, selectedTags interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetInterests", reflect.TypeOf((*MockInterestsRepository)(nil).GetInterests), ctx, searchKeyword, selectedTags)
}
