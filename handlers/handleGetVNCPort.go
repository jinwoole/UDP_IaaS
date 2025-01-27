package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	localLibvirt "udp_iaas/libvirt"
	"udp_iaas/types"
)

// handleGetVNCPort.go
func (app *App) handleGetVNCPort(w http.ResponseWriter, r *http.Request, vmName string) {
    // VNC 포트 조회
    vncPort, err := localLibvirt.GetVNCPort(app.Libvirt, vmName)
    if err != nil {
        http.Error(w, fmt.Sprintf("Failed to get VNC port: %v", err), http.StatusInternalServerError)
        return
    }

    // websockify 프로세스 시작
    wsPort, err := app.Websockify.Start(vncPort)
    if err != nil {
        http.Error(w, fmt.Sprintf("Failed to start websockify: %v", err), http.StatusInternalServerError)
        return
    }

    // 호스트 주소 추출
    host := r.Host
    if idx := strings.Index(host, ":"); idx != -1 {
        host = host[:idx]
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(types.VNCInfo{
        Port: wsPort,
        Host: host,
    })
}