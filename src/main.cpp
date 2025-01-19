#include "udp_iaas/LibvirtConnection.hpp"
#include <iostream>

int main() {
   try {
       udp_iaas::LibvirtConnection conn;
       conn.printConnectionInfo();

       // VM 설정
       udp_iaas::DomainConfig config{
           .name = "test-vm",
           .memoryMB = 1024,    // 2GB RAM
           .vcpus = 1,          // 2 vCPUs
           .diskPath = "/var/lib/libvirt/images/test-vm.qcow2",
           .diskSizeGB = 20     // 20GB 디스크
       };

       // VM 생성
       std::cout << "Creating VM..." << std::endl;
       auto domain = conn.createDomain(config);
       std::cout << "VM created successfully" << std::endl;

       // 사용자 입력 대기
       std::cout << "\nPress Enter to remove the VM..." << std::endl;
       std::cin.get();

       // VM 제거
       if (domain->remove()) {
           std::cout << "VM removed successfully" << std::endl;
       } else {
           std::cerr << "Failed to remove VM" << std::endl;
       }

   } catch (const std::exception& e) {
       std::cerr << "Error: " << e.what() << std::endl;
       return 1;
   }
   return 0;
}