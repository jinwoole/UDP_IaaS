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

    bool shutdownDomain(const std::string& name);
    bool destroyDomain(const std::string& name);
    bool undefineDomain(const std::string& name);

    // removeDomain: force-destroy + undefine
    bool removeDomain(const std::string& name);
};

} // namespace udp_iaas