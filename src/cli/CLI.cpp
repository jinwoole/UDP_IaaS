#include "cli/CLI.hpp"

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

    return config;
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

class ShutdownCommand : public Command {
private:
    LibvirtConnection& conn;
public:
    explicit ShutdownCommand(LibvirtConnection& conn) : conn(conn) {}
    void execute() override {
        std::cout << "Enter VM name to shutdown: ";
        std::string name;
        std::getline(std::cin, name);
        if (!conn.shutdownDomain(name))
            std::cerr << "Failed to shutdown VM\n";
        else
            std::cout << "VM shutdown issued successfully\n";
    }
    std::string help() const override {
        return "shutdown  Gracefully stop a virtual machine";
    }
};

class UndefineCommand : public Command {
private:
    LibvirtConnection& conn;
public:
    explicit UndefineCommand(LibvirtConnection& conn) : conn(conn) {}
    void execute() override {
        std::cout << "Enter VM name to undefine: ";
        std::string name;
        std::getline(std::cin, name);
        if (!conn.undefineDomain(name))
            std::cerr << "Failed to undefine VM\n";
        else
            std::cout << "VM undefined successfully\n";
    }
    std::string help() const override {
        return "undefine  Remove a virtual machine definition";
    }
};

class RemoveCommand : public Command {
private:
    LibvirtConnection& conn;
public:
    explicit RemoveCommand(LibvirtConnection& conn) : conn(conn) {}
    void execute() override {
        std::cout << "Enter VM name to remove: ";
        std::string name;
        std::getline(std::cin, name);
        if (!conn.removeDomain(name))
            std::cerr << "Failed to remove VM\n";
        else
            std::cout << "VM removed successfully\n";
    }
    std::string help() const override {
        return "remove    Force destroy and undefine a virtual machine";
    }
};

// CLI implementations
CLI::CLI() {
    // 커맨드 등록
    commands["create"] = std::make_unique<CreateCommand>(conn);
    commands["destroy"] = std::make_unique<DestroyCommand>(conn);
    commands["shutdown"] = std::make_unique<ShutdownCommand>(conn);
    commands["undefine"] = std::make_unique<UndefineCommand>(conn);
    commands["remove"]   = std::make_unique<RemoveCommand>(conn);
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
