package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	localLibvirt "udp_iaas/libvirt"
)

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
