// Code generated by MockGen. DO NOT EDIT.
// Source: ./controllers/price_management_controllers.go

// Package mock_controllers is a generated GoMock package.
package mock_controllers

import (
	reflect "reflect"
	controllers "server/controllers"

	gin "github.com/gin-gonic/gin"
	gomock "github.com/golang/mock/gomock"
)

// MockPriceManagementFetcher is a mock of PriceManagementFetcher interface.
type MockPriceManagementFetcher struct {
	ctrl     *gomock.Controller
	recorder *MockPriceManagementFetcherMockRecorder
}

// MockPriceManagementFetcherMockRecorder is the mock recorder for MockPriceManagementFetcher.
type MockPriceManagementFetcherMockRecorder struct {
	mock *MockPriceManagementFetcher
}

// NewMockPriceManagementFetcher creates a new mock instance.
func NewMockPriceManagementFetcher(ctrl *gomock.Controller) *MockPriceManagementFetcher {
	mock := &MockPriceManagementFetcher{ctrl: ctrl}
	mock.recorder = &MockPriceManagementFetcherMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPriceManagementFetcher) EXPECT() *MockPriceManagementFetcherMockRecorder {
	return m.recorder
}

// GetPriceInfoApi mocks base method.
func (m *MockPriceManagementFetcher) GetPriceInfoApi(c *gin.Context) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "GetPriceInfoApi", c)
}

// GetPriceInfoApi indicates an expected call of GetPriceInfoApi.
func (mr *MockPriceManagementFetcherMockRecorder) GetPriceInfoApi(c interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPriceInfoApi", reflect.TypeOf((*MockPriceManagementFetcher)(nil).GetPriceInfoApi), c)
}

// PriceCalc mocks base method.
func (m *MockPriceManagementFetcher) PriceCalc(moneyReceived, bouns, fixedCost, loan, private, insurance int) controllers.PriceInfo {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PriceCalc", moneyReceived, bouns, fixedCost, loan, private, insurance)
	ret0, _ := ret[0].(controllers.PriceInfo)
	return ret0
}

// PriceCalc indicates an expected call of PriceCalc.
func (mr *MockPriceManagementFetcherMockRecorder) PriceCalc(moneyReceived, bouns, fixedCost, loan, private, insurance interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PriceCalc", reflect.TypeOf((*MockPriceManagementFetcher)(nil).PriceCalc), moneyReceived, bouns, fixedCost, loan, private, insurance)
}
