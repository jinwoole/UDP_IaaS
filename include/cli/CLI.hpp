#pragma once
#include <string>
#include <memory>
#include <map>
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
    DomainConfig getUserInput() const;

public:
    explicit CreateCommand(LibvirtConnection& conn);
    void execute() override;
    std::string help() const override;
};

// Start command
class StartCommand : public Command {
private:
    LibvirtConnection& conn;
public:
    explicit StartCommand(LibvirtConnection& conn);
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

// Shutdown command
class ShutdownCommand : public Command {
private:
    LibvirtConnection& conn;
public:
    explicit ShutdownCommand(LibvirtConnection& conn);
    void execute() override;
    std::string help() const override;
};

// Undefine command
class UndefineCommand : public Command {
private:
    LibvirtConnection& conn;
public:
    explicit UndefineCommand(LibvirtConnection& conn);
    void execute() override;
    std::string help() const override;
};

// Remove command
class RemoveCommand : public Command {
private:
    LibvirtConnection& conn;
public:
    explicit RemoveCommand(LibvirtConnection& conn);
    void execute() override;
    std::string help() const override;
};

// Console command
class ConsoleCommand : public Command {
private:
    LibvirtConnection& conn;
public:
    explicit ConsoleCommand(LibvirtConnection& conn);
    void execute() override;
    std::string help() const override;
};

// Reboot command
class RebootCommand : public Command {
private:
    LibvirtConnection& conn;
public:
    explicit RebootCommand(LibvirtConnection& conn);
    void execute() override;
    std::string help() const override;
};

// CDROM 관리 command
class CDROMCommand : public Command {
private:
    LibvirtConnection& conn;
public:
    explicit CDROMCommand(LibvirtConnection& conn);
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