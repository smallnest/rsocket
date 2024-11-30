package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"strconv"

	"github.com/smallnest/rsocket"
)

var (
	serverAddr = flag.String("s", "0.0.0.0:8000", "server address")
)

func main() {
	flag.Parse()

	host, port, err := net.SplitHostPort(*serverAddr)
	tcpPort, _ := strconv.Atoi(port)

	// 创建RDMA socket
	ln, err := rsocket.NewTCPListener(host, tcpPort)
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()

	fmt.Printf("服务器正在监听 :%d\n", *serverAddr)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("接受连接失败: %v\n", err)
			return
		}

		// 处理新的连接
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	// 读取客户端数据
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		log.Printf("读取数据失败: %v\n", err)
		return
	}
	fmt.Printf("收到客户端消息: %s\n", buffer[:n])

	// 发送响应
	response := []byte("Server received your message!")
	n, err = conn.Write(response)
	if err != nil {
		log.Printf("发送响应失败: %v\n", err)
		return
	}
	fmt.Printf("发送了 %d 字节的响应\n", n)
}
