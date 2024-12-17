// Code generated by MockGen. DO NOT EDIT.
// Source: ./config/google_config.go

// Package mock_config is a generated GoMock package.
package mock_config

import (
	http "net/http"
	reflect "reflect"

	gin "github.com/gin-gonic/gin"
	gomock "github.com/golang/mock/gomock"
	oauth2 "golang.org/x/oauth2"
)

// MockGoogleConfig is a mock of GoogleConfig interface.
type MockGoogleConfig struct {
	ctrl     *gomock.Controller
	recorder *MockGoogleConfigMockRecorder
}

// MockGoogleConfigMockRecorder is the mock recorder for MockGoogleConfig.
type MockGoogleConfigMockRecorder struct {
	mock *MockGoogleConfig
}

// NewMockGoogleConfig creates a new mock instance.
func NewMockGoogleConfig(ctrl *gomock.Controller) *MockGoogleConfig {
	mock := &MockGoogleConfig{ctrl: ctrl}
	mock.recorder = &MockGoogleConfigMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockGoogleConfig) EXPECT() *MockGoogleConfigMockRecorder {
	return m.recorder
}

// Client mocks base method.
func (m *MockGoogleConfig) Client(c *gin.Context, googleAuth *oauth2.Config, token *oauth2.Token) *http.Client {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Client", c, googleAuth, token)
	ret0, _ := ret[0].(*http.Client)
	return ret0
}

// Client indicates an expected call of Client.
func (mr *MockGoogleConfigMockRecorder) Client(c, googleAuth, token interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Client", reflect.TypeOf((*MockGoogleConfig)(nil).Client), c, googleAuth, token)
}

// Exchange mocks base method.
func (m *MockGoogleConfig) Exchange(c *gin.Context, googleAuth *oauth2.Config, code string) (*oauth2.Token, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Exchange", c, googleAuth, code)
	ret0, _ := ret[0].(*oauth2.Token)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Exchange indicates an expected call of Exchange.
func (mr *MockGoogleConfigMockRecorder) Exchange(c, googleAuth, code interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Exchange", reflect.TypeOf((*MockGoogleConfig)(nil).Exchange), c, googleAuth, code)
}

// Get mocks base method.
func (m *MockGoogleConfig) Get(client *http.Client, url string) (*http.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", client, url)
	ret0, _ := ret[0].(*http.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockGoogleConfigMockRecorder) Get(client, url interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockGoogleConfig)(nil).Get), client, url)
}

// GoogleAuthURL mocks base method.
func (m *MockGoogleConfig) GoogleAuthURL(RedirectURI string) string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GoogleAuthURL", RedirectURI)
	ret0, _ := ret[0].(string)
	return ret0
}

// GoogleAuthURL indicates an expected call of GoogleAuthURL.
func (mr *MockGoogleConfigMockRecorder) GoogleAuthURL(RedirectURI interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GoogleAuthURL", reflect.TypeOf((*MockGoogleConfig)(nil).GoogleAuthURL), RedirectURI)
}

// GoogleOauth mocks base method.
func (m *MockGoogleConfig) GoogleOauth(RedirectURI string) *oauth2.Config {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GoogleOauth", RedirectURI)
	ret0, _ := ret[0].(*oauth2.Config)
	return ret0
}

// GoogleOauth indicates an expected call of GoogleOauth.
func (mr *MockGoogleConfigMockRecorder) GoogleOauth(RedirectURI interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GoogleOauth", reflect.TypeOf((*MockGoogleConfig)(nil).GoogleOauth), RedirectURI)
}
