package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	localLibvirt "udp_iaas/libvirt"
)

func (app *App) handleStartVM(w http.ResponseWriter, r *http.Request, vmName string) {
    go func() {
        if err := localLibvirt.StartVM(app.Libvirt, vmName); err != nil {
            log.Printf("Failed to start VM %s: %v", vmName, err)
            return
        }
        log.Printf("Successfully started VM: %s", vmName)
    }()

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusAccepted)
    json.NewEncoder(w).Encode(map[string]string{
        "message": fmt.Sprintf("VM start initiated for: %s", vmName),
    })
}