package system

import "libvirt.org/go/libvirt"

// SystemInfo는 시스템의 하드웨어 및 가상화 관련 정보를 포함하는 구조체입니다.
//
// Fields:
//   - CPUCount: 시스템의 CPU 코어 수
//   - MemorySize: 시스템의 총 메모리 크기 (KiB 단위)
//   - Hostname: 호스트 시스템의 이름
//   - NodeInfo: libvirt 노드 정보를 포함하는 구조체 포인터
//   - CPUModel: CPU 모델명
//   - CPUArch: CPU 아키텍처 (예: x86_64, arm64 등)
//   - LibVersion: libvirt 라이브러리 버전
//   - HostOS: 호스트 운영체제 정보
//   - VirtType: 가상화 타입 (예: KVM, QEMU 등)
//   - FreeMemory: 사용 가능한 메모리 크기
//   - MaxVCPUs: 지원되는 최대 가상 CPU 수
//   - IsConnected: libvirt 연결 상태 (true: 연결됨, false: 연결 안됨)

type SystemInfo struct {
	CPUCount    uint
	MemorySize  uint64 // in KiB
	Hostname    string
	NodeInfo    *libvirt.NodeInfo
	CPUModel    string
	CPUArch     string
	LibVersion  uint64
	HostOS      string
	VirtType    string
	FreeMemory  uint64
	MaxVCPUs    uint
	IsConnected bool
}

type ISOInfo struct {
	Name string
	Path string
}
