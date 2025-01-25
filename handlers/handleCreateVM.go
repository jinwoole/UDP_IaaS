package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	localLibvirt "udp_iaas/libvirt"
	"udp_iaas/types"
)

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
