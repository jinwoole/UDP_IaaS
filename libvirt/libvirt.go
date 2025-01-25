package libvirt

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"udp_iaas/types"

	libvirtgo "libvirt.org/go/libvirt"
)

const isoStoragePath = "/var/lib/vms/isos"

const vmXMLTemplate = `
<domain type='kvm'>
    <name>%s</name>
    <memory unit='MiB'>%d</memory>
    <vcpu>%d</vcpu>
    <os>
        <type arch='x86_64'>hvm</type>
        <boot dev='cdrom'/>    <!-- ISO에서 먼저 부팅 시도 -->
        <boot dev='hd'/>       <!-- 하드디스크는 두 번째 부팅 옵션 -->
    </os>
    <devices>
        <!-- 시스템 디스크 -->
        <disk type='file' device='disk'>
            <driver name='qemu' type='raw'/>
            <source file='/var/lib/vms/disks/%s.img'/>
            <target dev='vda' bus='virtio'/>
        </disk>
        <!-- ISO 마운트를 위한 CD-ROM 추가 -->
        <disk type='file' device='cdrom'>
            <driver name='qemu' type='raw'/>
            <source file='%s'/>
            <target dev='hdc' bus='ide'/>
            <readonly/>
        </disk>
        <interface type='network'>
            <source network='default'/>
            <model type='virtio'/>
        </interface>
        <console type='pty'/>
        <!-- VGA 디스플레이 추가 -->
        <video>
            <model type='vga'/>
        </video>
    </devices>
</domain>
`

// CreateVM creates a new VM
func CreateVM(conn *libvirtgo.Connect, req types.CreateVMRequest) error {
	// 시스템 디스크 생성
	diskPath := fmt.Sprintf("/var/lib/vms/disks/%s.img", req.Name)
	cmd := exec.Command("qemu-img", "create", "-f", "raw", diskPath, "10G")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create disk image: %w", err)
	}

	// ISO 파일 경로 처리
	var isoPath string
	if req.ISOName != "" { // This should now match the struct field
		isoPath = filepath.Join(isoStoragePath, req.ISOName)
		if _, err := os.Stat(isoPath); err != nil {
			os.Remove(diskPath) // cleanup
			return fmt.Errorf("ISO file not found: %w", err)
		}
	}

	// CD-ROM 디바이스 XML 조건부 생성
	cdromXML := ""
	if isoPath != "" {
		cdromXML = fmt.Sprintf(`
        <disk type='file' device='cdrom'>
            <driver name='qemu' type='raw'/>
            <source file='%s'/>
            <target dev='hdc' bus='ide'/>
            <readonly/>
        </disk>`, isoPath)
	}

	// VM XML 템플릿 수정
	xmlConfig := fmt.Sprintf(`
<domain type='kvm'>
    <name>%s</name>
    <memory unit='MiB'>%d</memory>
    <vcpu>%d</vcpu>
    <os>
        <type arch='x86_64'>hvm</type>
        <boot dev='cdrom'/>
        <boot dev='hd'/>
    </os>
    <devices>
        <disk type='file' device='disk'>
            <driver name='qemu' type='raw'/>
            <source file='%s'/>
            <target dev='vda' bus='virtio'/>
        </disk>
        %s
        <interface type='network'>
            <source network='default'/>
            <model type='virtio'/>
        </interface>
        <console type='pty'/>
        <video>
            <model type='vga'/>
        </video>
    </devices>
</domain>`,
		req.Name,   // VM 이름
		req.Memory, // 메모리 크기
		req.Cores,  // CPU 코어 수
		diskPath,   // 디스크 이미지 경로
		cdromXML,   // CD-ROM 설정 (ISO가 있는 경우에만)
	)

	// 도메인 정의 및 시작
	domain, err := conn.DomainDefineXML(xmlConfig)
	if err != nil {
		os.Remove(diskPath)
		return fmt.Errorf("failed to define domain: %w", err)
	}

	if err := domain.Create(); err != nil {
		domain.Undefine()
		os.Remove(diskPath)
		return fmt.Errorf("failed to start domain: %w", err)
	}

	return nil
}

