package main

import (
	"log"
	"net/http"
	"os"

	"udp_iaas/handlers"

	"libvirt.org/go/libvirt"
)

func main() {
	// VM 디스크 저장 디렉토리 확인 및 생성
	vmDiskPath := "/var/lib/vms/disks"
	vmIsoPath := "/var/lib/vms/isos"

	// 디스크 디렉토리 생성
	if _, err := os.Stat(vmDiskPath); os.IsNotExist(err) {
		if err := os.MkdirAll(vmDiskPath, 0755); err != nil {
			log.Fatalf("Failed to create VM disk directory: %v", err)
		}
	}

	// ISO 디렉토리 생성
	if _, err := os.Stat(vmIsoPath); os.IsNotExist(err) {
		if err := os.MkdirAll(vmIsoPath, 0755); err != nil {
			log.Fatalf("Failed to create VM ISO directory: %v", err)
		}
	}

	// libvirt 연결 설정
	conn, err := libvirt.NewConnect("qemu:///system")
	if err != nil {
		log.Fatalf("Failed to connect to libvirt: %v", err)
	}
	defer conn.Close()

	app := &handlers.App{
		Libvirt: conn,
	}

	// 정적 파일 서비스 설정
	http.Handle("/", http.FileServer(http.Dir("static")))

	http.HandleFunc("/api/vms/", app.HandleVMs)
	http.HandleFunc("/api/vms", app.HandleVMs)
	http.HandleFunc("/api/health", app.HandleHealthCheck)
	http.HandleFunc("/api/isos", app.HandleISOs)

	// 서버 시작
	log.Printf("Server starting on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
