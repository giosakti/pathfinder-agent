// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/pathfinder-cm/pathfinder-go-client/pfclient (interfaces: Pfclient)

// Package mock is a generated GoMock package.
package mock

import (
	gomock "github.com/golang/mock/gomock"
	pfmodel "github.com/pathfinder-cm/pathfinder-go-client/pfmodel"
	reflect "reflect"
)

// MockPfclient is a mock of Pfclient interface
type MockPfclient struct {
	ctrl     *gomock.Controller
	recorder *MockPfclientMockRecorder
}

// MockPfclientMockRecorder is the mock recorder for MockPfclient
type MockPfclientMockRecorder struct {
	mock *MockPfclient
}

// NewMockPfclient creates a new mock instance
func NewMockPfclient(ctrl *gomock.Controller) *MockPfclient {
	mock := &MockPfclient{ctrl: ctrl}
	mock.recorder = &MockPfclientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockPfclient) EXPECT() *MockPfclientMockRecorder {
	return m.recorder
}

// FetchContainersFromServer mocks base method
func (m *MockPfclient) FetchContainersFromServer(arg0 string) (*pfmodel.ContainerList, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FetchContainersFromServer", arg0)
	ret0, _ := ret[0].(*pfmodel.ContainerList)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FetchContainersFromServer indicates an expected call of FetchContainersFromServer
func (mr *MockPfclientMockRecorder) FetchContainersFromServer(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FetchContainersFromServer", reflect.TypeOf((*MockPfclient)(nil).FetchContainersFromServer), arg0)
}

// MarkContainerAsDeleted mocks base method
func (m *MockPfclient) MarkContainerAsDeleted(arg0, arg1 string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MarkContainerAsDeleted", arg0, arg1)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// MarkContainerAsDeleted indicates an expected call of MarkContainerAsDeleted
func (mr *MockPfclientMockRecorder) MarkContainerAsDeleted(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MarkContainerAsDeleted", reflect.TypeOf((*MockPfclient)(nil).MarkContainerAsDeleted), arg0, arg1)
}

// MarkContainerAsProvisionError mocks base method
func (m *MockPfclient) MarkContainerAsProvisionError(arg0, arg1 string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MarkContainerAsProvisionError", arg0, arg1)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// MarkContainerAsProvisionError indicates an expected call of MarkContainerAsProvisionError
func (mr *MockPfclientMockRecorder) MarkContainerAsProvisionError(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MarkContainerAsProvisionError", reflect.TypeOf((*MockPfclient)(nil).MarkContainerAsProvisionError), arg0, arg1)
}

// MarkContainerAsProvisioned mocks base method
func (m *MockPfclient) MarkContainerAsProvisioned(arg0, arg1 string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MarkContainerAsProvisioned", arg0, arg1)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// MarkContainerAsProvisioned indicates an expected call of MarkContainerAsProvisioned
func (mr *MockPfclientMockRecorder) MarkContainerAsProvisioned(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MarkContainerAsProvisioned", reflect.TypeOf((*MockPfclient)(nil).MarkContainerAsProvisioned), arg0, arg1)
}

// Register mocks base method
func (m *MockPfclient) Register(arg0, arg1 string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Register", arg0, arg1)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Register indicates an expected call of Register
func (mr *MockPfclientMockRecorder) Register(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Register", reflect.TypeOf((*MockPfclient)(nil).Register), arg0, arg1)
}

// StoreMetrics mocks base method
func (m *MockPfclient) StoreMetrics(arg0 *pfmodel.Metrics) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StoreMetrics", arg0)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// StoreMetrics indicates an expected call of StoreMetrics
func (mr *MockPfclientMockRecorder) StoreMetrics(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StoreMetrics", reflect.TypeOf((*MockPfclient)(nil).StoreMetrics), arg0)
}

// UpdateIpaddress mocks base method
func (m *MockPfclient) UpdateIpaddress(arg0, arg1, arg2 string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateIpaddress", arg0, arg1, arg2)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateIpaddress indicates an expected call of UpdateIpaddress
func (mr *MockPfclientMockRecorder) UpdateIpaddress(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateIpaddress", reflect.TypeOf((*MockPfclient)(nil).UpdateIpaddress), arg0, arg1, arg2)
}
