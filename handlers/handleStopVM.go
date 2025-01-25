package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	localLibvirt "udp_iaas/libvirt"
)

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
