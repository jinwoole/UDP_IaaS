package handlers

import (
	"net/http"
	"strings"
)

// handlers/handleVMs.go
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
 
	pathParts := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/vms/"), "/")
	if len(pathParts) == 0 || pathParts[0] == "" {
		http.Error(w, "VM name not specified", http.StatusBadRequest)
		return
	}
	vmName := pathParts[0]

	if len(pathParts) == 2 {
		switch pathParts[1] {
		case "stop":
			if r.Method == http.MethodPost {
				app.handleStopVM(w, r, vmName)
				return
			}
		case "start":
			if r.Method == http.MethodPost {
				app.handleStartVM(w, r, vmName)
				return
			}
		case "vnc":
			if r.Method == http.MethodGet {
				app.handleGetVNCPort(w, r, vmName)
				return
			}
		case "state":
            if r.Method == http.MethodGet {
                app.handleGetVMState(w, r, pathParts[0])
                return
            }
		}
		
		
	}
 
	// Check if it's a DELETE for that VM
	if len(pathParts) == 1 && r.Method == http.MethodDelete {
		app.handleDeleteVM(w, r, vmName)
		return
	}
 
	http.Error(w, "Invalid path", http.StatusNotFound)
 }