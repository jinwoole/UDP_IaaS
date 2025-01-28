package libvirt

import (
	"fmt"
	"strings"
	"time"

	libvirtgo "libvirt.org/go/libvirt"
)

const (
	isoStoragePath = "/var/lib/vms/isos"
)

func getISOPath(isoName string) string {
	return isoStoragePath + "/" + isoName
}

func DetachISO(conn *libvirtgo.Connect, name string) error {
    domain, err := conn.LookupDomainByName(name)
    if err != nil {
        return fmt.Errorf("failed to find domain: %w", err)
    }
    defer domain.Free()

    // Get current XML
    xmlDesc, err := domain.GetXMLDesc(0)
    if err != nil {
        return fmt.Errorf("failed to get domain XML: %w", err)
    }

    // 매우 단순하게 cdrom 디바이스 부분을 제거
    parts := strings.Split(xmlDesc, "\n")
    var newParts []string
    skipMode := false

    for _, line := range parts {
        if strings.Contains(line, `device='cdrom'`) {
            skipMode = true
            continue
        }
        
        // cdrom 디바이스 블록의 끝을 감지
        if skipMode && strings.Contains(line, "</disk>") {
            skipMode = false
            continue
        }
        
        if !skipMode {
            newParts = append(newParts, line)
        }
    }

    newXML := strings.Join(newParts, "\n")

    // Get current state
    state, _, err := domain.GetState()
    if err != nil {
        return fmt.Errorf("failed to get domain state: %w", err)
    }

    wasRunning := state == libvirtgo.DOMAIN_RUNNING
    
    // Stop if running
    if wasRunning {
        if err := domain.Shutdown(); err != nil {
            // Try force stop if normal shutdown fails
            if err := domain.Destroy(); err != nil {
                return fmt.Errorf("failed to stop domain: %w", err)
            }
        }
        // Wait for stop
        for i := 0; i < 30; i++ {
            state, _, err := domain.GetState()
            if err != nil {
                return fmt.Errorf("failed to get domain state: %w", err)
            }
            if state == libvirtgo.DOMAIN_SHUTOFF {
                break
            }
            time.Sleep(time.Second)
        }
    }

    // Define new configuration
    newDomain, err := conn.DomainDefineXML(newXML)
    if err != nil {
        return fmt.Errorf("failed to define new domain configuration: %w", err)
    }
    defer newDomain.Free()

    // Restart if it was running
    if wasRunning {
        if err := newDomain.Create(); err != nil {
            return fmt.Errorf("failed to start domain: %w", err)
        }
    }

    return nil
}