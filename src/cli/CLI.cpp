#include "cli/CLI.hpp"
#include <libvirt/libvirt-qemu.h>
#include <iostream>
#include <filesystem>
#include <array>
#include <vector>
#include <cstdio>
#include <memory>
#include <stdexcept>
#include <string>
#include <algorithm>
#include <unistd.h>

namespace udp_iaas {

// CreateCommand implementations
CreateCommand::CreateCommand(LibvirtConnection& conn) : conn(conn) {}

DomainConfig CreateCommand::getUserInput() const {
    DomainConfig config;
    
    std::cout << "Enter VM name: ";
    std::getline(std::cin, config.name);
    
    std::cout << "Enter memory size (MB): ";
    std::cin >> config.memoryMB;
    
    std::cout << "Enter number of vCPUs: ";
    std::cin >> config.vcpus;
    
    std::cin.ignore(); // 버퍼 클리어
    
    config.diskPath = "/var/lib/libvirt/images/" + config.name + ".qcow2";
    
    std::cout << "Enter disk size (GB): ";
    std::cin >> config.diskSizeGB;
    
    std::cin.ignore(); // 버퍼 클리어
    
    std::cout << "Do you want to install from ISO? (y/n): ";
    std::string response;
    std::getline(std::cin, response);
    
    if (response == "y" || response == "Y") {
        std::cout << "Available ISOs in ./iso directory:\n";
        std::filesystem::path isoDir = std::filesystem::current_path() / "iso";
        
        std::vector<std::string> isoFiles;
        if (std::filesystem::exists(isoDir)) {
            for (const auto& entry : std::filesystem::directory_iterator(isoDir)) {
                if (entry.path().extension() == ".iso") {
                    isoFiles.push_back(entry.path().filename().string());
                    std::cout << isoFiles.size() << ". " << entry.path().filename().string() << "\n";
                }
            }
        }
        
        if (isoFiles.empty()) {
            std::cout << "No ISO files found. Place ISO files in ./iso directory.\n";
        } else {
            std::cout << "Select ISO number (0 to skip): ";
            int selection;
            std::cin >> selection;
            std::cin.ignore();
            
            if (selection > 0 && selection <= static_cast<int>(isoFiles.size())) {
                config.isoPath = Domain::getISOPath(isoFiles[selection - 1]);
            }
        }
    }

    return config;
}

// StartCommand implementations
StartCommand::StartCommand(LibvirtConnection& conn) : conn(conn) {}

void StartCommand::execute() {
    std::cout << "Enter VM name to start: ";
    std::string name;
    std::getline(std::cin, name);

    virDomainPtr domain = virDomainLookupByName(conn.getConnection(), name.c_str());
    if (!domain) {
        std::cerr << "Failed to find VM: " << name << "\n";
        return;
    }

    if (virDomainCreate(domain) < 0) {
        std::cerr << "Failed to start VM\n";
    } else {
        std::cout << "VM started successfully\n";
    }

    virDomainFree(domain);
}

std::string StartCommand::help() const {
    return "start     Start a virtual machine";
}

void CreateCommand::execute() {
    auto config = getUserInput();
    auto domain = conn.createDomain(config);
    std::cout << "VM created successfully\n";
}

std::string CreateCommand::help() const {
    return "create    Create a new virtual machine";
}

// DestroyCommand implementations
DestroyCommand::DestroyCommand(LibvirtConnection& conn) : conn(conn) {}

void DestroyCommand::execute() {
    std::cout << "Enter VM name to destroy: ";
    std::string name;
    std::getline(std::cin, name);
    if (!conn.destroyDomain(name))
        std::cerr << "Failed to destroy VM\n";
    else
        std::cout << "VM destroyed successfully\n";
}

std::string DestroyCommand::help() const {
    return "destroy   Forcefully stop a virtual machine";
}

// ShutdownCommand implementations
ShutdownCommand::ShutdownCommand(LibvirtConnection& conn) : conn(conn) {}

void ShutdownCommand::execute() {
    std::cout << "Enter VM name to shutdown: ";
    std::string name;
    std::getline(std::cin, name);

    virDomainPtr domain = virDomainLookupByName(conn.getConnection(), name.c_str());
    if (!domain) {
        std::cerr << "Failed to find VM: " << name << "\n";
        return;
    }

    std::cout << "Shutdown methods:\n";
    std::cout << "1. ACPI Power Button (default)\n";
    std::cout << "2. Guest Agent Shutdown\n";
    std::cout << "3. Initctl Shutdown\n";
    std::cout << "4. Signal Shutdown\n";
    std::cout << "5. QEMU Monitor Power Down\n";  // 추가
    std::cout << "Select method (1-5): ";

    std::string choice;
    std::getline(std::cin, choice);

    bool success = false;
    if (choice == "5") {
        // QEMU monitor command를 사용한 종료
        const char* cmd = "{\"execute\": \"system_powerdown\" }";
        char *result = nullptr;
        
        if (virDomainQemuMonitorCommand(domain, cmd, &result, 0) < 0) {
            std::cerr << "Failed to send QEMU monitor command\n";
        } else {
            std::cout << "Power down command sent via QEMU monitor\n";
            success = true;
        }
        
        if (result) free(result);
    } else {
        // 기존 방식들
        unsigned int flags = 0;
        if (choice == "2")
            flags = VIR_DOMAIN_SHUTDOWN_GUEST_AGENT;
        else if (choice == "3")
            flags = VIR_DOMAIN_SHUTDOWN_INITCTL;
        else if (choice == "4")
            flags = VIR_DOMAIN_SHUTDOWN_SIGNAL;
        
        if (virDomainShutdownFlags(domain, flags) < 0) {
            std::cerr << "Failed to shutdown VM\n";
        } else {
            std::cout << "VM shutdown signal sent successfully\n";
            success = true;
        }
    }

    // 종료 명령이 성공했을 경우 상태 확인
    if (success) {
        std::cout << "Waiting for VM to shutdown";
        int timeout = 30;  // 30초 타임아웃
        while (timeout > 0) {
            virDomainInfo info;
            if (virDomainGetInfo(domain, &info) == 0) {
                if (info.state == VIR_DOMAIN_SHUTOFF) {
                    std::cout << "\nVM shutdown completed\n";
                    break;
                }
            }
            std::cout << ".";
            std::cout.flush();
            sleep(1);
            timeout--;
        }
        if (timeout == 0) {
            std::cout << "\nVM shutdown timed out\n";
        }
    }

    virDomainFree(domain);
}

std::string ShutdownCommand::help() const {
    return "shutdown  Gracefully stop a virtual machine";
}

// UndefineCommand implementations
UndefineCommand::UndefineCommand(LibvirtConnection& conn) : conn(conn) {}

void UndefineCommand::execute() {
    std::cout << "Enter VM name to undefine: ";
    std::string name;
    std::getline(std::cin, name);
    if (!conn.undefineDomain(name))
        std::cerr << "Failed to undefine VM\n";
    else
        std::cout << "VM undefined successfully\n";
}

std::string UndefineCommand::help() const {
    return "undefine  Remove a virtual machine definition";
}

// RemoveCommand implementations
RemoveCommand::RemoveCommand(LibvirtConnection& conn) : conn(conn) {}

void RemoveCommand::execute() {
    std::cout << "Enter VM name to remove: ";
    std::string name;
    std::getline(std::cin, name);
    if (!conn.removeDomain(name))
        std::cerr << "Failed to remove VM\n";
    else
        std::cout << "VM removed successfully\n";
}

std::string RemoveCommand::help() const {
    return "remove    Force destroy and undefine a virtual machine";
}

// RebootCommand implementations
RebootCommand::RebootCommand(LibvirtConnection& conn) : conn(conn) {}

void RebootCommand::execute() {
    std::cout << "Enter VM name to reboot: ";
    std::string name;
    std::getline(std::cin, name);

    virDomainPtr domain = virDomainLookupByName(conn.getConnection(), name.c_str());
    if (!domain) {
        std::cerr << "Failed to find VM: " << name << "\n";
        return;
    }

    if (virDomainReboot(domain, 0) < 0) {
        std::cerr << "Failed to reboot VM\n";
    } else {
        std::cout << "VM reboot initiated successfully\n";
    }

    virDomainFree(domain);
}

std::string RebootCommand::help() const {
    return "reboot    Reboot a virtual machine";
}

// ConsoleCommand implementations
ConsoleCommand::ConsoleCommand(LibvirtConnection& conn) : conn(conn) {}

void ConsoleCommand::execute() {
    std::cout << "Enter VM name to connect to console: ";
    std::string name;
    std::getline(std::cin, name);
    
    virDomainPtr domain = virDomainLookupByName(conn.getConnection(), name.c_str());
    if (!domain) {
        std::cerr << "Failed to find VM: " << name << "\n";
        return;
    }

    // Get VNC port
    std::string cmd = "virsh vncdisplay " + name;
    std::array<char, 128> buffer;
    std::string result;
    std::unique_ptr<FILE, decltype(&pclose)> pipe(popen(cmd.c_str(), "r"), pclose);
    if (!pipe) {
        std::cerr << "Failed to execute command\n";
        virDomainFree(domain);
        return;
    }
    while (fgets(buffer.data(), buffer.size(), pipe.get()) != nullptr) {
        result += buffer.data();
    }
    
    result.erase(std::remove(result.begin(), result.end(), '\n'), result.end());

    int vncPort = 5900;
    if (!result.empty() && result[0] == ':') {
        vncPort += std::stoi(result.substr(1));
    }

    // Get host IP address
    cmd = "hostname -I | awk '{print $1}'";
    pipe.reset(popen(cmd.c_str(), "r"));
    std::string hostIP;
    if (pipe && fgets(buffer.data(), buffer.size(), pipe.get()) != nullptr) {
        hostIP = buffer.data();
        hostIP.erase(std::remove(hostIP.begin(), hostIP.end(), '\n'), hostIP.end());
    } else {
        hostIP = "<your-server-ip>";
    }

    std::cout << "\nVM Console Access Options:\n";
    std::cout << "------------------------\n";
    std::cout << "1. VNC Connection:\n";
    std::cout << "   Local access:  localhost:" << vncPort << "\n";
    std::cout << "   Remote access: " << hostIP << ":" << vncPort << "\n";
    std::cout << "   Use any VNC viewer to connect\n\n";
    
    std::cout << "2. Web Browser Access:\n";
    std::cout << "   First, start noVNC server with:\n";
    std::cout << "   $ websockify --web=/usr/share/novnc --listen 0.0.0.0:6080 localhost:" << vncPort << "\n\n";
    std::cout << "   Then access via:\n";
    std::cout << "   Local:  http://localhost:6080/vnc.html\n";
    std::cout << "   Remote: http://" << hostIP << ":6080/vnc.html\n\n";
    
    std::cout << "3. Text Console (after OS installation):\n";
    std::cout << "   $ virsh console " << name << "\n";
    std::cout << "   (Use Ctrl+] to exit)\n\n";

    std::cout << "Note: Make sure your firewall allows incoming connections to ports " 
              << vncPort << " (VNC) and 6080 (noVNC)\n";

    virDomainFree(domain);
}

std::string ConsoleCommand::help() const {
    return "console   Connect to VM's console";
}

void CDROMCommand::execute() {
    std::cout << "Enter VM name: ";
    std::string name;
    std::getline(std::cin, name);

    virDomainPtr domain = virDomainLookupByName(conn.getConnection(), name.c_str());
    if (!domain) {
        std::cerr << "Failed to find VM: " << name << "\n";
        return;
    }

    std::cout << "CDROM operations:\n";
    std::cout << "1. Eject ISO\n";
    std::cout << "2. Insert ISO\n";
    std::cout << "Select operation (1-2): ";
    
    std::string choice;
    std::getline(std::cin, choice);

    // 플래그 설정: 실행 중이면 LIVE를, 꺼져있으면 CONFIG를 사용
    unsigned int flags = 0;
    virDomainInfo info;
    if (virDomainGetInfo(domain, &info) == 0) {
        if (info.state == VIR_DOMAIN_RUNNING)
            flags |= VIR_DOMAIN_AFFECT_LIVE;
        flags |= VIR_DOMAIN_AFFECT_CONFIG;  // 항상 설정은 저장
    }

    if (choice == "1") {
        // CDROM 분리
        if (virDomainDetachDeviceFlags(domain, "<disk device='cdrom'><target dev='hdc'/></disk>", flags) < 0) {
            std::cerr << "Failed to eject CDROM\n";
        } else {
            std::cout << "CDROM ejected successfully\n";
        }
    }
    else if (choice == "2") {
        // ISO 목록 표시 및 선택
        std::cout << "Available ISOs in ./iso directory:\n";
        std::filesystem::path isoDir = std::filesystem::current_path() / "iso";
        
        std::vector<std::string> isoFiles;
        if (std::filesystem::exists(isoDir)) {
            for (const auto& entry : std::filesystem::directory_iterator(isoDir)) {
                if (entry.path().extension() == ".iso") {
                    isoFiles.push_back(entry.path().string());
                    std::cout << isoFiles.size() << ". " << entry.path().filename().string() << "\n";
                }
            }
        }
        
        if (isoFiles.empty()) {
            std::cout << "No ISO files found in ./iso directory\n";
        } else {
            std::cout << "Select ISO number: ";
            std::string selection;
            std::getline(std::cin, selection);
            
            try {
                int idx = std::stoi(selection) - 1;
                if (idx >= 0 && idx < static_cast<int>(isoFiles.size())) {
                    // CDROM 장치 연결
                    std::string cdrom_xml = "<disk type='file' device='cdrom'>"
                                          "<driver name='qemu' type='raw'/>"
                                          "<source file='" + isoFiles[idx] + "'/>"
                                          "<target dev='hdc' bus='ide'/>"
                                          "<readonly/>"
                                          "</disk>";

                    if (virDomainAttachDeviceFlags(domain, cdrom_xml.c_str(), flags) < 0) {
                        std::cerr << "Failed to attach ISO\n";
                    } else {
                        std::cout << "ISO attached successfully\n";
                    }
                }
            } catch (const std::exception& e) {
                std::cerr << "Invalid selection\n";
            }
        }
    }

    // 현재 상태 확인
    std::string check_cmd = "virsh domblklist " + name + " | grep hdc";
    system(check_cmd.c_str());

    virDomainFree(domain);
}
// CDROMCommand 생성자 구현
CDROMCommand::CDROMCommand(LibvirtConnection& conn) : conn(conn) {}

// CDROMCommand help 함수 구현
std::string CDROMCommand::help() const {
    return "cdrom     Manage VM's CDROM (eject/insert ISO)";
}

// CLI implementations
CLI::CLI() {
    commands["create"] = std::make_unique<CreateCommand>(conn);
    commands["destroy"] = std::make_unique<DestroyCommand>(conn);
    commands["shutdown"] = std::make_unique<ShutdownCommand>(conn);
    commands["reboot"] = std::make_unique<RebootCommand>(conn);
    commands["start"] = std::make_unique<StartCommand>(conn);
    commands["undefine"] = std::make_unique<UndefineCommand>(conn);
    commands["remove"] = std::make_unique<RemoveCommand>(conn);
    commands["console"] = std::make_unique<ConsoleCommand>(conn);
    commands["cdrom"] = std::make_unique<CDROMCommand>(conn);
}
void CLI::run(int argc, char* argv[]) {
    if (argc < 2 || std::string(argv[1]) == "help") {
        printHelp();
        return;
    }

    std::string command = argv[1];
    auto it = commands.find(command);
    if (it != commands.end()) {
        it->second->execute();
    } else {
        std::cerr << "Unknown command: " << command << "\n";
        printHelp();
    }
}

void CLI::printHelp() const {
    std::cout << "Usage: udp_iaas <command>\n\n";
    std::cout << "Available commands:\n";
    for (const auto& [name, cmd] : commands) {
        std::cout << cmd->help() << "\n";
    }
}

} // namespace udp_iaas