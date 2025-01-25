package types

import "libvirt.org/go/libvirt"

// App holds references to external services or connections
type App struct {
	Libvirt *libvirt.Connect
}

// CreateVMRequest defines the incoming JSON body for VM creation
type CreateVMRequest struct {
	Name    string `json:"name"`
	Cores   int    `json:"cores"`
	Memory  int    `json:"memory"`
	ISOName string `json:"iso"` // Changed from ISO to ISOName to match usage
}

// VM represents a virtual machine status
type VM struct {
	Name   string `json:"name"`
	Cores  int    `json:"cores"`
	Memory int    `json:"memory"`
	State  string `json:"state"`
}

// ISO describes an uploaded ISO file
type ISO struct {
	Name string `json:"name"`
	Size int64  `json:"size"`
}
