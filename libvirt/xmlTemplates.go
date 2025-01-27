package libvirt

import "fmt"

func generateVMXML(name string, memory int, cores int, diskPath string, isoPath string) string {
	cdromXML := ""
	if isoPath != "" {
		cdromXML = fmt.Sprintf(`
        <disk type='file' device='cdrom'>
            <driver name='qemu' type='raw'/>
            <source file='%s'/>
            <target dev='hdc' bus='ide'/>
            <readonly/>
        </disk>`, isoPath)
	}

	return fmt.Sprintf(`
<domain type='kvm'>
    <name>%s</name>
    <memory unit='MiB'>%d</memory>
    <vcpu>%d</vcpu>
    <os>
        <type arch='x86_64'>hvm</type>
        <boot dev='cdrom'/>
        <boot dev='hd'/>
    </os>
    <devices>
        <disk type='file' device='disk'>
            <driver name='qemu' type='raw'/>
            <source file='%s'/>
            <target dev='vda' bus='virtio'/>
        </disk>
        %s
        <interface type='network'>
            <source network='default'/>
            <model type='virtio'/>
        </interface>
        <graphics type='vnc' port='-1' autoport='yes' listen='0.0.0.0'>
            <listen type='address' address='0.0.0.0'/>
        </graphics>
        <video>
            <model type='vga'/>
        </video>
    </devices>
</domain>`,
		name,
		memory,
		cores,
		diskPath,
		cdromXML,
	)
}