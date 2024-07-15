// Code generated by MockGen. DO NOT EDIT.
// Source: cache_manager.go

// Package cache is a generated GoMock package.
package cache

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockCacheManager is a mock of CacheManager interface.
type MockCacheManager struct {
	ctrl     *gomock.Controller
	recorder *MockCacheManagerMockRecorder
}

// MockCacheManagerMockRecorder is the mock recorder for MockCacheManager.
type MockCacheManagerMockRecorder struct {
	mock *MockCacheManager
}

// NewMockCacheManager creates a new mock instance.
func NewMockCacheManager(ctrl *gomock.Controller) *MockCacheManager {
	mock := &MockCacheManager{ctrl: ctrl}
	mock.recorder = &MockCacheManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCacheManager) EXPECT() *MockCacheManagerMockRecorder {
	return m.recorder
}

// Del mocks base method.
func (m *MockCacheManager) Del(ctx context.Context, key string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Del", ctx, key)
	ret0, _ := ret[0].(error)
	return ret0
}

// Del indicates an expected call of Del.
func (mr *MockCacheManagerMockRecorder) Del(ctx, key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Del", reflect.TypeOf((*MockCacheManager)(nil).Del), ctx, key)
}

// Get mocks base method.
func (m *MockCacheManager) Get(ctx context.Context, key string) *string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, key)
	ret0, _ := ret[0].(*string)
	return ret0
}

// Get indicates an expected call of Get.
func (mr *MockCacheManagerMockRecorder) Get(ctx, key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockCacheManager)(nil).Get), ctx, key)
}

// Set mocks base method.
func (m *MockCacheManager) Set(ctx context.Context, key string, val any) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Set", ctx, key, val)
	ret0, _ := ret[0].(error)
	return ret0
}

// Set indicates an expected call of Set.
func (mr *MockCacheManagerMockRecorder) Set(ctx, key, val interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Set", reflect.TypeOf((*MockCacheManager)(nil).Set), ctx, key, val)
}