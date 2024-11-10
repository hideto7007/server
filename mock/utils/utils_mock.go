// Code generated by MockGen. DO NOT EDIT.
// Source: ./utils/utils.go

// Package mock_utils is a generated GoMock package.
package mock_utils

import (
	reflect "reflect"

	jwt "github.com/dgrijalva/jwt-go"
	gomock "github.com/golang/mock/gomock"
)

// MockUtilsFetcher is a mock of UtilsFetcher interface.
type MockUtilsFetcher struct {
	ctrl     *gomock.Controller
	recorder *MockUtilsFetcherMockRecorder
}

// MockUtilsFetcherMockRecorder is the mock recorder for MockUtilsFetcher.
type MockUtilsFetcherMockRecorder struct {
	mock *MockUtilsFetcher
}

// NewMockUtilsFetcher creates a new mock instance.
func NewMockUtilsFetcher(ctrl *gomock.Controller) *MockUtilsFetcher {
	mock := &MockUtilsFetcher{ctrl: ctrl}
	mock.recorder = &MockUtilsFetcherMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUtilsFetcher) EXPECT() *MockUtilsFetcherMockRecorder {
	return m.recorder
}

// CompareHashPassword mocks base method.
func (m *MockUtilsFetcher) CompareHashPassword(hashedPassword, requestPassword string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CompareHashPassword", hashedPassword, requestPassword)
	ret0, _ := ret[0].(error)
	return ret0
}

// CompareHashPassword indicates an expected call of CompareHashPassword.
func (mr *MockUtilsFetcherMockRecorder) CompareHashPassword(hashedPassword, requestPassword interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CompareHashPassword", reflect.TypeOf((*MockUtilsFetcher)(nil).CompareHashPassword), hashedPassword, requestPassword)
}

// EncryptPassword mocks base method.
func (m *MockUtilsFetcher) EncryptPassword(password string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EncryptPassword", password)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// EncryptPassword indicates an expected call of EncryptPassword.
func (mr *MockUtilsFetcherMockRecorder) EncryptPassword(password interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EncryptPassword", reflect.TypeOf((*MockUtilsFetcher)(nil).EncryptPassword), password)
}

// GenerateJWT mocks base method.
func (m *MockUtilsFetcher) GenerateJWT(UserId, ExpirationDate int) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenerateJWT", UserId, ExpirationDate)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GenerateJWT indicates an expected call of GenerateJWT.
func (mr *MockUtilsFetcherMockRecorder) GenerateJWT(UserId, ExpirationDate interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateJWT", reflect.TypeOf((*MockUtilsFetcher)(nil).GenerateJWT), UserId, ExpirationDate)
}

// MapClaims mocks base method.
func (m *MockUtilsFetcher) MapClaims(token *jwt.Token) (jwt.MapClaims, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MapClaims", token)
	ret0, _ := ret[0].(jwt.MapClaims)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// MapClaims indicates an expected call of MapClaims.
func (mr *MockUtilsFetcherMockRecorder) MapClaims(token interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MapClaims", reflect.TypeOf((*MockUtilsFetcher)(nil).MapClaims), token)
}

// NewToken mocks base method.
func (m *MockUtilsFetcher) NewToken(UserId, ExpirationDate int) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewToken", UserId, ExpirationDate)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// NewToken indicates an expected call of NewToken.
func (mr *MockUtilsFetcherMockRecorder) NewToken(UserId, ExpirationDate interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewToken", reflect.TypeOf((*MockUtilsFetcher)(nil).NewToken), UserId, ExpirationDate)
}

// ParseWithClaims mocks base method.
func (m *MockUtilsFetcher) ParseWithClaims(validationToken string) (*jwt.Token, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ParseWithClaims", validationToken)
	ret0, _ := ret[0].(*jwt.Token)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ParseWithClaims indicates an expected call of ParseWithClaims.
func (mr *MockUtilsFetcherMockRecorder) ParseWithClaims(validationToken interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ParseWithClaims", reflect.TypeOf((*MockUtilsFetcher)(nil).ParseWithClaims), validationToken)
}

// RefreshToken mocks base method.
func (m *MockUtilsFetcher) RefreshToken(UserId, ExpirationDate int) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RefreshToken", UserId, ExpirationDate)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RefreshToken indicates an expected call of RefreshToken.
func (mr *MockUtilsFetcherMockRecorder) RefreshToken(UserId, ExpirationDate interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RefreshToken", reflect.TypeOf((*MockUtilsFetcher)(nil).RefreshToken), UserId, ExpirationDate)
}
