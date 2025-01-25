package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	localLibvirt "udp_iaas/libvirt"
	"udp_iaas/types"

	"libvirt.org/go/libvirt"
)

// App holds references to external services or connections
type App struct {
	Libvirt *libvirt.Connect
}

// VM 생성 요청을 처리하는 핸들러
func (app *App) handleCreateVM(w http.ResponseWriter, r *http.Request) {
	var req types.CreateVMRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// 기본값 설정
	if req.Cores <= 0 {
		req.Cores = 1
	}
	if req.Memory <= 0 {
		req.Memory = 1024 // 1GB
	}

	// VM 생성을 비동기로 실행
	go func() {
		if err := localLibvirt.CreateVM(app.Libvirt, req); err != nil {
			log.Printf("Failed to create VM %s: %v", req.Name, err)
			return
		}
		log.Printf("Successfully created VM: %s", req.Name)
	}()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]string{
		"message": fmt.Sprintf("VM creation started for: %s", req.Name),
	})
}

func (app *App) handleStopVM(w http.ResponseWriter, r *http.Request, vmName string) {
	// VM 중지를 비동기로 실행
	go func() {
		if err := localLibvirt.StopVM(app.Libvirt, vmName); err != nil {
			log.Printf("Failed to stop VM %s: %v", vmName, err)
			return
		}
		log.Printf("Successfully stopped VM: %s", vmName)
	}()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]string{
		"message": fmt.Sprintf("VM stop initiated for: %s", vmName),
	})
}

// VM 삭제 요청을 처리하는 핸들러
func (app *App) handleDeleteVM(w http.ResponseWriter, r *http.Request, vmName string) {
	// VM 삭제를 비동기로 실행
	go func() {
		if err := localLibvirt.DeleteVM(app.Libvirt, vmName); err != nil {
			log.Printf("Failed to delete VM %s: %v", vmName, err)
			return
		}
		log.Printf("Successfully deleted VM: %s", vmName)
	}()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]string{
		"message": fmt.Sprintf("VM deletion initiated for: %s", vmName),
	})
}

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

// 헬스체크 엔드포인트
func (app *App) HandleHealthCheck(w http.ResponseWriter, r *http.Request) {
	// libvirt 연결 상태 확인
	_, err := app.Libvirt.GetVersion()
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(map[string]string{
			"status": "unhealthy",
			"error":  "libvirt connection failed",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "healthy",
	})
}

func (app *App) handleListISOs(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	isos, err := localLibvirt.ListISOs()
	if err != nil {
		http.Error(w, "Failed to list ISOs", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(isos)
}

func (app *App) handleUploadISO(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	file, header, err := r.FormFile("iso")
	if err != nil {
		http.Error(w, "Failed to get file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	if err := localLibvirt.SaveISO(header.Filename, file); err != nil {
		http.Error(w, fmt.Sprintf("Failed to save ISO: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": fmt.Sprintf("Successfully uploaded: %s", header.Filename),
	})
}

func (app *App) HandleVMs(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/api/vms" {
		switch r.Method {
		case http.MethodGet:
			app.handleListVMs(w, r)
		case http.MethodPost:
			app.handleCreateVM(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
		return
	}

	// Parse VM name from URL
	pathParts := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/vms/"), "/")
	if len(pathParts) == 0 || pathParts[0] == "" {
		http.Error(w, "VM name not specified", http.StatusBadRequest)
		return
	}
	vmName := pathParts[0]

	// Check if it's "stop"
	if len(pathParts) == 2 && pathParts[1] == "stop" && r.Method == http.MethodPost {
		app.handleStopVM(w, r, vmName)
		return
	}

	// Check if it's a DELETE for that VM
	if len(pathParts) == 1 && r.Method == http.MethodDelete {
		app.handleDeleteVM(w, r, vmName)
		return
	}

	http.Error(w, "Invalid path", http.StatusNotFound)
}

func (app *App) HandleISOs(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		app.handleListISOs(w, r)
	case http.MethodPost:
		app.handleUploadISO(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
