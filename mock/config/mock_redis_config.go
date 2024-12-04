// Code generated by MockGen. DO NOT EDIT.
// Source: ./config/redis_config.go

// Package mock_config is a generated GoMock package.
package mock_config

import (
	reflect "reflect"
	time "time"

	gomock "github.com/golang/mock/gomock"
	redis "github.com/redis/go-redis/v9"
)

// MockRedisService is a mock of RedisService interface.
type MockRedisService struct {
	ctrl     *gomock.Controller
	recorder *MockRedisServiceMockRecorder
}

// MockRedisServiceMockRecorder is the mock recorder for MockRedisService.
type MockRedisServiceMockRecorder struct {
	mock *MockRedisService
}

// NewMockRedisService creates a new mock instance.
func NewMockRedisService(ctrl *gomock.Controller) *MockRedisService {
	mock := &MockRedisService{ctrl: ctrl}
	mock.recorder = &MockRedisServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRedisService) EXPECT() *MockRedisServiceMockRecorder {
	return m.recorder
}

// InitRedisClient mocks base method.
func (m *MockRedisService) InitRedisClient() *redis.Client {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InitRedisClient")
	ret0, _ := ret[0].(*redis.Client)
	return ret0
}

// InitRedisClient indicates an expected call of InitRedisClient.
func (mr *MockRedisServiceMockRecorder) InitRedisClient() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InitRedisClient", reflect.TypeOf((*MockRedisService)(nil).InitRedisClient))
}

// RedisDel mocks base method.
func (m *MockRedisService) RedisDel(key string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RedisDel", key)
	ret0, _ := ret[0].(error)
	return ret0
}

// RedisDel indicates an expected call of RedisDel.
func (mr *MockRedisServiceMockRecorder) RedisDel(key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RedisDel", reflect.TypeOf((*MockRedisService)(nil).RedisDel), key)
}

// RedisGet mocks base method.
func (m *MockRedisService) RedisGet(key string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RedisGet", key)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RedisGet indicates an expected call of RedisGet.
func (mr *MockRedisServiceMockRecorder) RedisGet(key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RedisGet", reflect.TypeOf((*MockRedisService)(nil).RedisGet), key)
}

// RedisSet mocks base method.
func (m *MockRedisService) RedisSet(key string, value interface{}, duration time.Duration) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RedisSet", key, value, duration)
	ret0, _ := ret[0].(error)
	return ret0
}

// RedisSet indicates an expected call of RedisSet.
func (mr *MockRedisServiceMockRecorder) RedisSet(key, value, duration interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RedisSet", reflect.TypeOf((*MockRedisService)(nil).RedisSet), key, value, duration)
}