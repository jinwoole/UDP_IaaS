package handlers

import (
	"net/http"
)

func (app *App) HandleISOs(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		app.handleListISOs(w, r)
	case http.MethodPost:
		app.handleUploadISO(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
