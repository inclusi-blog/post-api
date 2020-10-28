// Code generated by MockGen. DO NOT EDIT.
// Source: preview_posts_repository.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	db "post-api/models/db"
	reflect "reflect"
)

// MockPreviewPostsRepository is a mock of PreviewPostsRepository interface
type MockPreviewPostsRepository struct {
	ctrl     *gomock.Controller
	recorder *MockPreviewPostsRepositoryMockRecorder
}

// MockPreviewPostsRepositoryMockRecorder is the mock recorder for MockPreviewPostsRepository
type MockPreviewPostsRepositoryMockRecorder struct {
	mock *MockPreviewPostsRepository
}

// NewMockPreviewPostsRepository creates a new mock instance
func NewMockPreviewPostsRepository(ctrl *gomock.Controller) *MockPreviewPostsRepository {
	mock := &MockPreviewPostsRepository{ctrl: ctrl}
	mock.recorder = &MockPreviewPostsRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockPreviewPostsRepository) EXPECT() *MockPreviewPostsRepositoryMockRecorder {
	return m.recorder
}

// SavePreview mocks base method
func (m *MockPreviewPostsRepository) SavePreview(ctx context.Context, post db.PreviewPost) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SavePreview", ctx, post)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SavePreview indicates an expected call of SavePreview
func (mr *MockPreviewPostsRepositoryMockRecorder) SavePreview(ctx, post interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SavePreview", reflect.TypeOf((*MockPreviewPostsRepository)(nil).SavePreview), ctx, post)
}
