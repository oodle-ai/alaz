// Code generated by bpf2go; DO NOT EDIT.
//go:build 386 || amd64 || amd64p32 || arm || arm64 || loong64 || mips64le || mips64p32le || mipsle || ppc64le || riscv64

package c

import (
	"bytes"
	_ "embed"
	"fmt"
	"io"

	"github.com/cilium/ebpf"
)

type bpfCallEvent struct {
	Pid  uint32
	Tid  uint32
	Tx   uint64
	Type uint8
	_    [3]byte
	Seq  uint32
}

type bpfGoReadKey struct {
	Pid  uint32
	_    [4]byte
	Goid uint64
}

type bpfGoReqKey struct {
	Pid uint32
	_   [4]byte
	Fd  uint64
}

type bpfL7Event struct {
	Fd                  uint64
	WriteTimeNs         uint64
	Pid                 uint32
	Status              uint32
	Duration            uint64
	Protocol            uint8
	Method              uint8
	Padding             uint16
	Payload             [1024]uint8
	PayloadSize         uint32
	PayloadReadComplete uint8
	Failed              uint8
	IsTls               uint8
	_                   [1]byte
	Seq                 uint32
	Tid                 uint32
	KafkaApiVersion     int16
	_                   [2]byte
	Saddr               uint32
	Sport               uint16
	_                   [2]byte
	Daddr               uint32
	Dport               uint16
	_                   [2]byte
}

type bpfL7Request struct {
	WriteTimeNs         uint64
	Protocol            uint8
	Method              uint8
	Payload             [1024]uint8
	_                   [2]byte
	PayloadSize         uint32
	PayloadReadComplete uint8
	RequestType         uint8
	_                   [2]byte
	Seq                 uint32
	Tid                 uint32
	CorrelationId       int32
	ApiKey              int16
	ApiVersion          int16
	Saddr               uint32
	Sport               uint16
	_                   [2]byte
	Daddr               uint32
	Dport               uint16
	_                   [6]byte
}

type bpfLogMessage struct {
	Level    uint32
	LogMsg   [100]uint8
	FuncName [100]uint8
	Pid      uint32
	Arg1     uint64
	Arg2     uint64
	Arg3     uint64
}

type bpfP_event struct {
	Pid  uint32
	Type uint8
	_    [3]byte
}

type bpfSkInfo struct {
	Fd  uint64
	Pid uint32
	_   [4]byte
}

type bpfSocketKey struct {
	Fd    uint64
	Pid   uint32
	IsTls uint8
	_     [3]byte
}

// loadBpf returns the embedded CollectionSpec for bpf.
func loadBpf() (*ebpf.CollectionSpec, error) {
	reader := bytes.NewReader(_BpfBytes)
	spec, err := ebpf.LoadCollectionSpecFromReader(reader)
	if err != nil {
		return nil, fmt.Errorf("can't load bpf: %w", err)
	}

	return spec, err
}

// loadBpfObjects loads bpf and converts it into a struct.
//
// The following types are suitable as obj argument:
//
//	*bpfObjects
//	*bpfPrograms
//	*bpfMaps
//
// See ebpf.CollectionSpec.LoadAndAssign documentation for details.
func loadBpfObjects(obj interface{}, opts *ebpf.CollectionOptions) error {
	spec, err := loadBpf()
	if err != nil {
		return err
	}

	return spec.LoadAndAssign(obj, opts)
}

// bpfSpecs contains maps and programs before they are loaded into the kernel.
//
// It can be passed ebpf.CollectionSpec.Assign.
type bpfSpecs struct {
	bpfProgramSpecs
	bpfMapSpecs
}

