// Automatically generated by MockGen. DO NOT EDIT!
// Source: pkg/handler-launcher-com/cmd/v1/cmd.pb.go

package v1

import (
	gomock "github.com/golang/mock/gomock"
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Mock of CmdClient interface
type MockCmdClient struct {
	ctrl     *gomock.Controller
	recorder *_MockCmdClientRecorder
}

// Recorder for MockCmdClient (not exported)
type _MockCmdClientRecorder struct {
	mock *MockCmdClient
}

func NewMockCmdClient(ctrl *gomock.Controller) *MockCmdClient {
	mock := &MockCmdClient{ctrl: ctrl}
	mock.recorder = &_MockCmdClientRecorder{mock}
	return mock
}

func (_m *MockCmdClient) EXPECT() *_MockCmdClientRecorder {
	return _m.recorder
}

func (_m *MockCmdClient) SyncVirtualMachine(ctx context.Context, in *VMIRequest, opts ...grpc.CallOption) (*Response, error) {
	_s := []interface{}{ctx, in}
	for _, _x := range opts {
		_s = append(_s, _x)
	}
	ret := _m.ctrl.Call(_m, "SyncVirtualMachine", _s...)
	ret0, _ := ret[0].(*Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockCmdClientRecorder) SyncVirtualMachine(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	_s := append([]interface{}{arg0, arg1}, arg2...)
	return _mr.mock.ctrl.RecordCall(_mr.mock, "SyncVirtualMachine", _s...)
}

func (_m *MockCmdClient) PauseVirtualMachine(ctx context.Context, in *VMIRequest, opts ...grpc.CallOption) (*Response, error) {
	_s := []interface{}{ctx, in}
	for _, _x := range opts {
		_s = append(_s, _x)
	}
	ret := _m.ctrl.Call(_m, "PauseVirtualMachine", _s...)
	ret0, _ := ret[0].(*Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockCmdClientRecorder) PauseVirtualMachine(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	_s := append([]interface{}{arg0, arg1}, arg2...)
	return _mr.mock.ctrl.RecordCall(_mr.mock, "PauseVirtualMachine", _s...)
}

func (_m *MockCmdClient) UnpauseVirtualMachine(ctx context.Context, in *VMIRequest, opts ...grpc.CallOption) (*Response, error) {
	_s := []interface{}{ctx, in}
	for _, _x := range opts {
		_s = append(_s, _x)
	}
	ret := _m.ctrl.Call(_m, "UnpauseVirtualMachine", _s...)
	ret0, _ := ret[0].(*Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockCmdClientRecorder) UnpauseVirtualMachine(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	_s := append([]interface{}{arg0, arg1}, arg2...)
	return _mr.mock.ctrl.RecordCall(_mr.mock, "UnpauseVirtualMachine", _s...)
}

func (_m *MockCmdClient) FreezeVirtualMachine(ctx context.Context, in *FreezeRequest, opts ...grpc.CallOption) (*Response, error) {
	_s := []interface{}{ctx, in}
	for _, _x := range opts {
		_s = append(_s, _x)
	}
	ret := _m.ctrl.Call(_m, "FreezeVirtualMachine", _s...)
	ret0, _ := ret[0].(*Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockCmdClientRecorder) FreezeVirtualMachine(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	_s := append([]interface{}{arg0, arg1}, arg2...)
	return _mr.mock.ctrl.RecordCall(_mr.mock, "FreezeVirtualMachine", _s...)
}

func (_m *MockCmdClient) UnfreezeVirtualMachine(ctx context.Context, in *VMIRequest, opts ...grpc.CallOption) (*Response, error) {
	_s := []interface{}{ctx, in}
	for _, _x := range opts {
		_s = append(_s, _x)
	}
	ret := _m.ctrl.Call(_m, "UnfreezeVirtualMachine", _s...)
	ret0, _ := ret[0].(*Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockCmdClientRecorder) UnfreezeVirtualMachine(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	_s := append([]interface{}{arg0, arg1}, arg2...)
	return _mr.mock.ctrl.RecordCall(_mr.mock, "UnfreezeVirtualMachine", _s...)
}

func (_m *MockCmdClient) SoftRebootVirtualMachine(ctx context.Context, in *VMIRequest, opts ...grpc.CallOption) (*Response, error) {
	_s := []interface{}{ctx, in}
	for _, _x := range opts {
		_s = append(_s, _x)
	}
	ret := _m.ctrl.Call(_m, "SoftRebootVirtualMachine", _s...)
	ret0, _ := ret[0].(*Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockCmdClientRecorder) SoftRebootVirtualMachine(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	_s := append([]interface{}{arg0, arg1}, arg2...)
	return _mr.mock.ctrl.RecordCall(_mr.mock, "SoftRebootVirtualMachine", _s...)
}

func (_m *MockCmdClient) ShutdownVirtualMachine(ctx context.Context, in *VMIRequest, opts ...grpc.CallOption) (*Response, error) {
	_s := []interface{}{ctx, in}
	for _, _x := range opts {
		_s = append(_s, _x)
	}
	ret := _m.ctrl.Call(_m, "ShutdownVirtualMachine", _s...)
	ret0, _ := ret[0].(*Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockCmdClientRecorder) ShutdownVirtualMachine(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	_s := append([]interface{}{arg0, arg1}, arg2...)
	return _mr.mock.ctrl.RecordCall(_mr.mock, "ShutdownVirtualMachine", _s...)
}

func (_m *MockCmdClient) KillVirtualMachine(ctx context.Context, in *VMIRequest, opts ...grpc.CallOption) (*Response, error) {
	_s := []interface{}{ctx, in}
	for _, _x := range opts {
		_s = append(_s, _x)
	}
	ret := _m.ctrl.Call(_m, "KillVirtualMachine", _s...)
	ret0, _ := ret[0].(*Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockCmdClientRecorder) KillVirtualMachine(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	_s := append([]interface{}{arg0, arg1}, arg2...)
	return _mr.mock.ctrl.RecordCall(_mr.mock, "KillVirtualMachine", _s...)
}

func (_m *MockCmdClient) DeleteVirtualMachine(ctx context.Context, in *VMIRequest, opts ...grpc.CallOption) (*Response, error) {
	_s := []interface{}{ctx, in}
	for _, _x := range opts {
		_s = append(_s, _x)
	}
	ret := _m.ctrl.Call(_m, "DeleteVirtualMachine", _s...)
	ret0, _ := ret[0].(*Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockCmdClientRecorder) DeleteVirtualMachine(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	_s := append([]interface{}{arg0, arg1}, arg2...)
	return _mr.mock.ctrl.RecordCall(_mr.mock, "DeleteVirtualMachine", _s...)
}

func (_m *MockCmdClient) MigrateVirtualMachine(ctx context.Context, in *MigrationRequest, opts ...grpc.CallOption) (*Response, error) {
	_s := []interface{}{ctx, in}
	for _, _x := range opts {
		_s = append(_s, _x)
	}
	ret := _m.ctrl.Call(_m, "MigrateVirtualMachine", _s...)
	ret0, _ := ret[0].(*Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockCmdClientRecorder) MigrateVirtualMachine(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	_s := append([]interface{}{arg0, arg1}, arg2...)
	return _mr.mock.ctrl.RecordCall(_mr.mock, "MigrateVirtualMachine", _s...)
}

func (_m *MockCmdClient) SyncMigrationTarget(ctx context.Context, in *VMIRequest, opts ...grpc.CallOption) (*Response, error) {
	_s := []interface{}{ctx, in}
	for _, _x := range opts {
		_s = append(_s, _x)
	}
	ret := _m.ctrl.Call(_m, "SyncMigrationTarget", _s...)
	ret0, _ := ret[0].(*Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockCmdClientRecorder) SyncMigrationTarget(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	_s := append([]interface{}{arg0, arg1}, arg2...)
	return _mr.mock.ctrl.RecordCall(_mr.mock, "SyncMigrationTarget", _s...)
}

func (_m *MockCmdClient) CancelVirtualMachineMigration(ctx context.Context, in *VMIRequest, opts ...grpc.CallOption) (*Response, error) {
	_s := []interface{}{ctx, in}
	for _, _x := range opts {
		_s = append(_s, _x)
	}
	ret := _m.ctrl.Call(_m, "CancelVirtualMachineMigration", _s...)
	ret0, _ := ret[0].(*Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockCmdClientRecorder) CancelVirtualMachineMigration(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	_s := append([]interface{}{arg0, arg1}, arg2...)
	return _mr.mock.ctrl.RecordCall(_mr.mock, "CancelVirtualMachineMigration", _s...)
}

func (_m *MockCmdClient) SignalTargetPodCleanup(ctx context.Context, in *VMIRequest, opts ...grpc.CallOption) (*Response, error) {
	_s := []interface{}{ctx, in}
	for _, _x := range opts {
		_s = append(_s, _x)
	}
	ret := _m.ctrl.Call(_m, "SignalTargetPodCleanup", _s...)
	ret0, _ := ret[0].(*Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockCmdClientRecorder) SignalTargetPodCleanup(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	_s := append([]interface{}{arg0, arg1}, arg2...)
	return _mr.mock.ctrl.RecordCall(_mr.mock, "SignalTargetPodCleanup", _s...)
}

func (_m *MockCmdClient) FinalizeVirtualMachineMigration(ctx context.Context, in *VMIRequest, opts ...grpc.CallOption) (*Response, error) {
	_s := []interface{}{ctx, in}
	for _, _x := range opts {
		_s = append(_s, _x)
	}
	ret := _m.ctrl.Call(_m, "FinalizeVirtualMachineMigration", _s...)
	ret0, _ := ret[0].(*Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockCmdClientRecorder) FinalizeVirtualMachineMigration(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	_s := append([]interface{}{arg0, arg1}, arg2...)
	return _mr.mock.ctrl.RecordCall(_mr.mock, "FinalizeVirtualMachineMigration", _s...)
}

func (_m *MockCmdClient) HotplugHostDevices(ctx context.Context, in *VMIRequest, opts ...grpc.CallOption) (*Response, error) {
	_s := []interface{}{ctx, in}
	for _, _x := range opts {
		_s = append(_s, _x)
	}
	ret := _m.ctrl.Call(_m, "HotplugHostDevices", _s...)
	ret0, _ := ret[0].(*Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockCmdClientRecorder) HotplugHostDevices(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	_s := append([]interface{}{arg0, arg1}, arg2...)
	return _mr.mock.ctrl.RecordCall(_mr.mock, "HotplugHostDevices", _s...)
}

func (_m *MockCmdClient) GetDomain(ctx context.Context, in *EmptyRequest, opts ...grpc.CallOption) (*DomainResponse, error) {
	_s := []interface{}{ctx, in}
	for _, _x := range opts {
		_s = append(_s, _x)
	}
	ret := _m.ctrl.Call(_m, "GetDomain", _s...)
	ret0, _ := ret[0].(*DomainResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockCmdClientRecorder) GetDomain(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	_s := append([]interface{}{arg0, arg1}, arg2...)
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetDomain", _s...)
}

func (_m *MockCmdClient) GetDomainStats(ctx context.Context, in *EmptyRequest, opts ...grpc.CallOption) (*DomainStatsResponse, error) {
	_s := []interface{}{ctx, in}
	for _, _x := range opts {
		_s = append(_s, _x)
	}
	ret := _m.ctrl.Call(_m, "GetDomainStats", _s...)
	ret0, _ := ret[0].(*DomainStatsResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockCmdClientRecorder) GetDomainStats(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	_s := append([]interface{}{arg0, arg1}, arg2...)
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetDomainStats", _s...)
}

func (_m *MockCmdClient) GetGuestInfo(ctx context.Context, in *EmptyRequest, opts ...grpc.CallOption) (*GuestInfoResponse, error) {
	_s := []interface{}{ctx, in}
	for _, _x := range opts {
		_s = append(_s, _x)
	}
	ret := _m.ctrl.Call(_m, "GetGuestInfo", _s...)
	ret0, _ := ret[0].(*GuestInfoResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockCmdClientRecorder) GetGuestInfo(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	_s := append([]interface{}{arg0, arg1}, arg2...)
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetGuestInfo", _s...)
}

func (_m *MockCmdClient) GetUsers(ctx context.Context, in *EmptyRequest, opts ...grpc.CallOption) (*GuestUserListResponse, error) {
	_s := []interface{}{ctx, in}
	for _, _x := range opts {
		_s = append(_s, _x)
	}
	ret := _m.ctrl.Call(_m, "GetUsers", _s...)
	ret0, _ := ret[0].(*GuestUserListResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockCmdClientRecorder) GetUsers(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	_s := append([]interface{}{arg0, arg1}, arg2...)
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetUsers", _s...)
}

func (_m *MockCmdClient) GetFilesystems(ctx context.Context, in *EmptyRequest, opts ...grpc.CallOption) (*GuestFilesystemsResponse, error) {
	_s := []interface{}{ctx, in}
	for _, _x := range opts {
		_s = append(_s, _x)
	}
	ret := _m.ctrl.Call(_m, "GetFilesystems", _s...)
	ret0, _ := ret[0].(*GuestFilesystemsResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockCmdClientRecorder) GetFilesystems(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	_s := append([]interface{}{arg0, arg1}, arg2...)
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetFilesystems", _s...)
}

func (_m *MockCmdClient) Ping(ctx context.Context, in *EmptyRequest, opts ...grpc.CallOption) (*Response, error) {
	_s := []interface{}{ctx, in}
	for _, _x := range opts {
		_s = append(_s, _x)
	}
	ret := _m.ctrl.Call(_m, "Ping", _s...)
	ret0, _ := ret[0].(*Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockCmdClientRecorder) Ping(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	_s := append([]interface{}{arg0, arg1}, arg2...)
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Ping", _s...)
}

func (_m *MockCmdClient) Exec(ctx context.Context, in *ExecRequest, opts ...grpc.CallOption) (*ExecResponse, error) {
	_s := []interface{}{ctx, in}
	for _, _x := range opts {
		_s = append(_s, _x)
	}
	ret := _m.ctrl.Call(_m, "Exec", _s...)
	ret0, _ := ret[0].(*ExecResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockCmdClientRecorder) Exec(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	_s := append([]interface{}{arg0, arg1}, arg2...)
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Exec", _s...)
}

func (_m *MockCmdClient) GuestPing(ctx context.Context, in *GuestPingRequest, opts ...grpc.CallOption) (*GuestPingResponse, error) {
	_s := []interface{}{ctx, in}
	for _, _x := range opts {
		_s = append(_s, _x)
	}
	ret := _m.ctrl.Call(_m, "GuestPing", _s...)
	ret0, _ := ret[0].(*GuestPingResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockCmdClientRecorder) GuestPing(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	_s := append([]interface{}{arg0, arg1}, arg2...)
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GuestPing", _s...)
}

func (_m *MockCmdClient) VirtualMachineMemoryDump(ctx context.Context, in *MemoryDumpRequest, opts ...grpc.CallOption) (*Response, error) {
	_s := []interface{}{ctx, in}
	for _, _x := range opts {
		_s = append(_s, _x)
	}
	ret := _m.ctrl.Call(_m, "VirtualMachineMemoryDump", _s...)
	ret0, _ := ret[0].(*Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockCmdClientRecorder) VirtualMachineMemoryDump(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	_s := append([]interface{}{arg0, arg1}, arg2...)
	return _mr.mock.ctrl.RecordCall(_mr.mock, "VirtualMachineMemoryDump", _s...)
}

func (_m *MockCmdClient) GetVmmVersion(ctx context.Context, in *EmptyRequest, opts ...grpc.CallOption) (*VmmVersionResponse, error) {
	_s := []interface{}{ctx, in}
	for _, _x := range opts {
		_s = append(_s, _x)
	}
	ret := _m.ctrl.Call(_m, "GetVmmVersion", _s...)
	ret0, _ := ret[0].(*VmmVersionResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockCmdClientRecorder) GetVmmVersion(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	_s := append([]interface{}{arg0, arg1}, arg2...)
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetVmmVersion", _s...)
}

// Mock of CmdServer interface
type MockCmdServer struct {
	ctrl     *gomock.Controller
	recorder *_MockCmdServerRecorder
}

// Recorder for MockCmdServer (not exported)
type _MockCmdServerRecorder struct {
	mock *MockCmdServer
}

func NewMockCmdServer(ctrl *gomock.Controller) *MockCmdServer {
	mock := &MockCmdServer{ctrl: ctrl}
	mock.recorder = &_MockCmdServerRecorder{mock}
	return mock
}

func (_m *MockCmdServer) EXPECT() *_MockCmdServerRecorder {
	return _m.recorder
}

func (_m *MockCmdServer) SyncVirtualMachine(_param0 context.Context, _param1 *VMIRequest) (*Response, error) {
	ret := _m.ctrl.Call(_m, "SyncVirtualMachine", _param0, _param1)
	ret0, _ := ret[0].(*Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockCmdServerRecorder) SyncVirtualMachine(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "SyncVirtualMachine", arg0, arg1)
}

func (_m *MockCmdServer) PauseVirtualMachine(_param0 context.Context, _param1 *VMIRequest) (*Response, error) {
	ret := _m.ctrl.Call(_m, "PauseVirtualMachine", _param0, _param1)
	ret0, _ := ret[0].(*Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockCmdServerRecorder) PauseVirtualMachine(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "PauseVirtualMachine", arg0, arg1)
}

func (_m *MockCmdServer) UnpauseVirtualMachine(_param0 context.Context, _param1 *VMIRequest) (*Response, error) {
	ret := _m.ctrl.Call(_m, "UnpauseVirtualMachine", _param0, _param1)
	ret0, _ := ret[0].(*Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockCmdServerRecorder) UnpauseVirtualMachine(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "UnpauseVirtualMachine", arg0, arg1)
}

func (_m *MockCmdServer) FreezeVirtualMachine(_param0 context.Context, _param1 *FreezeRequest) (*Response, error) {
	ret := _m.ctrl.Call(_m, "FreezeVirtualMachine", _param0, _param1)
	ret0, _ := ret[0].(*Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockCmdServerRecorder) FreezeVirtualMachine(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "FreezeVirtualMachine", arg0, arg1)
}

func (_m *MockCmdServer) UnfreezeVirtualMachine(_param0 context.Context, _param1 *VMIRequest) (*Response, error) {
	ret := _m.ctrl.Call(_m, "UnfreezeVirtualMachine", _param0, _param1)
	ret0, _ := ret[0].(*Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockCmdServerRecorder) UnfreezeVirtualMachine(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "UnfreezeVirtualMachine", arg0, arg1)
}

func (_m *MockCmdServer) SoftRebootVirtualMachine(_param0 context.Context, _param1 *VMIRequest) (*Response, error) {
	ret := _m.ctrl.Call(_m, "SoftRebootVirtualMachine", _param0, _param1)
	ret0, _ := ret[0].(*Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockCmdServerRecorder) SoftRebootVirtualMachine(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "SoftRebootVirtualMachine", arg0, arg1)
}

func (_m *MockCmdServer) ShutdownVirtualMachine(_param0 context.Context, _param1 *VMIRequest) (*Response, error) {
	ret := _m.ctrl.Call(_m, "ShutdownVirtualMachine", _param0, _param1)
	ret0, _ := ret[0].(*Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockCmdServerRecorder) ShutdownVirtualMachine(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "ShutdownVirtualMachine", arg0, arg1)
}

func (_m *MockCmdServer) KillVirtualMachine(_param0 context.Context, _param1 *VMIRequest) (*Response, error) {
	ret := _m.ctrl.Call(_m, "KillVirtualMachine", _param0, _param1)
	ret0, _ := ret[0].(*Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockCmdServerRecorder) KillVirtualMachine(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "KillVirtualMachine", arg0, arg1)
}

func (_m *MockCmdServer) DeleteVirtualMachine(_param0 context.Context, _param1 *VMIRequest) (*Response, error) {
	ret := _m.ctrl.Call(_m, "DeleteVirtualMachine", _param0, _param1)
	ret0, _ := ret[0].(*Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockCmdServerRecorder) DeleteVirtualMachine(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "DeleteVirtualMachine", arg0, arg1)
}

func (_m *MockCmdServer) MigrateVirtualMachine(_param0 context.Context, _param1 *MigrationRequest) (*Response, error) {
	ret := _m.ctrl.Call(_m, "MigrateVirtualMachine", _param0, _param1)
	ret0, _ := ret[0].(*Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockCmdServerRecorder) MigrateVirtualMachine(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "MigrateVirtualMachine", arg0, arg1)
}

func (_m *MockCmdServer) SyncMigrationTarget(_param0 context.Context, _param1 *VMIRequest) (*Response, error) {
	ret := _m.ctrl.Call(_m, "SyncMigrationTarget", _param0, _param1)
	ret0, _ := ret[0].(*Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockCmdServerRecorder) SyncMigrationTarget(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "SyncMigrationTarget", arg0, arg1)
}

func (_m *MockCmdServer) CancelVirtualMachineMigration(_param0 context.Context, _param1 *VMIRequest) (*Response, error) {
	ret := _m.ctrl.Call(_m, "CancelVirtualMachineMigration", _param0, _param1)
	ret0, _ := ret[0].(*Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockCmdServerRecorder) CancelVirtualMachineMigration(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "CancelVirtualMachineMigration", arg0, arg1)
}

func (_m *MockCmdServer) SignalTargetPodCleanup(_param0 context.Context, _param1 *VMIRequest) (*Response, error) {
	ret := _m.ctrl.Call(_m, "SignalTargetPodCleanup", _param0, _param1)
	ret0, _ := ret[0].(*Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockCmdServerRecorder) SignalTargetPodCleanup(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "SignalTargetPodCleanup", arg0, arg1)
}

func (_m *MockCmdServer) FinalizeVirtualMachineMigration(_param0 context.Context, _param1 *VMIRequest) (*Response, error) {
	ret := _m.ctrl.Call(_m, "FinalizeVirtualMachineMigration", _param0, _param1)
	ret0, _ := ret[0].(*Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockCmdServerRecorder) FinalizeVirtualMachineMigration(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "FinalizeVirtualMachineMigration", arg0, arg1)
}

func (_m *MockCmdServer) HotplugHostDevices(_param0 context.Context, _param1 *VMIRequest) (*Response, error) {
	ret := _m.ctrl.Call(_m, "HotplugHostDevices", _param0, _param1)
	ret0, _ := ret[0].(*Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockCmdServerRecorder) HotplugHostDevices(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "HotplugHostDevices", arg0, arg1)
}

func (_m *MockCmdServer) GetDomain(_param0 context.Context, _param1 *EmptyRequest) (*DomainResponse, error) {
	ret := _m.ctrl.Call(_m, "GetDomain", _param0, _param1)
	ret0, _ := ret[0].(*DomainResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockCmdServerRecorder) GetDomain(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetDomain", arg0, arg1)
}

func (_m *MockCmdServer) GetDomainStats(_param0 context.Context, _param1 *EmptyRequest) (*DomainStatsResponse, error) {
	ret := _m.ctrl.Call(_m, "GetDomainStats", _param0, _param1)
	ret0, _ := ret[0].(*DomainStatsResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockCmdServerRecorder) GetDomainStats(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetDomainStats", arg0, arg1)
}

func (_m *MockCmdServer) GetGuestInfo(_param0 context.Context, _param1 *EmptyRequest) (*GuestInfoResponse, error) {
	ret := _m.ctrl.Call(_m, "GetGuestInfo", _param0, _param1)
	ret0, _ := ret[0].(*GuestInfoResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockCmdServerRecorder) GetGuestInfo(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetGuestInfo", arg0, arg1)
}

func (_m *MockCmdServer) GetUsers(_param0 context.Context, _param1 *EmptyRequest) (*GuestUserListResponse, error) {
	ret := _m.ctrl.Call(_m, "GetUsers", _param0, _param1)
	ret0, _ := ret[0].(*GuestUserListResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockCmdServerRecorder) GetUsers(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetUsers", arg0, arg1)
}

func (_m *MockCmdServer) GetFilesystems(_param0 context.Context, _param1 *EmptyRequest) (*GuestFilesystemsResponse, error) {
	ret := _m.ctrl.Call(_m, "GetFilesystems", _param0, _param1)
	ret0, _ := ret[0].(*GuestFilesystemsResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockCmdServerRecorder) GetFilesystems(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetFilesystems", arg0, arg1)
}

func (_m *MockCmdServer) Ping(_param0 context.Context, _param1 *EmptyRequest) (*Response, error) {
	ret := _m.ctrl.Call(_m, "Ping", _param0, _param1)
	ret0, _ := ret[0].(*Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockCmdServerRecorder) Ping(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Ping", arg0, arg1)
}

func (_m *MockCmdServer) Exec(_param0 context.Context, _param1 *ExecRequest) (*ExecResponse, error) {
	ret := _m.ctrl.Call(_m, "Exec", _param0, _param1)
	ret0, _ := ret[0].(*ExecResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockCmdServerRecorder) Exec(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Exec", arg0, arg1)
}

func (_m *MockCmdServer) GuestPing(_param0 context.Context, _param1 *GuestPingRequest) (*GuestPingResponse, error) {
	ret := _m.ctrl.Call(_m, "GuestPing", _param0, _param1)
	ret0, _ := ret[0].(*GuestPingResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockCmdServerRecorder) GuestPing(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GuestPing", arg0, arg1)
}

func (_m *MockCmdServer) VirtualMachineMemoryDump(_param0 context.Context, _param1 *MemoryDumpRequest) (*Response, error) {
	ret := _m.ctrl.Call(_m, "VirtualMachineMemoryDump", _param0, _param1)
	ret0, _ := ret[0].(*Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockCmdServerRecorder) VirtualMachineMemoryDump(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "VirtualMachineMemoryDump", arg0, arg1)
}

func (_m *MockCmdServer) GetVmmVersion(_param0 context.Context, _param1 *EmptyRequest) (*VmmVersionResponse, error) {
	ret := _m.ctrl.Call(_m, "GetVmmVersion", _param0, _param1)
	ret0, _ := ret[0].(*VmmVersionResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockCmdServerRecorder) GetVmmVersion(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetVmmVersion", arg0, arg1)
}
