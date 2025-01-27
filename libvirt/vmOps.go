package libvirt

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"
	"udp_iaas/types"

	libvirtgo "libvirt.org/go/libvirt"
)

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
	if req.ISOName != "" {
		isoPath = getISOPath(req.ISOName)
		if _, err := os.Stat(isoPath); err != nil {
			os.Remove(diskPath) // cleanup
			return fmt.Errorf("ISO file not found: %w", err)
		}
	}

	// VM XML 생성 및 도메인 정의
	xmlConfig := generateVMXML(req.Name, req.Memory, req.Cores, diskPath, isoPath)
	domain, err := conn.DomainDefineXML(xmlConfig)
	if err != nil {
		os.Remove(diskPath)
		return fmt.Errorf("failed to define domain: %w", err)
	}

	// VM 시작
	if err := domain.Create(); err != nil {
		domain.Undefine()
		os.Remove(diskPath)
		return fmt.Errorf("failed to start domain: %w", err)
	}

	return nil
}

// StartVM starts a stopped VM
func StartVM(conn *libvirtgo.Connect, name string) error {
    domain, err := conn.LookupDomainByName(name)
    if err != nil {
        return fmt.Errorf("failed to find domain: %w", err)
    }
    defer domain.Free()

    state, _, err := domain.GetState()
    if err != nil {
        return fmt.Errorf("failed to get domain state: %w", err)
    }

    if state == libvirtgo.DOMAIN_RUNNING {
        return nil
    }

    if err := domain.Create(); err != nil {
        return fmt.Errorf("failed to start domain: %w", err)
    }

    return nil
}

func GetVNCPort(conn *libvirtgo.Connect, name string) (int, error) {
    domain, err := conn.LookupDomainByName(name)
    if err != nil {
        return 0, fmt.Errorf("failed to find domain: %w", err)
    }
    defer domain.Free()

    // VNC 포트는 5900부터 순차적으로 할당됨
    return 5900, nil
}

// StopVM stops a running VM
func StopVM(conn *libvirtgo.Connect, name string) error {
	domain, err := conn.LookupDomainByName(name)
	if err != nil {
		return fmt.Errorf("failed to find domain: %w", err)
	}
	defer domain.Free()

	state, _, err := domain.GetState()
	if err != nil {
		return fmt.Errorf("failed to get domain state: %w", err)
	}

	if state == libvirtgo.DOMAIN_SHUTOFF {
		return nil
	}

	// 정상 종료 시도
	if err := domain.Shutdown(); err != nil {
		return fmt.Errorf("failed to shutdown domain: %w", err)
	}

	// 최대 30초 대기
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

	// 강제 종료
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

	// VM 중지
	if err := StopVM(conn, name); err != nil {
		return fmt.Errorf("failed to stop domain: %w", err)
	}

	// 도메인 정의 제거
	if err := domain.Undefine(); err != nil {
		return fmt.Errorf("failed to undefine domain: %w", err)
	}

	// 디스크 파일 정리
	diskPath := fmt.Sprintf("/var/lib/vms/disks/%s.img", name)
	if err := os.Remove(diskPath); err != nil {
		log.Printf("Warning: failed to remove disk file %s: %v", diskPath, err)
	}

	return nil
}