// Code generated by MockGen. DO NOT EDIT.
// Source: ./controllers/google_auth_controllers.go

// Package mock_controllers is a generated GoMock package.
package mock_controllers

import (
	reflect "reflect"

	gin "github.com/gin-gonic/gin"
	gomock "github.com/golang/mock/gomock"
)

// MockGoogleService is a mock of GoogleService interface.
type MockGoogleService struct {
	ctrl     *gomock.Controller
	recorder *MockGoogleServiceMockRecorder
}

// MockGoogleServiceMockRecorder is the mock recorder for MockGoogleService.
type MockGoogleServiceMockRecorder struct {
	mock *MockGoogleService
}

// NewMockGoogleService creates a new mock instance.
func NewMockGoogleService(ctrl *gomock.Controller) *MockGoogleService {
	mock := &MockGoogleService{ctrl: ctrl}
	mock.recorder = &MockGoogleServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockGoogleService) EXPECT() *MockGoogleServiceMockRecorder {
	return m.recorder
}

// GoogleDelete mocks base method.
func (m *MockGoogleService) GoogleDelete(c *gin.Context) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "GoogleDelete", c)
}

// GoogleDelete indicates an expected call of GoogleDelete.
func (mr *MockGoogleServiceMockRecorder) GoogleDelete(c interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GoogleDelete", reflect.TypeOf((*MockGoogleService)(nil).GoogleDelete), c)
}

// GoogleDeleteCallback mocks base method.
func (m *MockGoogleService) GoogleDeleteCallback(c *gin.Context) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "GoogleDeleteCallback", c)
}

// GoogleDeleteCallback indicates an expected call of GoogleDeleteCallback.
func (mr *MockGoogleServiceMockRecorder) GoogleDeleteCallback(c interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GoogleDeleteCallback", reflect.TypeOf((*MockGoogleService)(nil).GoogleDeleteCallback), c)
}

// GoogleSignIn mocks base method.
func (m *MockGoogleService) GoogleSignIn(c *gin.Context) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "GoogleSignIn", c)
}

// GoogleSignIn indicates an expected call of GoogleSignIn.
func (mr *MockGoogleServiceMockRecorder) GoogleSignIn(c interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GoogleSignIn", reflect.TypeOf((*MockGoogleService)(nil).GoogleSignIn), c)
}

// GoogleSignInCallback mocks base method.
func (m *MockGoogleService) GoogleSignInCallback(c *gin.Context) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "GoogleSignInCallback", c)
}

// GoogleSignInCallback indicates an expected call of GoogleSignInCallback.
func (mr *MockGoogleServiceMockRecorder) GoogleSignInCallback(c interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GoogleSignInCallback", reflect.TypeOf((*MockGoogleService)(nil).GoogleSignInCallback), c)
}

// GoogleSignUp mocks base method.
func (m *MockGoogleService) GoogleSignUp(c *gin.Context) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "GoogleSignUp", c)
}

// GoogleSignUp indicates an expected call of GoogleSignUp.
func (mr *MockGoogleServiceMockRecorder) GoogleSignUp(c interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GoogleSignUp", reflect.TypeOf((*MockGoogleService)(nil).GoogleSignUp), c)
}

// GoogleSignUpCallback mocks base method.
func (m *MockGoogleService) GoogleSignUpCallback(c *gin.Context) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "GoogleSignUpCallback", c)
}

// GoogleSignUpCallback indicates an expected call of GoogleSignUpCallback.
func (mr *MockGoogleServiceMockRecorder) GoogleSignUpCallback(c interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GoogleSignUpCallback", reflect.TypeOf((*MockGoogleService)(nil).GoogleSignUpCallback), c)
}