package system

import (
	"udp_iaas/data/system"
)

func (s *Service) GetISO() ([]system.ISOInfo, error) {
	info, err := s.data.GetISOFiles()
	if err != nil {
		return nil, err
	}
	return info, nil
}
