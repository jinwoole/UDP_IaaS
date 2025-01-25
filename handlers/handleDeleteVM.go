package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	localLibvirt "udp_iaas/libvirt"
)

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