// bpfSpecs contains programs before they are loaded into the kernel.
//
// It can be passed ebpf.CollectionSpec.Assign.
type bpfProgramSpecs struct {
	GoTlsConnReadEnter  *ebpf.ProgramSpec `ebpf:"go_tls_conn_read_enter"`
	GoTlsConnReadExit   *ebpf.ProgramSpec `ebpf:"go_tls_conn_read_exit"`
	GoTlsConnWriteEnter *ebpf.ProgramSpec `ebpf:"go_tls_conn_write_enter"`
	InetSockSetState    *ebpf.ProgramSpec `ebpf:"inet_sock_set_state"`
	SchedProcessExec    *ebpf.ProgramSpec `ebpf:"sched_process_exec"`
	SchedProcessExit    *ebpf.ProgramSpec `ebpf:"sched_process_exit"`
	SchedProcessFork    *ebpf.ProgramSpec `ebpf:"sched_process_fork"`
	SslReadEnterV102    *ebpf.ProgramSpec `ebpf:"ssl_read_enter_v1_0_2"`
	SslReadEnterV111    *ebpf.ProgramSpec `ebpf:"ssl_read_enter_v1_1_1"`
	SslReadEnterV3      *ebpf.ProgramSpec `ebpf:"ssl_read_enter_v3"`
	SslRetRead          *ebpf.ProgramSpec `ebpf:"ssl_ret_read"`
	SslWriteV102        *ebpf.ProgramSpec `ebpf:"ssl_write_v1_0_2"`
	SslWriteV111        *ebpf.ProgramSpec `ebpf:"ssl_write_v1_1_1"`
	SslWriteV3          *ebpf.ProgramSpec `ebpf:"ssl_write_v3"`
	SysEnterConnect     *ebpf.ProgramSpec `ebpf:"sys_enter_connect"`
	SysEnterRead        *ebpf.ProgramSpec `ebpf:"sys_enter_read"`
	SysEnterRecvfrom    *ebpf.ProgramSpec `ebpf:"sys_enter_recvfrom"`
	SysEnterSendto      *ebpf.ProgramSpec `ebpf:"sys_enter_sendto"`
	SysEnterWrite       *ebpf.ProgramSpec `ebpf:"sys_enter_write"`
	SysEnterWritev      *ebpf.ProgramSpec `ebpf:"sys_enter_writev"`
	SysExitConnect      *ebpf.ProgramSpec `ebpf:"sys_exit_connect"`
	SysExitRead         *ebpf.ProgramSpec `ebpf:"sys_exit_read"`
	SysExitRecvfrom     *ebpf.ProgramSpec `ebpf:"sys_exit_recvfrom"`
	SysExitSendto       *ebpf.ProgramSpec `ebpf:"sys_exit_sendto"`
	SysExitWrite        *ebpf.ProgramSpec `ebpf:"sys_exit_write"`
	SysExitWritev       *ebpf.ProgramSpec `ebpf:"sys_exit_writev"`
}

// bpfMapSpecs contains maps before they are loaded into the kernel.
//
// It can be passed ebpf.CollectionSpec.Assign.
type bpfMapSpecs struct {
	ActiveL7Requests   *ebpf.MapSpec `ebpf:"active_l7_requests"`
	ActiveReads        *ebpf.MapSpec `ebpf:"active_reads"`
	ActiveWrites       *ebpf.MapSpec `ebpf:"active_writes"`
	ContainerPids      *ebpf.MapSpec `ebpf:"container_pids"`
	FdByPidTgid        *ebpf.MapSpec `ebpf:"fd_by_pid_tgid"`
	GoActiveL7Requests *ebpf.MapSpec `ebpf:"go_active_l7_requests"`
	GoActiveReads      *ebpf.MapSpec `ebpf:"go_active_reads"`
	GoL7RequestHeap    *ebpf.MapSpec `ebpf:"go_l7_request_heap"`
	IngressEgressCalls *ebpf.MapSpec `ebpf:"ingress_egress_calls"`
	IngressEgressHeap  *ebpf.MapSpec `ebpf:"ingress_egress_heap"`
	L7EventHeap        *ebpf.MapSpec `ebpf:"l7_event_heap"`
	L7Events           *ebpf.MapSpec `ebpf:"l7_events"`
	L7RequestHeap      *ebpf.MapSpec `ebpf:"l7_request_heap"`
	LogHeap            *ebpf.MapSpec `ebpf:"log_heap"`
	LogMap             *ebpf.MapSpec `ebpf:"log_map"`
	ProcEventHeap      *ebpf.MapSpec `ebpf:"proc_event_heap"`
	ProcEvents         *ebpf.MapSpec `ebpf:"proc_events"`
	SockMap            *ebpf.MapSpec `ebpf:"sock_map"`
	SockMapTemp        *ebpf.MapSpec `ebpf:"sock_map_temp"`
	TcpConnectEvents   *ebpf.MapSpec `ebpf:"tcp_connect_events"`
	TcpListenEvents    *ebpf.MapSpec `ebpf:"tcp_listen_events"`
}

