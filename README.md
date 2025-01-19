# UDP_IaaS, Ultra-light Deployment Platform
알파인 리눅스 환경에서 구동되는 VM 가상화 플랫폼입니다.  
알파인 리눅스를 호스트OS로 초저사양 x86 홈서버에서 가상화를 지원하는 것을 목적으로 하며, 전혀 슈퍼하지 않습니다.  
가능하면 싱글 바이너리 파일로 만들어 설치를 좀 쉽게 하려고 합니다.  


* 설치 전 커뮤니티 저장소 활성화 필수  
```bash
apk add qemu-system libvirt libvirt-dev g++ make cmake pkgconfig

find /usr -name "libvirt.h"

pkg-config --cflags libvirt


rc-service libvirtd start

#kvm 확인
lsmod | grep kvm
#없으면 로드
modprobe kvm
modprobe kvm_intel  # Intel CPU의 경우
# or
modprobe kvm_amd    # AMD CPU의 경우

#kvm 호환성 확인
virt-host-validate

```