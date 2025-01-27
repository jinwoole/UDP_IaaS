package handlers

import (
	"udp_iaas/types"

	"libvirt.org/go/libvirt"
)

type App struct {
    Libvirt    *libvirt.Connect
    Websockify *types.WebsockifyManager
}