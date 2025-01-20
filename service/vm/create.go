package vm

import (
	"fmt"
	vmData "udp_iaas/data/vm"
)

func (s *Service) CreateVM(req vmData.CreateVMRequest) error {
	if err := s.validateVMRequest(req); err != nil {
		return err
	}

	xmlConfig := s.generateVMConfig(req)
	return s.data.CreateDomain(xmlConfig)
}

func (s *Service) validateVMRequest(req vmData.CreateVMRequest) error {
	if req.Name == "" {
		return fmt.Errorf("name is required")
	}
	if req.Memory < 512 {
		return fmt.Errorf("memory must be at least 512MB")
	}
	if req.VCPU < 1 {
		return fmt.Errorf("vcpu must be at least 1")
	}
	if req.ImagePath == "" {
		return fmt.Errorf("image path is required")
	}
	return nil
}

func (s *Service) generateVMConfig(req vmData.CreateVMRequest) string {
	xmlTemplate := `
    <domain type='kvm'>
        <name>%s</name>
        <memory unit='MiB'>%d</memory>
        <vcpu>%d</vcpu>
        <os>
            <type arch='x86_64'>hvm</type>
            <boot dev='hd'/>
        </os>
        <devices>
            <disk type='file' device='disk'>
                <driver name='qemu' type='qcow2'/>
                <source file='%s'/>
                <target dev='vda' bus='virtio'/>
            </disk>
            <interface type='network'>
                <source network='default'/>
                <model type='virtio'/>
            </interface>
        </devices>
    </domain>
    `

	return fmt.Sprintf(xmlTemplate,
		req.Name,
		req.Memory,
		req.VCPU,
		req.ImagePath,
	)
}