// bpfObjects contains all objects after they have been loaded into the kernel.
//
// It can be passed to loadBpfObjects or ebpf.CollectionSpec.LoadAndAssign.
type bpfObjects struct {
	bpfPrograms
	bpfMaps
}

func (o *bpfObjects) Close() error {
	return _BpfClose(
		&o.bpfPrograms,
		&o.bpfMaps,
	)
}

// bpfMaps contains all maps after they have been loaded into the kernel.
//
// It can be passed to loadBpfObjects or ebpf.CollectionSpec.LoadAndAssign.
type bpfMaps struct {
	ActiveL7Requests   *ebpf.Map `ebpf:"active_l7_requests"`
	ActiveReads        *ebpf.Map `ebpf:"active_reads"`
	ActiveWrites       *ebpf.Map `ebpf:"active_writes"`
	ContainerPids      *ebpf.Map `ebpf:"container_pids"`
	FdByPidTgid        *ebpf.Map `ebpf:"fd_by_pid_tgid"`
	GoActiveL7Requests *ebpf.Map `ebpf:"go_active_l7_requests"`
	GoActiveReads      *ebpf.Map `ebpf:"go_active_reads"`
	GoL7RequestHeap    *ebpf.Map `ebpf:"go_l7_request_heap"`
	IngressEgressCalls *ebpf.Map `ebpf:"ingress_egress_calls"`
	IngressEgressHeap  *ebpf.Map `ebpf:"ingress_egress_heap"`
	L7EventHeap        *ebpf.Map `ebpf:"l7_event_heap"`
	L7Events           *ebpf.Map `ebpf:"l7_events"`
	L7RequestHeap      *ebpf.Map `ebpf:"l7_request_heap"`
	LogHeap            *ebpf.Map `ebpf:"log_heap"`
	LogMap             *ebpf.Map `ebpf:"log_map"`
	ProcEventHeap      *ebpf.Map `ebpf:"proc_event_heap"`
	ProcEvents         *ebpf.Map `ebpf:"proc_events"`
	SockMap            *ebpf.Map `ebpf:"sock_map"`
	SockMapTemp        *ebpf.Map `ebpf:"sock_map_temp"`
	TcpConnectEvents   *ebpf.Map `ebpf:"tcp_connect_events"`
	TcpListenEvents    *ebpf.Map `ebpf:"tcp_listen_events"`
}

func (m *bpfMaps) Close() error {
	return _BpfClose(
		m.ActiveL7Requests,
		m.ActiveReads,
		m.ActiveWrites,
		m.ContainerPids,
		m.FdByPidTgid,
		m.GoActiveL7Requests,
		m.GoActiveReads,
		m.GoL7RequestHeap,
		m.IngressEgressCalls,
		m.IngressEgressHeap,
		m.L7EventHeap,
		m.L7Events,
		m.L7RequestHeap,
		m.LogHeap,
		m.LogMap,
		m.ProcEventHeap,
		m.ProcEvents,
		m.SockMap,
		m.SockMapTemp,
		m.TcpConnectEvents,
		m.TcpListenEvents,
	)
}

