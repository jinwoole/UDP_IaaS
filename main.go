package main

import (
	"log"

	"libvirt.org/go/libvirt"

	systemData "udp_iaas/data/system"
	vmData "udp_iaas/data/vm"
	systemService "udp_iaas/service/system"
	vmService "udp_iaas/service/vm"

	"udp_iaas/web/router"

	_ "udp_iaas/docs"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title UDP IaaS API
// @version 1.0
// @description API documentation for UDP IaaS
// @host localhost:8116
// @BasePath /

func main() {
	// libvirt 연결
	conn, err := libvirt.NewConnect("qemu:///system")
	if err != nil {
		log.Fatal("Failed to connect to libvirt:", err)
	}
	defer conn.Close()

	// 각 계층 초기화
	vmData := vmData.NewData(conn)
	vmService := vmService.NewService(vmData)
	systemData := systemData.NewData(conn)
	systemService := systemService.NewService(systemData)

	// 라우터 설정
	r := router.NewRouter(vmService, systemService)
	e := r.SetupRoutes()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Enable CORS
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:8116", "http://192.168.50.96:8116"},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE},
	}))
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// 서버 시작
	log.Println("Server starting on :8116...")
	if err := e.Start(":8116"); err != nil {
		log.Fatal(err)
	}
}
