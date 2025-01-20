package vm

import (
	"net/http"
	vmData "udp_iaas/data/vm"

	"github.com/labstack/echo/v4"
)

// @Summary Create new VM
// @Description Create a new virtual machine
// @Tags VM
// @Accept json
// @Produce json
// @Param request body vmData.CreateVMRequest true "VM Creation Request"
// @Success 201 {object} map[string]string "VM created successfully"
// @Failure 400 {object} map[string]string "Invalid request format"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /vm/create [post]
func (h *Handler) Create(c echo.Context) error {
	var req vmData.CreateVMRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request format",
		})
	}

	err := h.service.CreateVM(req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, map[string]string{
		"message": "VM created successfully",
	})
}
