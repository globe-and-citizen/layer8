// Code generated by MockGen. DO NOT EDIT.
// Source: server/internals/service/service.go

// Package mocks is a generated GoMock package.
package mocks

import (
	entities "globe-and-citizen/layer8/server/entities"
	models "globe-and-citizen/layer8/server/models"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	oauth2 "golang.org/x/oauth2"
)

// MockServiceInterface is a mock of ServiceInterface interface.
type MockServiceInterface struct {
	ctrl     *gomock.Controller
	recorder *MockServiceInterfaceMockRecorder
}

// MockServiceInterfaceMockRecorder is the mock recorder for MockServiceInterface.
type MockServiceInterfaceMockRecorder struct {
	mock *MockServiceInterface
}

// NewMockServiceInterface creates a new mock instance.
func NewMockServiceInterface(ctrl *gomock.Controller) *MockServiceInterface {
	mock := &MockServiceInterface{ctrl: ctrl}
	mock.recorder = &MockServiceInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockServiceInterface) EXPECT() *MockServiceInterfaceMockRecorder {
	return m.recorder
}

// AccessResourcesWithToken mocks base method.
func (m *MockServiceInterface) AccessResourcesWithToken(token string) (map[string]interface{}, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AccessResourcesWithToken", token)
	ret0, _ := ret[0].(map[string]interface{})
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AccessResourcesWithToken indicates an expected call of AccessResourcesWithToken.
func (mr *MockServiceInterfaceMockRecorder) AccessResourcesWithToken(token interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AccessResourcesWithToken", reflect.TypeOf((*MockServiceInterface)(nil).AccessResourcesWithToken), token)
}

// AddTestClient mocks base method.
func (m *MockServiceInterface) AddTestClient() (*models.Client, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddTestClient")
	ret0, _ := ret[0].(*models.Client)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddTestClient indicates an expected call of AddTestClient.
func (mr *MockServiceInterfaceMockRecorder) AddTestClient() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddTestClient", reflect.TypeOf((*MockServiceInterface)(nil).AddTestClient))
}

// ExchangeCodeForToken mocks base method.
func (m *MockServiceInterface) ExchangeCodeForToken(config *oauth2.Config, code string) (*oauth2.Token, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ExchangeCodeForToken", config, code)
	ret0, _ := ret[0].(*oauth2.Token)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ExchangeCodeForToken indicates an expected call of ExchangeCodeForToken.
func (mr *MockServiceInterfaceMockRecorder) ExchangeCodeForToken(config, code interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ExchangeCodeForToken", reflect.TypeOf((*MockServiceInterface)(nil).ExchangeCodeForToken), config, code)
}

// GenerateAuthorizationURL mocks base method.
func (m *MockServiceInterface) GenerateAuthorizationURL(config *oauth2.Config, userID int64, headerMap map[string]string) (*entities.AuthURL, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenerateAuthorizationURL", config, userID, headerMap)
	ret0, _ := ret[0].(*entities.AuthURL)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GenerateAuthorizationURL indicates an expected call of GenerateAuthorizationURL.
func (mr *MockServiceInterfaceMockRecorder) GenerateAuthorizationURL(config, userID, headerMap interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateAuthorizationURL", reflect.TypeOf((*MockServiceInterface)(nil).GenerateAuthorizationURL), config, userID, headerMap)
}

// GetClient mocks base method.
func (m *MockServiceInterface) GetClient(id string) (*models.Client, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetClient", id)
	ret0, _ := ret[0].(*models.Client)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetClient indicates an expected call of GetClient.
func (mr *MockServiceInterfaceMockRecorder) GetClient(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetClient", reflect.TypeOf((*MockServiceInterface)(nil).GetClient), id)
}

// GetUserByToken mocks base method.
func (m *MockServiceInterface) GetUserByToken(token string) (*models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByToken", token)
	ret0, _ := ret[0].(*models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByToken indicates an expected call of GetUserByToken.
func (mr *MockServiceInterfaceMockRecorder) GetUserByToken(token interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByToken", reflect.TypeOf((*MockServiceInterface)(nil).GetUserByToken), token)
}

// LoginUser mocks base method.
func (m *MockServiceInterface) LoginUser(username, password string) (map[string]interface{}, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LoginUser", username, password)
	ret0, _ := ret[0].(map[string]interface{})
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LoginUser indicates an expected call of LoginUser.
func (mr *MockServiceInterfaceMockRecorder) LoginUser(username, password interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LoginUser", reflect.TypeOf((*MockServiceInterface)(nil).LoginUser), username, password)
}
