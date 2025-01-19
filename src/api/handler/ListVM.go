package handler

import (
	"encoding/xml"
	"net/http"
	"udp_iaas/src/libvirt"

	"github.com/labstack/echo/v4"
)

type VMHandler struct {
	client *libvirt.Client
}

func NewVMHandler(client *libvirt.Client) *VMHandler {
	return &VMHandler{client: client}
}

func (h *VMHandler) ListVMs(c echo.Context) error {
	domains, err := h.client.GetConnection().ListAllDomains(0)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
	}

	vms := make([]map[string]interface{}, 0)
	for _, domain := range domains {
		name, err := domain.GetName()
		if err != nil {
			continue
		}

		// 기본 상태 정보
		state, _, err := domain.GetState()
		if err != nil {
			continue
		}
		stateStr := "unknown"
		switch state {
		case 1:
			stateStr = "running"
		case 3:
			stateStr = "paused"
		case 4:
			stateStr = "shutdown"
		case 5:
			stateStr = "shutoff"
		}

		// CPU, 메모리 등 리소스 정보
		info, err := domain.GetInfo()
		if err != nil {
			continue
		}

		// 디스크 정보 가져오기
		diskXML, err := domain.GetXMLDesc(0)
		if err != nil {
			continue
		}

		// XML 파싱하여 디스크, ISO 정보 추출
		type disk struct {
			Device string `xml:"device,attr"`
			Source struct {
				File string `xml:"file,attr"`
			} `xml:"source"`
		}
		type devices struct {
			Disks []disk `xml:"disk"`
		}
		type domainXML struct {
			Devices devices `xml:"devices"`
		}

		var xmlDoc domainXML
		if err := xml.Unmarshal([]byte(diskXML), &xmlDoc); err != nil {
			continue
		}

		// 디스크와 ISO 정보 분리
		var disks []string
		var isos []string
		for _, d := range xmlDoc.Devices.Disks {
			if d.Device == "disk" {
				disks = append(disks, d.Source.File)
			} else if d.Device == "cdrom" && d.Source.File != "" {
				isos = append(isos, d.Source.File)
			}
		}

		var memoryUsed uint64
		if stateStr == "shutoff" {
			memoryUsed = 0
		} else {
			memoryUsed = info.Memory / 1024
		}

		vmInfo := map[string]interface{}{
			"name":         name,
			"state":        stateStr,
			"vcpus":        info.NrVirtCpu,
			"memory_max":   info.MaxMem / 1024, // KB를 MB로 변환
			"memory_used":  memoryUsed,
			"disks":        disks,
			"mounted_isos": isos,
		}

		vms = append(vms, vmInfo)
		domain.Free()
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"vms": vms,
	})
}
