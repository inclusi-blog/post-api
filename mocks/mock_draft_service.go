// Code generated by MockGen. DO NOT EDIT.
// Source: service/draft_service.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	models "post-api/models"
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
func (m *MockDraftService) SaveDraft(postData models.UpsertDraft, ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveDraft", postData, ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveDraft indicates an expected call of SaveDraft
func (mr *MockDraftServiceMockRecorder) SaveDraft(postData, ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveDraft", reflect.TypeOf((*MockDraftService)(nil).SaveDraft), postData, ctx)
}