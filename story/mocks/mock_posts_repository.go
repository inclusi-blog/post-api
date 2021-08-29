// Code generated by MockGen. DO NOT EDIT.
// Source: posts_repository.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	helper "post-api/helper"
	db "post-api/story/models/db"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
)

// MockPostsRepository is a mock of PostsRepository interface.
type MockPostsRepository struct {
	ctrl     *gomock.Controller
	recorder *MockPostsRepositoryMockRecorder
}

// MockPostsRepositoryMockRecorder is the mock recorder for MockPostsRepository.
type MockPostsRepositoryMockRecorder struct {
	mock *MockPostsRepository
}

// NewMockPostsRepository creates a new mock instance.
func NewMockPostsRepository(ctrl *gomock.Controller) *MockPostsRepository {
	mock := &MockPostsRepository{ctrl: ctrl}
	mock.recorder = &MockPostsRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPostsRepository) EXPECT() *MockPostsRepositoryMockRecorder {
	return m.recorder
}

// AddInterests mocks base method.
func (m *MockPostsRepository) AddInterests(ctx context.Context, transaction helper.Transaction, postID uuid.UUID, interests []uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddInterests", ctx, transaction, postID, interests)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddInterests indicates an expected call of AddInterests.
func (mr *MockPostsRepositoryMockRecorder) AddInterests(ctx, transaction, postID, interests interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddInterests", reflect.TypeOf((*MockPostsRepository)(nil).AddInterests), ctx, transaction, postID, interests)
}

// CreatePost mocks base method.
func (m *MockPostsRepository) CreatePost(ctx context.Context, tx helper.Transaction, post db.PublishPost) (uuid.UUID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreatePost", ctx, tx, post)
	ret0, _ := ret[0].(uuid.UUID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreatePost indicates an expected call of CreatePost.
func (mr *MockPostsRepositoryMockRecorder) CreatePost(ctx, tx, post interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreatePost", reflect.TypeOf((*MockPostsRepository)(nil).CreatePost), ctx, tx, post)
}

// Like mocks base method.
func (m *MockPostsRepository) Like(ctx context.Context, postID, userID uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Like", ctx, postID, userID)
	ret0, _ := ret[0].(error)
	return ret0
}

// Like indicates an expected call of Like.
func (mr *MockPostsRepositoryMockRecorder) Like(ctx, postID, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Like", reflect.TypeOf((*MockPostsRepository)(nil).Like), ctx, postID, userID)
}

// UnLike mocks base method.
func (m *MockPostsRepository) UnLike(ctx context.Context, postID, userID uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UnLike", ctx, postID, userID)
	ret0, _ := ret[0].(error)
	return ret0
}

// UnLike indicates an expected call of UnLike.
func (mr *MockPostsRepositoryMockRecorder) UnLike(ctx, postID, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UnLike", reflect.TypeOf((*MockPostsRepository)(nil).UnLike), ctx, postID, userID)
}