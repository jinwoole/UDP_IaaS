package vm

import (
	vmservice "udp_iaas/service/vm"
)

type Handler struct {
	service *vmservice.Service
}

func NewHandler(service *vmservice.Service) *Handler {
	return &Handler{
		service: service,
	}
}
