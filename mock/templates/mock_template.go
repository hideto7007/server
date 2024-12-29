// Code generated by MockGen. DO NOT EDIT.
// Source: ./templates/template.go

// Package mock_templates is a generated GoMock package.
package mock_templates

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockEmailTemplateService is a mock of EmailTemplateService interface.
type MockEmailTemplateService struct {
	ctrl     *gomock.Controller
	recorder *MockEmailTemplateServiceMockRecorder
}

// MockEmailTemplateServiceMockRecorder is the mock recorder for MockEmailTemplateService.
type MockEmailTemplateServiceMockRecorder struct {
	mock *MockEmailTemplateService
}

// NewMockEmailTemplateService creates a new mock instance.
func NewMockEmailTemplateService(ctrl *gomock.Controller) *MockEmailTemplateService {
	mock := &MockEmailTemplateService{ctrl: ctrl}
	mock.recorder = &MockEmailTemplateServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockEmailTemplateService) EXPECT() *MockEmailTemplateServiceMockRecorder {
	return m.recorder
}

// DeleteSignInTemplate mocks base method.
func (m *MockEmailTemplateService) DeleteSignInTemplate(Name, UserName, DateTime string) (string, string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteSignInTemplate", Name, UserName, DateTime)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// DeleteSignInTemplate indicates an expected call of DeleteSignInTemplate.
func (mr *MockEmailTemplateServiceMockRecorder) DeleteSignInTemplate(Name, UserName, DateTime interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteSignInTemplate", reflect.TypeOf((*MockEmailTemplateService)(nil).DeleteSignInTemplate), Name, UserName, DateTime)
}

// PostSignInEditTemplate mocks base method.
func (m *MockEmailTemplateService) PostSignInEditTemplate(Update, UpdateValue, DateTime string) (string, string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PostSignInEditTemplate", Update, UpdateValue, DateTime)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// PostSignInEditTemplate indicates an expected call of PostSignInEditTemplate.
func (mr *MockEmailTemplateServiceMockRecorder) PostSignInEditTemplate(Update, UpdateValue, DateTime interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PostSignInEditTemplate", reflect.TypeOf((*MockEmailTemplateService)(nil).PostSignInEditTemplate), Update, UpdateValue, DateTime)
}

// PostSignInTemplate mocks base method.
func (m *MockEmailTemplateService) PostSignInTemplate(UserName, DateTime string) (string, string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PostSignInTemplate", UserName, DateTime)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// PostSignInTemplate indicates an expected call of PostSignInTemplate.
func (mr *MockEmailTemplateServiceMockRecorder) PostSignInTemplate(UserName, DateTime interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PostSignInTemplate", reflect.TypeOf((*MockEmailTemplateService)(nil).PostSignInTemplate), UserName, DateTime)
}

// PostSignUpTemplate mocks base method.
func (m *MockEmailTemplateService) PostSignUpTemplate(Name, UserName, DateTime string) (string, string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PostSignUpTemplate", Name, UserName, DateTime)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// PostSignUpTemplate indicates an expected call of PostSignUpTemplate.
func (mr *MockEmailTemplateServiceMockRecorder) PostSignUpTemplate(Name, UserName, DateTime interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PostSignUpTemplate", reflect.TypeOf((*MockEmailTemplateService)(nil).PostSignUpTemplate), Name, UserName, DateTime)
}

// SignOutTemplate mocks base method.
func (m *MockEmailTemplateService) SignOutTemplate(UserName, DateTime string) (string, string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SignOutTemplate", UserName, DateTime)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// SignOutTemplate indicates an expected call of SignOutTemplate.
func (mr *MockEmailTemplateServiceMockRecorder) SignOutTemplate(UserName, DateTime interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SignOutTemplate", reflect.TypeOf((*MockEmailTemplateService)(nil).SignOutTemplate), UserName, DateTime)
}

// TemporayPostSignUpTemplate mocks base method.
func (m *MockEmailTemplateService) TemporayPostSignUpTemplate(Name, ConfirmCode string) (string, string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TemporayPostSignUpTemplate", Name, ConfirmCode)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// TemporayPostSignUpTemplate indicates an expected call of TemporayPostSignUpTemplate.
func (mr *MockEmailTemplateServiceMockRecorder) TemporayPostSignUpTemplate(Name, ConfirmCode interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TemporayPostSignUpTemplate", reflect.TypeOf((*MockEmailTemplateService)(nil).TemporayPostSignUpTemplate), Name, ConfirmCode)
}
