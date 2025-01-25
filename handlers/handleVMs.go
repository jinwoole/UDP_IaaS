package handlers

import (
	"net/http"
	"strings"
)

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
