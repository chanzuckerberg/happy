// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/chanzuckerberg/happy/pkg/stack_mgr (interfaces: StackServiceIface)

// Package stack_mgr is a generated GoMock package.
package stack_mgr

import (
	reflect "reflect"

	config "github.com/chanzuckerberg/happy/pkg/config"
	workspace_repo "github.com/chanzuckerberg/happy/pkg/workspace_repo"
	gomock "github.com/golang/mock/gomock"
)

// MockStackServiceIface is a mock of StackServiceIface interface.
type MockStackServiceIface struct {
	ctrl     *gomock.Controller
	recorder *MockStackServiceIfaceMockRecorder
}

// MockStackServiceIfaceMockRecorder is the mock recorder for MockStackServiceIface.
type MockStackServiceIfaceMockRecorder struct {
	mock *MockStackServiceIface
}

// NewMockStackServiceIface creates a new mock instance.
func NewMockStackServiceIface(ctrl *gomock.Controller) *MockStackServiceIface {
	mock := &MockStackServiceIface{ctrl: ctrl}
	mock.recorder = &MockStackServiceIfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStackServiceIface) EXPECT() *MockStackServiceIfaceMockRecorder {
	return m.recorder
}

// Add mocks base method.
func (m *MockStackServiceIface) Add(arg0 string) (*Stack, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Add", arg0)
	ret0, _ := ret[0].(*Stack)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Add indicates an expected call of Add.
func (mr *MockStackServiceIfaceMockRecorder) Add(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Add", reflect.TypeOf((*MockStackServiceIface)(nil).Add), arg0)
}

// GetConfig mocks base method.
func (m *MockStackServiceIface) GetConfig() config.HappyConfig {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetConfig")
	ret0, _ := ret[0].(config.HappyConfig)
	return ret0
}

// GetConfig indicates an expected call of GetConfig.
func (mr *MockStackServiceIfaceMockRecorder) GetConfig() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetConfig", reflect.TypeOf((*MockStackServiceIface)(nil).GetConfig))
}

// GetStackWorkspace mocks base method.
func (m *MockStackServiceIface) GetStackWorkspace(arg0 string) (workspace_repo.Workspace, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetStackWorkspace", arg0)
	ret0, _ := ret[0].(workspace_repo.Workspace)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetStackWorkspace indicates an expected call of GetStackWorkspace.
func (mr *MockStackServiceIfaceMockRecorder) GetStackWorkspace(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetStackWorkspace", reflect.TypeOf((*MockStackServiceIface)(nil).GetStackWorkspace), arg0)
}

// GetStacks mocks base method.
func (m *MockStackServiceIface) GetStacks() (map[string]*Stack, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetStacks")
	ret0, _ := ret[0].(map[string]*Stack)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetStacks indicates an expected call of GetStacks.
func (mr *MockStackServiceIfaceMockRecorder) GetStacks() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetStacks", reflect.TypeOf((*MockStackServiceIface)(nil).GetStacks))
}

// NewStackMeta mocks base method.
func (m *MockStackServiceIface) NewStackMeta(arg0 string) *StackMeta {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewStackMeta", arg0)
	ret0, _ := ret[0].(*StackMeta)
	return ret0
}

// NewStackMeta indicates an expected call of NewStackMeta.
func (mr *MockStackServiceIfaceMockRecorder) NewStackMeta(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewStackMeta", reflect.TypeOf((*MockStackServiceIface)(nil).NewStackMeta), arg0)
}

// Remove mocks base method.
func (m *MockStackServiceIface) Remove(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Remove", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Remove indicates an expected call of Remove.
func (mr *MockStackServiceIfaceMockRecorder) Remove(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Remove", reflect.TypeOf((*MockStackServiceIface)(nil).Remove), arg0)
}
