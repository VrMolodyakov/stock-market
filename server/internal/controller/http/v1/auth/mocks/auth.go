// Code generated by MockGen. DO NOT EDIT.
// Source: ./internal/controller/http/v1/auth/services.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"
	time "time"

	entity "github.com/VrMolodyakov/stock-market/internal/domain/entity"
	gomock "github.com/golang/mock/gomock"
)

// MockUserService is a mock of UserService interface.
type MockUserService struct {
	ctrl     *gomock.Controller
	recorder *MockUserServiceMockRecorder
}

// MockUserServiceMockRecorder is the mock recorder for MockUserService.
type MockUserServiceMockRecorder struct {
	mock *MockUserService
}

// NewMockUserService creates a new mock instance.
func NewMockUserService(ctrl *gomock.Controller) *MockUserService {
	mock := &MockUserService{ctrl: ctrl}
	mock.recorder = &MockUserServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserService) EXPECT() *MockUserServiceMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockUserService) Create(ctx context.Context, username, password string) (entity.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, username, password)
	ret0, _ := ret[0].(entity.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockUserServiceMockRecorder) Create(ctx, username, password interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockUserService)(nil).Create), ctx, username, password)
}

// Get mocks base method.
func (m *MockUserService) Get(ctx context.Context, username string) (entity.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, username)
	ret0, _ := ret[0].(entity.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockUserServiceMockRecorder) Get(ctx, username interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockUserService)(nil).Get), ctx, username)
}

// MockTokenHandler is a mock of TokenHandler interface.
type MockTokenHandler struct {
	ctrl     *gomock.Controller
	recorder *MockTokenHandlerMockRecorder
}

// MockTokenHandlerMockRecorder is the mock recorder for MockTokenHandler.
type MockTokenHandlerMockRecorder struct {
	mock *MockTokenHandler
}

// NewMockTokenHandler creates a new mock instance.
func NewMockTokenHandler(ctrl *gomock.Controller) *MockTokenHandler {
	mock := &MockTokenHandler{ctrl: ctrl}
	mock.recorder = &MockTokenHandlerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTokenHandler) EXPECT() *MockTokenHandlerMockRecorder {
	return m.recorder
}

// CreateAccessToken mocks base method.
func (m *MockTokenHandler) CreateAccessToken(ttl time.Duration, payload interface{}) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateAccessToken", ttl, payload)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateAccessToken indicates an expected call of CreateAccessToken.
func (mr *MockTokenHandlerMockRecorder) CreateAccessToken(ttl, payload interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateAccessToken", reflect.TypeOf((*MockTokenHandler)(nil).CreateAccessToken), ttl, payload)
}

// CreateRefreshToken mocks base method.
func (m *MockTokenHandler) CreateRefreshToken(ttl time.Duration, payload interface{}) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateRefreshToken", ttl, payload)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateRefreshToken indicates an expected call of CreateRefreshToken.
func (mr *MockTokenHandlerMockRecorder) CreateRefreshToken(ttl, payload interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateRefreshToken", reflect.TypeOf((*MockTokenHandler)(nil).CreateRefreshToken), ttl, payload)
}

// ValidateRefreshToken mocks base method.
func (m *MockTokenHandler) ValidateRefreshToken(token string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ValidateRefreshToken", token)
	ret0, _ := ret[0].(error)
	return ret0
}

// ValidateRefreshToken indicates an expected call of ValidateRefreshToken.
func (mr *MockTokenHandlerMockRecorder) ValidateRefreshToken(token interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidateRefreshToken", reflect.TypeOf((*MockTokenHandler)(nil).ValidateRefreshToken), token)
}

// MockTokenService is a mock of TokenService interface.
type MockTokenService struct {
	ctrl     *gomock.Controller
	recorder *MockTokenServiceMockRecorder
}

// MockTokenServiceMockRecorder is the mock recorder for MockTokenService.
type MockTokenServiceMockRecorder struct {
	mock *MockTokenService
}

// NewMockTokenService creates a new mock instance.
func NewMockTokenService(ctrl *gomock.Controller) *MockTokenService {
	mock := &MockTokenService{ctrl: ctrl}
	mock.recorder = &MockTokenServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTokenService) EXPECT() *MockTokenServiceMockRecorder {
	return m.recorder
}

// Find mocks base method.
func (m *MockTokenService) Find(refreshToken string) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Find", refreshToken)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Find indicates an expected call of Find.
func (mr *MockTokenServiceMockRecorder) Find(refreshToken interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Find", reflect.TypeOf((*MockTokenService)(nil).Find), refreshToken)
}

// Remove mocks base method.
func (m *MockTokenService) Remove(refreshToken string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Remove", refreshToken)
	ret0, _ := ret[0].(error)
	return ret0
}

// Remove indicates an expected call of Remove.
func (mr *MockTokenServiceMockRecorder) Remove(refreshToken interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Remove", reflect.TypeOf((*MockTokenService)(nil).Remove), refreshToken)
}

// Save mocks base method.
func (m *MockTokenService) Save(refreshToken string, userId int, expireAt time.Duration) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save", refreshToken, userId, expireAt)
	ret0, _ := ret[0].(error)
	return ret0
}

// Save indicates an expected call of Save.
func (mr *MockTokenServiceMockRecorder) Save(refreshToken, userId, expireAt interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockTokenService)(nil).Save), refreshToken, userId, expireAt)
}