package vm

type VM struct {
	Name  string `json:"name"`
	State uint   `json:"state"`
}

type CreateVMRequest struct {
	Name       string `json:"name"`
	Memory     uint   `json:"memory"`
	VCPU       uint   `json:"vcpu"`
	ImagePath  string `json:"image_path"`
	IsoPath    string `json:"iso_path"`
	DiskSizeGB uint   `json:"disk_size_gb"`
}
