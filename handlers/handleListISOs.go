package handlers

import (
	"encoding/json"
	"net/http"

	localLibvirt "udp_iaas/libvirt"
)

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
