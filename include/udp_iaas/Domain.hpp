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

    // VM 제어
    bool create(const DomainConfig& config);
    bool shutdown(); //종료
    bool destroy(); //강제종료
    bool undefine(); //정의 제거
    bool remove(); // destory + undefine




private:
    // XML 설정 생성
    std::string generateXML(const DomainConfig& config) const;
    
    // 디스크 이미지 생성
    bool createDiskImage(const std::string& path, unsigned long sizeGB) const;
};

} // namespace udp_iaas