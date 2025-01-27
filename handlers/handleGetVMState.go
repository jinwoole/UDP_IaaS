package handlers

import (
	"encoding/json"
	"net/http"

	"libvirt.org/go/libvirt"
)

// handlers/handleVMs.go
func (app *App) handleGetVMState(w http.ResponseWriter, r *http.Request, vmName string) {
	domain, err := app.Libvirt.LookupDomainByName(vmName)
	if err != nil {
		http.Error(w, "VM not found", http.StatusNotFound)
		return
	}
	defer domain.Free()
 
	info, err := domain.GetInfo()
	if err != nil {
		http.Error(w, "Failed to get VM info", http.StatusInternalServerError)
		return
	}
 
	state := "unknown"
	switch info.State {
	case libvirt.DOMAIN_RUNNING:
		state = "running"
	case libvirt.DOMAIN_PAUSED:
		state = "paused"
	case libvirt.DOMAIN_SHUTOFF:
		state = "stopped"
	}
 
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"state": state,
	})
 }