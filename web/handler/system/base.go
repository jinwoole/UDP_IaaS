package system

import (
	systemService "udp_iaas/service/system"
)

type Handler struct {
	service *systemService.Service
}

func NewHandler(service *systemService.Service) *Handler {
	return &Handler{
		service: service,
	}
}
