package handlers

import (
	"encoding/json"
	"net/http"
)

// 헬스체크 엔드포인트
func (app *App) HandleHealthCheck(w http.ResponseWriter, r *http.Request) {
	// libvirt 연결 상태 확인
	_, err := app.Libvirt.GetVersion()
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(map[string]string{
			"status": "unhealthy",
			"error":  "libvirt connection failed",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "healthy",
	})
}
