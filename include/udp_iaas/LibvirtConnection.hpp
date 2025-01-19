#pragma once
#include "Domain.hpp"
#include <libvirt/libvirt.h>
#include <memory>

namespace udp_iaas {

class LibvirtConnection {
private:
    virConnectPtr conn;

public:
    LibvirtConnection();
    ~LibvirtConnection();
    
    void printConnectionInfo() const;

    // Domain 생성 메서드
    std::unique_ptr<Domain> createDomain(const DomainConfig& config) {
        return std::make_unique<Domain>(conn, config);
    }
};

} // namespace udp_iaas