package router

import (
	"udp_iaas/src/api/handler"
	"udp_iaas/src/libvirt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func NewRouter(client *libvirt.Client) *echo.Echo {
	e := echo.New()

	// 미들웨어 설정
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// 핸들러 초기화
	statusHandler := handler.NewStatusHandler(client)
	vmHandler := handler.NewVMHandler(client)

	// API 그룹 생성
	api := e.Group("/api/v1")

	// 라우트 설정
	api.GET("/status", statusHandler.GetStatus)
	api.GET("/vms", vmHandler.ListVMs)

	e.Static(("/"), "web/dist")
	e.File("/", "web/dist/index.html")

	return e
}
