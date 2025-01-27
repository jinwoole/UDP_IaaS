package libvirt

const (
	isoStoragePath = "/var/lib/vms/isos"
)

func getISOPath(isoName string) string {
	return isoStoragePath + "/" + isoName
}