func StopVM(conn *libvirtgo.Connect, name string) error {
	// 도메인을 찾습니다
	domain, err := conn.LookupDomainByName(name)
	if err != nil {
		return fmt.Errorf("failed to find domain: %w", err)
	}
	defer domain.Free()

	// 도메인의 상태를 확인합니다
	state, _, err := domain.GetState()
	if err != nil {
		return fmt.Errorf("failed to get domain state: %w", err)
	}

	// 이미 멈춰있다면 아무것도 하지 않습니다
	if state == libvirtgo.DOMAIN_SHUTOFF {
		return nil
	}

	// 정상적인 종료를 시도합니다
	if err := domain.Shutdown(); err != nil {
		return fmt.Errorf("failed to shutdown domain: %w", err)
	}

	// 최대 30초 동안 종료되기를 기다립니다
	for i := 0; i < 30; i++ {
		state, _, err := domain.GetState()
		if err != nil {
			return fmt.Errorf("failed to get domain state: %w", err)
		}
		if state == libvirtgo.DOMAIN_SHUTOFF {
			return nil
		}
		time.Sleep(time.Second)
	}

	// 여전히 종료되지 않았다면 강제 종료를 시도합니다
	if err := domain.Destroy(); err != nil {
		return fmt.Errorf("failed to force stop domain: %w", err)
	}

	return nil
}

// DeleteVM deletes an existing VM
func DeleteVM(conn *libvirtgo.Connect, name string) error {
	domain, err := conn.LookupDomainByName(name)
	if err != nil {
		return fmt.Errorf("failed to find domain: %w", err)
	}
	defer domain.Free()

	// 도메인이 실행 중이라면 중지
	if err := StopVM(conn, name); err != nil {
		return fmt.Errorf("failed to stop domain: %w", err)
	}

	// 도메인 제거
	if err := domain.Undefine(); err != nil {
		return fmt.Errorf("failed to undefine domain: %w", err)
	}

	// 디스크 파일 정리
	diskPath := fmt.Sprintf("/var/lib/vms/disks/%s.img", name)
	if err := os.Remove(diskPath); err != nil {
		// 디스크 삭제 실패는 로그만 하고 계속 진행
		log.Printf("Warning: failed to remove disk file %s: %v", diskPath, err)
	}

	return nil
}

// ListISOs returns the list of ISO files
func ListISOs() ([]types.ISO, error) {
	// 디렉토리가 없다면 생성
	if err := os.MkdirAll(isoStoragePath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create iso directory: %w", err)
	}

	entries, err := os.ReadDir(isoStoragePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read iso directory: %w", err)
	}

	var isos []types.ISO // Make sure this is types.ISO, not just ISO
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(strings.ToLower(entry.Name()), ".iso") {
			info, err := entry.Info()
			if err != nil {
				continue
			}

			isos = append(isos, types.ISO{ // Make sure this is types.ISO, not just ISO
				Name: entry.Name(),
				Size: info.Size(),
			})
		}
	}

	return isos, nil
}

// SaveISO saves an uploaded ISO to disk
func SaveISO(fileName string, file io.Reader) error {
	if !strings.HasSuffix(strings.ToLower(fileName), ".iso") {
		return fmt.Errorf("file must have .iso extension")
	}

	// 디렉토리가 없다면 생성
	if err := os.MkdirAll(isoStoragePath, 0755); err != nil {
		return fmt.Errorf("failed to create iso directory: %w", err)
	}

	// 파일 생성
	dst, err := os.Create(filepath.Join(isoStoragePath, fileName))
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer dst.Close()

	// 파일 복사
	if _, err := io.Copy(dst, file); err != nil {
		return fmt.Errorf("failed to save file: %w", err)
	}

	return nil
}
