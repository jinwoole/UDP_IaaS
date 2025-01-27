package main

import (
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"udp_iaas/handlers"
	"udp_iaas/types"

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
   startCmd := exec.Command("rc-service", "libvirtd", "start")
   if err := startCmd.Run(); err != nil {
       return err
   }

   updateCmd := exec.Command("rc-update", "add", "libvirtd")
   return updateCmd.Run()
}

func main() {
   vmDiskPath := "/var/lib/vms/disks"
   vmIsoPath := "/var/lib/vms/isos"

   if _, err := os.Stat(vmDiskPath); os.IsNotExist(err) {
       if err := os.MkdirAll(vmDiskPath, 0755); err != nil {
           log.Fatalf("Failed to create VM disk directory: %v", err)
       }
   }

   if _, err := os.Stat(vmIsoPath); os.IsNotExist(err) {
       if err := os.MkdirAll(vmIsoPath, 0755); err != nil {
           log.Fatalf("Failed to create VM ISO directory: %v", err)
       }
   }

   if !isLibvirtdRunning() {
       log.Println("libvirtd service is not running. Attempting to start...")
       if err := startLibvirtd(); err != nil {
           log.Fatalf("Failed to start libvirtd service: %v", err)
       }
       log.Println("libvirtd service started successfully")
   }

   conn, err := libvirt.NewConnect("qemu:///system")
   if err != nil {
       log.Fatalf("Failed to connect to libvirt: %v", err)
   }
   defer conn.Close()

   websockifyMgr := types.NewWebsockifyManager()

   app := &handlers.App{
       Libvirt:    conn,
       Websockify: websockifyMgr,
   }

   http.Handle("/", http.FileServer(http.Dir("static")))

   http.HandleFunc("/api/vms/", app.HandleVMs)
   http.HandleFunc("/api/vms", app.HandleVMs)
   http.HandleFunc("/api/health", app.HandleHealthCheck)
   http.HandleFunc("/api/isos", app.HandleISOs)

   serverAddr := ":8080"
   log.Printf("Server starting on %s", serverAddr)
   if err := http.ListenAndServe(serverAddr, nil); err != nil {
       log.Fatalf("Server failed: %v", err)
   }
}