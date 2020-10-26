// Code generated by MockGen. DO NOT EDIT.
// Source: draft_service.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	golaerror "github.com/gola-glitch/gola-utils/golaerror"
	gomock "github.com/golang/mock/gomock"
	models "post-api/models"
	request "post-api/models/request"
	reflect "reflect"
)

// MockDraftService is a mock of DraftService interface
type MockDraftService struct {
	ctrl     *gomock.Controller
	recorder *MockDraftServiceMockRecorder
}

// MockDraftServiceMockRecorder is the mock recorder for MockDraftService
type MockDraftServiceMockRecorder struct {
	mock *MockDraftService
}

// NewMockDraftService creates a new mock instance
func NewMockDraftService(ctrl *gomock.Controller) *MockDraftService {
	mock := &MockDraftService{ctrl: ctrl}
	mock.recorder = &MockDraftServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockDraftService) EXPECT() *MockDraftServiceMockRecorder {
	return m.recorder
}

// SaveDraft mocks base method
func (m *MockDraftService) SaveDraft(postData models.UpsertDraft, ctx context.Context) *golaerror.Error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveDraft", postData, ctx)
	ret0, _ := ret[0].(*golaerror.Error)
	return ret0
}

// SaveDraft indicates an expected call of SaveDraft
func (mr *MockDraftServiceMockRecorder) SaveDraft(postData, ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveDraft", reflect.TypeOf((*MockDraftService)(nil).SaveDraft), postData, ctx)
}

// UpsertTagline mocks base method
func (m *MockDraftService) UpsertTagline(taglineRequest request.TaglineSaveRequest, ctx context.Context) *golaerror.Error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpsertTagline", taglineRequest, ctx)
	ret0, _ := ret[0].(*golaerror.Error)
	return ret0
}

// UpsertTagline indicates an expected call of UpsertTagline
func (mr *MockDraftServiceMockRecorder) UpsertTagline(taglineRequest, ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpsertTagline", reflect.TypeOf((*MockDraftService)(nil).UpsertTagline), taglineRequest, ctx)
}

// UpsertInterests mocks base method
func (m *MockDraftService) UpsertInterests(interestRequest request.InterestsSaveRequest, ctx context.Context) *golaerror.Error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpsertInterests", interestRequest, ctx)
	ret0, _ := ret[0].(*golaerror.Error)
	return ret0
}

// UpsertInterests indicates an expected call of UpsertInterests
func (mr *MockDraftServiceMockRecorder) UpsertInterests(interestRequest, ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpsertInterests", reflect.TypeOf((*MockDraftService)(nil).UpsertInterests), interestRequest, ctx)
}
