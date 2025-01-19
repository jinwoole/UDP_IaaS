# UDP_IaaS, Ultra-light Deployment Platform
알파인 리눅스 환경에서 구동되는 VM 가상화 플랫폼입니다.  
알파인 리눅스를 호스트OS로 초저사양 x86 홈서버에서 가상화를 지원하는 것을 목적으로 하며, 전혀 슈퍼하지 않습니다. 일단 클라우드 네이티브 하게 꾸며보려 합니다.  


* 알파인 환경설정: 설치 전 커뮤니티 저장소 활성화 필수  
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

Handler 비즈니스 로직 처리 담당  
Router URL 경로와 핸들러, 미들웨어 설정  
main 의존성 주입, 서버 시작  