#pragma once
#include <libvirt/libvirt.h>
#include <string>
#include <memory>
#include <optional>
#include <filesystem>

namespace udp_iaas {

struct DomainConfig {
    std::string name;
    unsigned int memoryMB;
    unsigned int vcpus;
    std::string diskPath;
    unsigned long diskSizeGB;
    std::optional<std::string> isoPath;  // ISO 파일 경로 (선택적)
};

class Domain {
private:
    virDomainPtr domain;
    virConnectPtr conn;  // 참조용 포인터, 소유권 없음

public:
    // 새로운 VM 생성을 위한 생성자
    Domain(virConnectPtr conn, const DomainConfig& config);
    
    // 기존 VM을 래핑하기 위한 생성자
    Domain(virConnectPtr conn, virDomainPtr domain);
    
    ~Domain();

    // VM 제어
    bool create(const DomainConfig& config);
    bool shutdown();
    bool destroy();
    bool undefine();
    bool remove();
    bool openConsole() const;

    static std::string getISOPath(const std::string& isoName);


private:
    std::string generateXML(const DomainConfig& config) const;
    bool createDiskImage(const std::string& path, unsigned long sizeGB) const;
};

} // namespace udp_iaas