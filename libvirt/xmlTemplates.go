package libvirt

import "fmt"

func generateVMXML(name string, memory int, cores int, diskPath string, isoPath string) string {
    cdromXML := ""
    if isoPath != "" {
        cdromXML = fmt.Sprintf(`
        <disk type='file' device='cdrom'>
            <driver name='qemu' type='raw' cache='none' io='native'/>
            <source file='%s'/>
            <target dev='sda' bus='sata'/>
            <readonly/>
        </disk>`, isoPath)
    }
 
    return fmt.Sprintf(`
 <domain type='kvm'>
    <name>%s</name>
    <memory unit='MiB'>%d</memory>
    <vcpu placement='static'>%d</vcpu>
    <cpu mode='host-passthrough'>
        <cache mode='passthrough'/>
        <feature policy='require' name='vmx'/>
    </cpu>
    <features>
        <acpi/>
        <apic/>
        <kvm>
            <hidden state='on'/>
        </kvm>
    </features>
    <os>
        <type arch='x86_64' machine='pc-q35-8.0'>hvm</type>
        <boot dev='cdrom'/>
        <boot dev='hd'/>
    </os>
    <devices>
        <disk type='file' device='disk'>
            <driver name='qemu' type='raw' cache='none' io='native' discard='unmap'/>
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
            <model type='virtio' heads='1'/>
        </video>
    </devices>
    <iothreads>1</iothreads>
 </domain>`,
        name,
        memory,
        cores,
        diskPath,
        cdromXML,
    )
 }