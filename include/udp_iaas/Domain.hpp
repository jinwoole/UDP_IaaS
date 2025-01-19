#pragma once
#include <libvirt/libvirt.h>
#include <string>
#include <memory>

namespace udp_iaas {

struct DomainConfig {
    std::string name;
    unsigned int memoryMB;
    unsigned int vcpus;
    std::string diskPath;
    unsigned long diskSizeGB;
};

class Domain {
private:
    virDomainPtr domain;
    virConnectPtr conn;  // 참조용 포인터, 소유권 없음

public:
    Domain(virConnectPtr conn, const DomainConfig& config);
    ~Domain();

    // 새로운 VM 생성
    bool create(const DomainConfig& config);

private:
    // XML 설정 생성
    std::string generateXML(const DomainConfig& config) const;
    
    // 디스크 이미지 생성
    bool createDiskImage(const std::string& path, unsigned long sizeGB) const;
};

} // namespace udp_iaas