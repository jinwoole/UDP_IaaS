package system

import (
	"udp_iaas/data/system"
)

func (s *Service) GetInfo() ([]system.SystemInfo, error) {
	info, err := s.data.GetSystemInfo()
	if err != nil {
		return nil, err
	}
	return []system.SystemInfo{*info}, nil
}
