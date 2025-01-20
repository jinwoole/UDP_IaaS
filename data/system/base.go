package system

import (
	"libvirt.org/go/libvirt"
)

type Data struct {
	conn *libvirt.Connect
}

func NewData(conn *libvirt.Connect) *Data {
	return &Data{
		conn: conn,
	}
}
