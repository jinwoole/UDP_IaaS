package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	localLibvirt "udp_iaas/libvirt"
)

func (app *App) handleDetachISO(w http.ResponseWriter, r *http.Request, vmName string) {
    // VM에서 ISO 제거를 비동기로 실행
    go func() {
        if err := localLibvirt.DetachISO(app.Libvirt, vmName); err != nil {
            log.Printf("Failed to detach ISO from VM %s: %v", vmName, err)
            return
        }
        log.Printf("Successfully detached ISO from VM: %s", vmName)
    }()

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusAccepted)
    json.NewEncoder(w).Encode(map[string]string{
        "message": fmt.Sprintf("ISO detachment initiated for VM: %s", vmName),
    })
}