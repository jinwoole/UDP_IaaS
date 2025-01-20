package router

import (
	systemService "udp_iaas/service/system"
	vmservice "udp_iaas/service/vm"
	systemhandler "udp_iaas/web/handler/system"
	vmhandler "udp_iaas/web/handler/vm"

	"github.com/labstack/echo/v4"
)

type Router struct {
	e             *echo.Echo
	vmService     *vmservice.Service
	systemService *systemService.Service
}

func NewRouter(vmService *vmservice.Service, systemService *systemService.Service) *Router {
	return &Router{
		e:             echo.New(),
		vmService:     vmService,
		systemService: systemService,
	}
}

func (r *Router) SetupRoutes() *echo.Echo {
	// 정적 파일 라우팅 설정
	r.e.Static("/", "console") // 루트는 콘솔

	// VM 관련 핸들러 초기화
	vmHandler := vmhandler.NewHandler(r.vmService)

	// VM 라우트 그룹 설정
	vm := r.e.Group("/vm")
	vm.GET("", vmHandler.List)
	vm.POST("/create", vmHandler.Create)

	// System 라우트 그룹 설정
	systemHandler := systemhandler.NewHandler(r.systemService)
	system := r.e.Group("/system")
	system.GET("/info", systemHandler.GetInfo)
	system.GET("/info/iso", systemHandler.GetISOInfo)

	return r.e
}
