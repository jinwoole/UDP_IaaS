package main

import (
	"log"
	"udp_iaas/src/api/router"
	"udp_iaas/src/libvirt"
	"udp_iaas/src/storage/badger"
)

func main() {
	// BadgerDB 클라이언트 초기화
	dbClient, err := badger.NewClient("./data/db")
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer dbClient.Close()

	// libvirt 클라이언트 초기화
	libvirtClient, err := libvirt.NewClient()
	if err != nil {
		log.Fatal("Failed to initialize libvirt:", err)
	}
	defer libvirtClient.Close()

	// 라우터 설정
	e := router.NewRouter(libvirtClient)

	// 서버 시작
	e.Logger.Fatal(e.Start(":8080"))
}
