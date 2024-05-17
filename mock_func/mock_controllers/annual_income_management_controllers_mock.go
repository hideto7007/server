// Code generated by MockGen. DO NOT EDIT.
// Source: ./controllers/annual_income_management_controllers.go

// Package mock_controllers is a generated GoMock package.
package mock_controllers

import (
	reflect "reflect"

	gin "github.com/gin-gonic/gin"
	gomock "github.com/golang/mock/gomock"
)

// MockIncomeDataFetcher is a mock of IncomeDataFetcher interface.
type MockIncomeDataFetcher struct {
	ctrl     *gomock.Controller
	recorder *MockIncomeDataFetcherMockRecorder
}

// MockIncomeDataFetcherMockRecorder is the mock recorder for MockIncomeDataFetcher.
type MockIncomeDataFetcherMockRecorder struct {
	mock *MockIncomeDataFetcher
}

// NewMockIncomeDataFetcher creates a new mock instance.
func NewMockIncomeDataFetcher(ctrl *gomock.Controller) *MockIncomeDataFetcher {
	mock := &MockIncomeDataFetcher{ctrl: ctrl}
	mock.recorder = &MockIncomeDataFetcherMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIncomeDataFetcher) EXPECT() *MockIncomeDataFetcherMockRecorder {
	return m.recorder
}

// DeleteIncomeDataApi mocks base method.
func (m *MockIncomeDataFetcher) DeleteIncomeDataApi(c *gin.Context) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "DeleteIncomeDataApi", c)
}

// DeleteIncomeDataApi indicates an expected call of DeleteIncomeDataApi.
func (mr *MockIncomeDataFetcherMockRecorder) DeleteIncomeDataApi(c interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteIncomeDataApi", reflect.TypeOf((*MockIncomeDataFetcher)(nil).DeleteIncomeDataApi), c)
}

// GetDateRangeApi mocks base method.
func (m *MockIncomeDataFetcher) GetDateRangeApi(c *gin.Context) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "GetDateRangeApi", c)
}

// GetDateRangeApi indicates an expected call of GetDateRangeApi.
func (mr *MockIncomeDataFetcherMockRecorder) GetDateRangeApi(c interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDateRangeApi", reflect.TypeOf((*MockIncomeDataFetcher)(nil).GetDateRangeApi), c)
}

// GetIncomeDataInRangeApi mocks base method.
func (m *MockIncomeDataFetcher) GetIncomeDataInRangeApi(c *gin.Context) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "GetIncomeDataInRangeApi", c)
}

// GetIncomeDataInRangeApi indicates an expected call of GetIncomeDataInRangeApi.
func (mr *MockIncomeDataFetcherMockRecorder) GetIncomeDataInRangeApi(c interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetIncomeDataInRangeApi", reflect.TypeOf((*MockIncomeDataFetcher)(nil).GetIncomeDataInRangeApi), c)
}

// GetYearIncomeAndDeductionApi mocks base method.
func (m *MockIncomeDataFetcher) GetYearIncomeAndDeductionApi(c *gin.Context) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "GetYearIncomeAndDeductionApi", c)
}

// GetYearIncomeAndDeductionApi indicates an expected call of GetYearIncomeAndDeductionApi.
func (mr *MockIncomeDataFetcherMockRecorder) GetYearIncomeAndDeductionApi(c interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetYearIncomeAndDeductionApi", reflect.TypeOf((*MockIncomeDataFetcher)(nil).GetYearIncomeAndDeductionApi), c)
}

// InsertIncomeDataApi mocks base method.
func (m *MockIncomeDataFetcher) InsertIncomeDataApi(c *gin.Context) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "InsertIncomeDataApi", c)
}

// InsertIncomeDataApi indicates an expected call of InsertIncomeDataApi.
func (mr *MockIncomeDataFetcherMockRecorder) InsertIncomeDataApi(c interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertIncomeDataApi", reflect.TypeOf((*MockIncomeDataFetcher)(nil).InsertIncomeDataApi), c)
}

// UpdateIncomeDataApi mocks base method.
func (m *MockIncomeDataFetcher) UpdateIncomeDataApi(c *gin.Context) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "UpdateIncomeDataApi", c)
}

// UpdateIncomeDataApi indicates an expected call of UpdateIncomeDataApi.
func (mr *MockIncomeDataFetcherMockRecorder) UpdateIncomeDataApi(c interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateIncomeDataApi", reflect.TypeOf((*MockIncomeDataFetcher)(nil).UpdateIncomeDataApi), c)
}