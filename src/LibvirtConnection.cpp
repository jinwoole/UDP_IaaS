#include "udp_iaas/LibvirtConnection.hpp"
#include <iostream>
#include <stdexcept>

namespace udp_iaas {

LibvirtConnection::LibvirtConnection() {
    conn = virConnectOpen("qemu:///system");
    if (!conn) {
        throw std::runtime_error("Failed to connect to QEMU/KVM");
    }
}

LibvirtConnection::~LibvirtConnection() {
    if (conn) virConnectClose(conn);
}

void LibvirtConnection::printConnectionInfo() const {
    const char* hostname = virConnectGetHostname(conn);
    if (hostname) {
        std::cout << "Connected to host: " << hostname << std::endl;
        free((void*)hostname);
    }

    unsigned long hvVer;
    if (virConnectGetVersion(conn, &hvVer) == 0) {
        std::cout << "Hypervisor version: "
                  << (hvVer/1000000) << "."
                  << (hvVer/1000) % 1000 << "."
                  << hvVer % 1000 << std::endl;
    }

    unsigned long libVer;
    if (virConnectGetLibVersion(conn, &libVer) == 0) {
        std::cout << "Libvirt version: "
                  << (libVer/1000000) << "."
                  << (libVer/1000) % 1000 << "."
                  << libVer % 1000 << std::endl;
    }
}

}  // namespace udp_iaas