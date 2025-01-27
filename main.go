package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"udp_iaas/handlers"

	"libvirt.org/go/libvirt"
)

func isLibvirtdRunning() bool {
	cmd := exec.Command("rc-service", "libvirtd", "status")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return false
	}
	return strings.Contains(string(output), "started")
}

func startLibvirtd() error {
	// Start the service
	startCmd := exec.Command("rc-service", "libvirtd", "start")
	if err := startCmd.Run(); err != nil {
		return err
	}

	// Add to default runlevel
	updateCmd := exec.Command("rc-update", "add", "libvirtd")
	return updateCmd.Run()
}

func startWebsockify(port int) error {
    cmd := exec.Command("websockify",
        "--web", "/usr/share/novnc/",
        fmt.Sprintf("%d", port),           // 웹소켓 포트 (5901)
        fmt.Sprintf("localhost:%d", 5900), // VNC 포트
    )
    return cmd.Start()
}

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

	// libvirtd 서비스 상태 확인 및 시작
	if !isLibvirtdRunning() {
		log.Println("libvirtd service is not running. Attempting to start...")
		if err := startLibvirtd(); err != nil {
			log.Fatalf("Failed to start libvirtd service: %v", err)
		}
		log.Println("libvirtd service started successfully")
	}

	// libvirt 연결 설정
	conn, err := libvirt.NewConnect("qemu:///system")
	if err != nil {
		log.Fatalf("Failed to connect to libvirt: %v", err)
	}
	defer conn.Close()

    app := &handlers.App{
        Libvirt: conn,
        StartWebsockify: startWebsockify,  // 함수 전달
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
