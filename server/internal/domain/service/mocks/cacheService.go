// Code generated by MockGen. DO NOT EDIT.
// Source: ./internal/domain/service/cacheService.go

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"
	time "time"

	gomock "github.com/golang/mock/gomock"
)

// MockStockStorage is a mock of StockStorage interface.
type MockStockStorage struct {
	ctrl     *gomock.Controller
	recorder *MockStockStorageMockRecorder
}

// MockStockStorageMockRecorder is the mock recorder for MockStockStorage.
type MockStockStorageMockRecorder struct {
	mock *MockStockStorage
}

// NewMockStockStorage creates a new mock instance.
func NewMockStockStorage(ctrl *gomock.Controller) *MockStockStorage {
	mock := &MockStockStorage{ctrl: ctrl}
	mock.recorder = &MockStockStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStockStorage) EXPECT() *MockStockStorageMockRecorder {
	return m.recorder
}

// Get mocks base method.
func (m *MockStockStorage) Get(symbol string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", symbol)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockStockStorageMockRecorder) Get(symbol interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockStockStorage)(nil).Get), symbol)
}

// Set mocks base method.
func (m *MockStockStorage) Set(symbol, stockInfo string, expireAt time.Duration) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Set", symbol, stockInfo, expireAt)
	ret0, _ := ret[0].(error)
	return ret0
}

// Set indicates an expected call of Set.
func (mr *MockStockStorageMockRecorder) Set(symbol, stockInfo, expireAt interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Set", reflect.TypeOf((*MockStockStorage)(nil).Set), symbol, stockInfo, expireAt)
}
