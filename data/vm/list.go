package vm

import (
	"libvirt.org/go/libvirt"
)

func (d *Data) GetAllDomains() ([]libvirt.Domain, error) {
	return d.conn.ListAllDomains(libvirt.CONNECT_LIST_DOMAINS_ACTIVE |
		libvirt.CONNECT_LIST_DOMAINS_INACTIVE)
}
