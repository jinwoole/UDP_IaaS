package system

import (
	"net/http"

	systemData "udp_iaas/data/system"

	"github.com/labstack/echo/v4"
)

// Define a structured error response
type ErrorResponse struct {
	Error string `json:"error"`
}

// GetISOInfo retrieves ISO information.
//
// @Summary Get ISO List
// @Description Get list of all available ISO images
// @Tags ISO
// @Accept json
// @Produce json
// @Success 200 {array} systemData.ISOInfo "List of ISO images"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /system/info/iso [get]
func (h *Handler) GetISOInfo(c echo.Context) error {
	var isos []systemData.ISOInfo
	isos, err := h.service.GetISO()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, isos)
}
