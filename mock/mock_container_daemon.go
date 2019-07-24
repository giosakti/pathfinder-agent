// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/pathfinder-cm/pathfinder-agent/daemon (interfaces: ContainerDaemon)

// Package mock is a generated GoMock package.
package mock

import (
	gomock "github.com/golang/mock/gomock"
	pfmodel "github.com/pathfinder-cm/pathfinder-go-client/pfmodel"
	reflect "reflect"
)

// MockContainerDaemon is a mock of ContainerDaemon interface
type MockContainerDaemon struct {
	ctrl     *gomock.Controller
	recorder *MockContainerDaemonMockRecorder
}

// MockContainerDaemonMockRecorder is the mock recorder for MockContainerDaemon
type MockContainerDaemonMockRecorder struct {
	mock *MockContainerDaemon
}

// NewMockContainerDaemon creates a new mock instance
func NewMockContainerDaemon(ctrl *gomock.Controller) *MockContainerDaemon {
	mock := &MockContainerDaemon{ctrl: ctrl}
	mock.recorder = &MockContainerDaemonMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockContainerDaemon) EXPECT() *MockContainerDaemonMockRecorder {
	return m.recorder
}

// CreateContainer mocks base method
func (m *MockContainerDaemon) CreateContainer(arg0 pfmodel.Container) (bool, string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateContainer", arg0)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// CreateContainer indicates an expected call of CreateContainer
func (mr *MockContainerDaemonMockRecorder) CreateContainer(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateContainer", reflect.TypeOf((*MockContainerDaemon)(nil).CreateContainer), arg0)
}

// CreateContainerBootstrapFile mocks base method
func (m *MockContainerDaemon) CreateContainerBootstrapFile(arg0 pfmodel.Container) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateContainerBootstrapFile", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateContainerBootstrapFile indicates an expected call of CreateContainerBootstrapFile
func (mr *MockContainerDaemonMockRecorder) CreateContainerBootstrapFile(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateContainerBootstrapFile", reflect.TypeOf((*MockContainerDaemon)(nil).CreateContainerBootstrapFile), arg0)
}

// DeleteContainer mocks base method
func (m *MockContainerDaemon) DeleteContainer(arg0 string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteContainer", arg0)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteContainer indicates an expected call of DeleteContainer
func (mr *MockContainerDaemonMockRecorder) DeleteContainer(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteContainer", reflect.TypeOf((*MockContainerDaemon)(nil).DeleteContainer), arg0)
}

// ExecContainerBootstrap mocks base method
func (m *MockContainerDaemon) ExecContainerBootstrap(arg0 pfmodel.Container) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ExecContainerBootstrap", arg0)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ExecContainerBootstrap indicates an expected call of ExecContainerBootstrap
func (mr *MockContainerDaemonMockRecorder) ExecContainerBootstrap(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ExecContainerBootstrap", reflect.TypeOf((*MockContainerDaemon)(nil).ExecContainerBootstrap), arg0)
}

// ListContainers mocks base method
func (m *MockContainerDaemon) ListContainers() (*pfmodel.ContainerList, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListContainers")
	ret0, _ := ret[0].(*pfmodel.ContainerList)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListContainers indicates an expected call of ListContainers
func (mr *MockContainerDaemonMockRecorder) ListContainers() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListContainers", reflect.TypeOf((*MockContainerDaemon)(nil).ListContainers))
}
