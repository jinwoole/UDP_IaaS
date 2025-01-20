package vm

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Handler) List(c echo.Context) error {
	vms, err := h.service.GetAllVMs()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}
	return c.JSON(http.StatusOK, vms)
}
