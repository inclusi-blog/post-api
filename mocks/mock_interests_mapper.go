// Code generated by MockGen. DO NOT EDIT.
// Source: interests_mapper.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	db "post-api/models/db"
	response "post-api/models/response"
	reflect "reflect"
)

// MockInterestsMapper is a mock of InterestsMapper interface
type MockInterestsMapper struct {
	ctrl     *gomock.Controller
	recorder *MockInterestsMapperMockRecorder
}

// MockInterestsMapperMockRecorder is the mock recorder for MockInterestsMapper
type MockInterestsMapperMockRecorder struct {
	mock *MockInterestsMapper
}

// NewMockInterestsMapper creates a new mock instance
func NewMockInterestsMapper(ctrl *gomock.Controller) *MockInterestsMapper {
	mock := &MockInterestsMapper{ctrl: ctrl}
	mock.recorder = &MockInterestsMapperMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockInterestsMapper) EXPECT() *MockInterestsMapperMockRecorder {
	return m.recorder
}

// MapUserFollowedInterest mocks base method
func (m *MockInterestsMapper) MapUserFollowedInterest(ctx context.Context, dbCategoriesAndInterests []db.CategoryAndInterest, userFollowingInterests []string) []response.CategoryAndInterest {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MapUserFollowedInterest", ctx, dbCategoriesAndInterests, userFollowingInterests)
	ret0, _ := ret[0].([]response.CategoryAndInterest)
	return ret0
}

// MapUserFollowedInterest indicates an expected call of MapUserFollowedInterest
func (mr *MockInterestsMapperMockRecorder) MapUserFollowedInterest(ctx, dbCategoriesAndInterests, userFollowingInterests interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MapUserFollowedInterest", reflect.TypeOf((*MockInterestsMapper)(nil).MapUserFollowedInterest), ctx, dbCategoriesAndInterests, userFollowingInterests)
}