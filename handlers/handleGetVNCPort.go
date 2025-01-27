package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	localLibvirt "udp_iaas/libvirt"
)

func (app *App) handleGetVNCPort(w http.ResponseWriter, r *http.Request, vmName string) {
    vncPort, err := localLibvirt.GetVNCPort(app.Libvirt, vmName)
    if err != nil {
        http.Error(w, fmt.Sprintf("Failed to get VNC port: %v", err), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]int{
        "port": vncPort + 1,
    })
}