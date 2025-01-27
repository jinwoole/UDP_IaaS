package types

import (
	"os/exec"
	"sync"
)

// WebsockifyManager manages websockify processes
type WebsockifyManager struct {
    processes map[int]*exec.Cmd
    mu        sync.Mutex
}

// CreateVMRequest defines the incoming JSON body for VM creation
type CreateVMRequest struct {
    Name    string `json:"name"`
    Cores   int    `json:"cores"`
    Memory  int    `json:"memory"`
    ISOName string `json:"iso"`
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

// VNCInfo represents VNC connection information
type VNCInfo struct {
    Port int    `json:"port"`
    Host string `json:"host"`
}