package vm

import (
	vmData "udp_iaas/data/vm"
)

func (s *Service) GetAllVMs() ([]vmData.VM, error) {
	domains, err := s.data.GetAllDomains()
	if err != nil {
		return nil, err
	}

	var vms []vmData.VM
	for _, domain := range domains {
		name, err := domain.GetName()
		if err != nil {
			domain.Free()
			continue
		}

		state, _, err := domain.GetState()
		if err != nil {
			domain.Free()
			continue
		}

		vms = append(vms, vmData.VM{
			Name:  name,
			State: uint(state),
		})
		domain.Free()
	}

	return vms, nil
}
