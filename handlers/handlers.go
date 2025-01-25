package handlers

import (
	"libvirt.org/go/libvirt"
)

// App holds references to external services or connections
type App struct {
	Libvirt *libvirt.Connect
}
