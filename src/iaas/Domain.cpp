#include "udp_iaas/Domain.hpp"
#include <stdexcept>
#include <sstream>
#include <iostream>
#include <filesystem>

namespace udp_iaas {

Domain::Domain(virConnectPtr connection, const DomainConfig& config) 
    : conn(connection), domain(nullptr) {
    if (!create(config)) {
        throw std::runtime_error("Failed to create domain: " + config.name);
    }
}

Domain::Domain(virConnectPtr conn, virDomainPtr domain) 
    : conn(conn), domain(domain) {
    if (!domain) {
        throw std::runtime_error("Invalid domain pointer");
    }
}

Domain::~Domain() {
    if (domain) {
        virDomainFree(domain);
    }
}

bool Domain::create(const DomainConfig& config) {
    // 1. 먼저 디스크 이미지를 생성
    // qcow2 형식의 가상 디스크를 생성합니다
    if (!createDiskImage(config.diskPath, config.diskSizeGB)) {
        std::cerr << "Failed to create disk image" << std::endl;
        return false;
    }

    // 2. 도메인 XML 설정을 생성
    // libvirt는 XML 기반으로 VM 설정을 관리합니다
    std::string xmlConfig = generateXML(config);

    // 3. XML 설정으로 도메인을 정의
    // virDomainDefineXML: 영구적인 도메인 정의를 생성합니다
    // virDomainCreateXML: 임시 도메인을 생성하고 즉시 시작합니다
    domain = virDomainDefineXML(conn, xmlConfig.c_str());
    if (!domain) {
        std::cerr << "Failed to define domain from XML" << std::endl;
        return false;
    }

    // 4. 정의된 도메인을 시작
    // virDomainCreate: 비활성 상태의 도메인을 시작합니다
    if (virDomainCreate(domain) < 0) {
        std::cerr << "Failed to start domain" << std::endl;
        virDomainUndefine(domain);  // 도메인 제거
        return false;
    }

    return true;
}

bool Domain::shutdown() {
    // virDomainShutdown: ACPI 신호를 보내 게스트 OS를 정상적으로 종료
    if (virDomainShutdown(domain) < 0) {
        std::cerr << "Failed to shutdown domain" << std::endl;
        return false;
    }
    return true;
}

bool Domain::destroy() {
    // virDomainDestroy: 즉시 강제 종료 (kill)
    if (virDomainDestroy(domain) < 0) {
        std::cerr << "Failed to destroy domain" << std::endl;
        return false;
    }
    return true;
}

bool Domain::undefine() {
    // virDomainUndefine: 도메인 정의 제거
    if (virDomainUndefine(domain) < 0) {
        std::cerr << "Failed to undefine domain" << std::endl;
        return false;
    }
    return true;
}

bool Domain::remove() {
    // 먼저 강제 종료 후 정의 제거
    if (!destroy()) {
        return false;
    }
    return undefine();
}

std::string Domain::getISOPath(const std::string& isoName) {
    std::filesystem::path execPath = std::filesystem::current_path();
    std::filesystem::path isoPath = execPath / "iso" / isoName;
    return isoPath.string();
}


std::string Domain::generateXML(const DomainConfig& config) const {
    std::stringstream xml;
    
    xml << "<domain type='kvm'>\n"
        << "  <name>" << config.name << "</name>\n"
        << "  <memory unit='KiB'>" << config.memoryMB * 1024 << "</memory>\n"
        << "  <currentMemory unit='KiB'>" << config.memoryMB * 1024 << "</currentMemory>\n"
        << "  <vcpu placement='static'>" << config.vcpus << "</vcpu>\n"
        << "  <os>\n"
        << "    <type arch='x86_64' machine='pc'>hvm</type>\n";
    
    // ISO가 지정된 경우 CDROM에서 부팅
    if (config.isoPath) {
        xml << "    <boot dev='cdrom'/>\n"
            << "    <boot dev='hd'/>\n";
    } else {
        xml << "    <boot dev='hd'/>\n";
    }
    
    xml << "  </os>\n"
        << "  <cpu mode='host-model'/>\n"
        << "  <devices>\n"
        // 하드 디스크
        << "    <disk type='file' device='disk'>\n"
        << "      <driver name='qemu' type='qcow2'/>\n"
        << "      <source file='" << config.diskPath << "'/>\n"
        << "      <target dev='vda' bus='virtio'/>\n"
        << "    </disk>\n";

    // ISO 마운트를 위한 CDROM 추가
    if (config.isoPath) {
        xml << "    <disk type='file' device='cdrom'>\n"
            << "      <driver name='qemu' type='raw'/>\n"
            << "      <source file='" << config.isoPath.value() << "'/>\n"
            << "      <target dev='hdc' bus='ide'/>\n"
            << "      <readonly/>\n"
            << "    </disk>\n";
    }

    // 네트워크
    xml << "    <interface type='network'>\n"
        << "      <source network='default'/>\n"
        << "      <model type='virtio'/>\n"
        << "    </interface>\n"
        // 직렬 콘솔
        << "    <serial type='pty'>\n"
        << "      <target type='isa-serial' port='0'/>\n"
        << "    </serial>\n"
        << "    <console type='pty'>\n"
        << "      <target type='serial' port='0'/>\n"
        << "    </console>\n"
        // VNC 설정
        << "    <graphics type='vnc' port='-1' autoport='yes' listen='0.0.0.0'>\n"
        << "      <listen type='address' address='0.0.0.0'/>\n"
        << "    </graphics>\n"
        << "  </devices>\n"
        << "</domain>";

    return xml.str();
}

bool Domain::createDiskImage(const std::string& path, unsigned long sizeGB) const {
    // qemu-img 명령을 사용하여 디스크 이미지 생성
    // qcow2 포맷: Copy-On-Write 지원하는 QEMU 디스크 이미지 포맷
    std::string cmd = "qemu-img create -f qcow2 " + path + " " + 
                     std::to_string(sizeGB) + "G";
    
    return (system(cmd.c_str()) == 0);
}

bool Domain::openConsole() const {
    if (!domain) {
        std::cerr << "Domain not initialized\n";
        return false;
    }

    // virsh console 명령어 실행
    std::string cmd = "virsh console " + std::string(virDomainGetName(domain));
    return (system(cmd.c_str()) == 0);
}

} // namespace udp_iaas