// bpfPrograms contains all programs after they have been loaded into the kernel.
//
// It can be passed to loadBpfObjects or ebpf.CollectionSpec.LoadAndAssign.
type bpfPrograms struct {
	GoTlsConnReadEnter  *ebpf.Program `ebpf:"go_tls_conn_read_enter"`
	GoTlsConnReadExit   *ebpf.Program `ebpf:"go_tls_conn_read_exit"`
	GoTlsConnWriteEnter *ebpf.Program `ebpf:"go_tls_conn_write_enter"`
	InetSockSetState    *ebpf.Program `ebpf:"inet_sock_set_state"`
	SchedProcessExec    *ebpf.Program `ebpf:"sched_process_exec"`
	SchedProcessExit    *ebpf.Program `ebpf:"sched_process_exit"`
	SchedProcessFork    *ebpf.Program `ebpf:"sched_process_fork"`
	SslReadEnterV102    *ebpf.Program `ebpf:"ssl_read_enter_v1_0_2"`
	SslReadEnterV111    *ebpf.Program `ebpf:"ssl_read_enter_v1_1_1"`
	SslReadEnterV3      *ebpf.Program `ebpf:"ssl_read_enter_v3"`
	SslRetRead          *ebpf.Program `ebpf:"ssl_ret_read"`
	SslWriteV102        *ebpf.Program `ebpf:"ssl_write_v1_0_2"`
	SslWriteV111        *ebpf.Program `ebpf:"ssl_write_v1_1_1"`
	SslWriteV3          *ebpf.Program `ebpf:"ssl_write_v3"`
	SysEnterConnect     *ebpf.Program `ebpf:"sys_enter_connect"`
	SysEnterRead        *ebpf.Program `ebpf:"sys_enter_read"`
	SysEnterRecvfrom    *ebpf.Program `ebpf:"sys_enter_recvfrom"`
	SysEnterSendto      *ebpf.Program `ebpf:"sys_enter_sendto"`
	SysEnterWrite       *ebpf.Program `ebpf:"sys_enter_write"`
	SysEnterWritev      *ebpf.Program `ebpf:"sys_enter_writev"`
	SysExitConnect      *ebpf.Program `ebpf:"sys_exit_connect"`
	SysExitRead         *ebpf.Program `ebpf:"sys_exit_read"`
	SysExitRecvfrom     *ebpf.Program `ebpf:"sys_exit_recvfrom"`
	SysExitSendto       *ebpf.Program `ebpf:"sys_exit_sendto"`
	SysExitWrite        *ebpf.Program `ebpf:"sys_exit_write"`
	SysExitWritev       *ebpf.Program `ebpf:"sys_exit_writev"`
}

func (p *bpfPrograms) Close() error {
	return _BpfClose(
		p.GoTlsConnReadEnter,
		p.GoTlsConnReadExit,
		p.GoTlsConnWriteEnter,
		p.InetSockSetState,
		p.SchedProcessExec,
		p.SchedProcessExit,
		p.SchedProcessFork,
		p.SslReadEnterV102,
		p.SslReadEnterV111,
		p.SslReadEnterV3,
		p.SslRetRead,
		p.SslWriteV102,
		p.SslWriteV111,
		p.SslWriteV3,
		p.SysEnterConnect,
		p.SysEnterRead,
		p.SysEnterRecvfrom,
		p.SysEnterSendto,
		p.SysEnterWrite,
		p.SysEnterWritev,
		p.SysExitConnect,
		p.SysExitRead,
		p.SysExitRecvfrom,
		p.SysExitSendto,
		p.SysExitWrite,
		p.SysExitWritev,
	)
}

func _BpfClose(closers ...io.Closer) error {
	for _, closer := range closers {
		if err := closer.Close(); err != nil {
			return err
		}
	}
	return nil
}

// Do not access this directly.
//
//go:embed bpf_bpfel.o
var _BpfBytes []byte
