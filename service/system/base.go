package system

import (
	"udp_iaas/data/system"
)

type Service struct {
	data *system.Data
}

func NewService(data *system.Data) *Service {
	return &Service{
		data: data,
	}
}
