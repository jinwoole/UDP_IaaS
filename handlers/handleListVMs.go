package handlers

import (
	"encoding/json"
	"net/http"

	"udp_iaas/types"

	"libvirt.org/go/libvirt"
)

// VM 목록 조회 요청을 처리하는 핸들러
func (app *App) handleListVMs(w http.ResponseWriter, r *http.Request) {
	domains, err := app.Libvirt.ListAllDomains(libvirt.CONNECT_LIST_DOMAINS_ACTIVE |
		libvirt.CONNECT_LIST_DOMAINS_INACTIVE)
	if err != nil {
		http.Error(w, "Failed to list domains", http.StatusInternalServerError)
		return
	}

	vms := []types.VM{}
	for _, domain := range domains {
		name, err := domain.GetName()
		if err != nil {
			continue
		}

		info, err := domain.GetInfo()
		if err != nil {
			continue
		}

		// VM 상태를 문자열로 변환
		state := "unknown"
		switch info.State {
		case libvirt.DOMAIN_RUNNING:
			state = "running"
		case libvirt.DOMAIN_PAUSED:
			state = "paused"
		case libvirt.DOMAIN_SHUTOFF:
			state = "stopped"
		case libvirt.DOMAIN_CRASHED:
			state = "crashed"
		case libvirt.DOMAIN_PMSUSPENDED:
			state = "suspended"
		}

		vms = append(vms, types.VM{
			Name:   name,
			Cores:  int(info.NrVirtCpu),
			Memory: int(info.Memory),
			State:  state,
		})

		domain.Free()
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(vms)
}
