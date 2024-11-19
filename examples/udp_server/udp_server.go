package main

import (
	"fmt"
	"log"
	"net"
	"syscall"

	"github.com/smallnest/rsocket"
)

func main() {
	// 创建RDMA UDP socket
	fd, err := rsocket.Socket(rsocket.AF_INET, rsocket.SOCK_DGRAM, 0)
	if err != nil {
		log.Fatalf("socket err: %v", err)
	}
	defer rsocket.Close(fd)

	// 绑定到地址
	addr := &net.UDPAddr{IP: net.ParseIP("0.0.0.0"), Port: 8001}
	sa := &syscall.SockaddrInet4{
		Port: addr.Port,
	}
	copy(sa.Addr[:], addr.IP.To4())

	if err := rsocket.Bind(fd, sa); err != nil {
		log.Fatalf("bind err: %v", err)
	}

	fmt.Printf("UDP 服务器正在监听 :%d\n", addr.Port)

	// 循环接收数据
	buffer := make([]byte, 65507) // UDP最大包大小
	for {
		n, err := rsocket.Read(fd, buffer)
		if err != nil {
			log.Printf("读取数据失败: %v\n", err)
			continue
		}

		message := string(buffer[:n])
		fmt.Printf("收到消息: %s\n", message)

		// 发送响应
		response := []byte("Server received your UDP message!")
		n, err = rsocket.Write(fd, response)
		if err != nil {
			log.Printf("发送响应失败: %v\n", err)
			continue
		}
		fmt.Printf("发送了 %d 字节的响应\n", n)
	}
}
