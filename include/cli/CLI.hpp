// include/udp_iaas/CLI.hpp
#pragma once
#include <iostream>
#include <string>
#include <memory>
#include <map>
#include <vector>
#include <functional>
#include "udp_iaas/LibvirtConnection.hpp"

namespace udp_iaas {

// 기본 Command 클래스
class Command {
public:
    virtual ~Command() = default;
    virtual void execute() = 0;
    virtual std::string help() const = 0;
};

// Create command
class CreateCommand : public Command {
private:
    LibvirtConnection& conn;

    // 사용자로부터 VM 설정 정보 얻기
    DomainConfig getUserInput() const;

public:
    explicit CreateCommand(LibvirtConnection& conn);
    
    void execute() override;
    std::string help() const override;
};

// Destroy command
class DestroyCommand : public Command {
private:
    LibvirtConnection& conn;

public:
    explicit DestroyCommand(LibvirtConnection& conn);
    
    void execute() override;
    std::string help() const override;
};

// CLI 관리 클래스
class CLI {
private:
    LibvirtConnection conn;
    std::map<std::string, std::unique_ptr<Command>> commands;

public:
    CLI();

    void run(int argc, char* argv[]);
    void printHelp() const;
};

} // namespace udp_iaas