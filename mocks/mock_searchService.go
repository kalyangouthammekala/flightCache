// Code generated by MockGen. DO NOT EDIT.
// Source: awesomeProject1/searchService (interfaces: SearchService)

// Package mocks is a generated GoMock package.
package mocks

import (
	models "awesomeProject1/models"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockSearchService is a mock of SearchService interface.
type MockSearchService struct {
	ctrl     *gomock.Controller
	recorder *MockSearchServiceMockRecorder
}

// MockSearchServiceMockRecorder is the mock recorder for MockSearchService.
type MockSearchServiceMockRecorder struct {
	mock *MockSearchService
}

// NewMockSearchService creates a new mock instance.
func NewMockSearchService(ctrl *gomock.Controller) *MockSearchService {
	mock := &MockSearchService{ctrl: ctrl}
	mock.recorder = &MockSearchServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSearchService) EXPECT() *MockSearchServiceMockRecorder {
	return m.recorder
}

// Search mocks base method.
func (m *MockSearchService) Search(arg0 *models.SearchRequest) string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Search", arg0)
	ret0, _ := ret[0].(string)
	return ret0
}

// Search indicates an expected call of Search.
func (mr *MockSearchServiceMockRecorder) Search(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Search", reflect.TypeOf((*MockSearchService)(nil).Search), arg0)
}
