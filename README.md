# 필수 패키지 설치
apk add qemu-system-x86_64 libvirt qemu-img

# 이건 앞으로 필요할
apk add novnc websockify

# libvirt 서비스 시작
rc-service libvirtd start
rc-update add libvirtd