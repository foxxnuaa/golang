package main

import (
	"log"
	"performance/test"
	"runtime"
)

var printDetail = false

func main() {
	// 设置P最大数量
	runtime.GOMAXPROCS(runtime.NumCPU())

	serverAddr := "127.0.0.1:9797"

	// 初始化服务器
	server := test.NewTcpServer()
	defer server.Close()
	log.Printf("Startup TCP server(%s)...\n", serverAddr)
	err := server.Listen(serverAddr)
	if err != nil {
		log.Printf("TCP Server startup failing! (addr=%s)!\n", serverAddr)
	}

	ch := make(chan byte)
	<-ch
}
