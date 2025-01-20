package vm

import (
	vmdata "udp_iaas/data/vm"
)

type Service struct {
	data *vmdata.Data
}

func NewService(data *vmdata.Data) *Service {
	return &Service{
		data: data,
	}
}
