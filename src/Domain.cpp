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

std::string Domain::generateXML(const DomainConfig& config) const {
    // XML 생성을 위한 문자열 스트림
    std::stringstream xml;
    
    // 기본 XML 구조 생성
    // type='kvm': KVM 하이퍼바이저 사용
    // 각 요소들은 libvirt 문서에 정의된 스키마를 따릅니다
    xml << "<domain type='kvm'>\n"
        << "  <name>" << config.name << "</name>\n"
        // 메모리 설정 (KB 단위로 변환)
        << "  <memory unit='KiB'>" << config.memoryMB * 1024 << "</memory>\n"
        << "  <currentMemory unit='KiB'>" << config.memoryMB * 1024 << "</currentMemory>\n"
        // VCPU 설정
        << "  <vcpu placement='static'>" << config.vcpus << "</vcpu>\n"
        << "  <os>\n"
        // BIOS 타입 설정
        << "    <type arch='x86_64' machine='pc'>hvm</type>\n"
        << "    <boot dev='hd'/>\n"
        << "  </os>\n"
        // CPU 설정
        << "  <cpu mode='host-model'/>\n"
        // 디바이스 설정
        << "  <devices>\n"
        // 디스크 설정 (qcow2 포맷 사용)
        << "    <disk type='file' device='disk'>\n"
        << "      <driver name='qemu' type='qcow2'/>\n"
        << "      <source file='" << config.diskPath << "'/>\n"
        << "      <target dev='vda' bus='virtio'/>\n"
        << "    </disk>\n"
        // 기본 네트워크 설정 (NAT 모드)
        << "    <interface type='network'>\n"
        << "      <source network='default'/>\n"
        << "      <model type='virtio'/>\n"
        << "    </interface>\n"
        // 그래픽 설정 (VNC)
        << "    <graphics type='vnc' port='-1' autoport='yes' listen='127.0.0.1'/>\n"
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

} // namespace udp_iaas