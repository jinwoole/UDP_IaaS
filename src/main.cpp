#include "udp_iaas/LibvirtConnection.hpp"
#include <iostream>

int main() {
    try {
        udp_iaas::LibvirtConnection conn;
        conn.printConnectionInfo();

        // VM 설정
        udp_iaas::DomainConfig config{
            .name = "test-vm",
            .memoryMB = 2048,
            .vcpus = 2,
            .diskPath = "/var/lib/libvirt/images/test-vm.qcow2",
            .diskSizeGB = 20
        };

        // VM 생성
        auto domain = conn.createDomain(config);
        std::cout << "VM created successfully" << std::endl;

    } catch (const std::exception& e) {
        std::cerr << "Error: " << e.what() << std::endl;
        return 1;
    }
    return 0;
}