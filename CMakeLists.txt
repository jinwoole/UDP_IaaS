cmake_minimum_required(VERSION 3.10)
project(UDP_IaaS)

set(CMAKE_CXX_STANDARD 17)
set(CMAKE_CXX_STANDARD_REQUIRED ON)

find_package(PkgConfig REQUIRED)
pkg_check_modules(LIBVIRT REQUIRED libvirt)

include_directories(
    ${PROJECT_SOURCE_DIR}/include
    ${LIBVIRT_INCLUDE_DIRS}
)

add_executable(udp_iaas
    src/cli/CLI.cpp
    src/main.cpp
    src/iaas/LibvirtConnection.cpp
    src/iaas/Domain.cpp
)

target_link_libraries(udp_iaas PRIVATE 
    ${LIBVIRT_LIBRARIES}
    virt-qemu
)