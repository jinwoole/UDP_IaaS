package handler

import (
	"net/http"
	"udp_iaas/src/libvirt"

	"github.com/labstack/echo/v4"
)

type StatusHandler struct {
	client *libvirt.Client
}

func NewStatusHandler(client *libvirt.Client) *StatusHandler {
	return &StatusHandler{client: client}
}

func (h *StatusHandler) GetStatus(c echo.Context) error {
	isConnected, err := h.client.IsConnected()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
	}

	hostname, err := h.client.GetConnection().GetHostname()
	status := map[string]interface{}{
		"connected": isConnected,
	}
	if err == nil {
		status["hostname"] = hostname
	}

	return c.JSON(http.StatusOK, status)
}
