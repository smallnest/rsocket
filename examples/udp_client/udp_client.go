package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"strconv"
	"syscall"

	"github.com/smallnest/rsocket"
)

var (
	clientAddr = flag.String("c", "127.0.0.1:7000", "client address")
	serverAddr = flag.String("s", "127.0.0.1:8000", "server address")
)

func main() {
	flag.Parse()

	// 创建RDMA UDP socket
	fd, err := rsocket.Socket(rsocket.AF_INET, rsocket.SOCK_DGRAM, 0)
	if err != nil {
		log.Fatal(err)
	}
	defer rsocket.Close(fd)

	// // 连接到服务器（UDP中这步是可选的，但可以简化后续的通信）
	// if err := rsocket.Connect(fd, sa); err != nil {
	// 	log.Fatal("连接失败:", err)
	// }
	// fmt.Println("UDP socket已就绪")
	localhost, localport, _ := net.SplitHostPort(*clientAddr)
	localUDPPort, _ := strconv.Atoi(localport)

	cAddr := &net.UDPAddr{IP: net.ParseIP(localhost), Port: localUDPPort}
	ca := &syscall.SockaddrInet4{
		Port: cAddr.Port,
	}
	copy(ca.Addr[:], cAddr.IP.To4())

	if err := rsocket.Bind(fd, ca); err != nil {
		log.Fatalf("bind err: %v", err)
	}

	fmt.Printf("UDP 服务器正在监听 :%s\n", *clientAddr)

	// 准备服务器地址
	host, port, _ := net.SplitHostPort(*serverAddr)
	udpPort, _ := strconv.Atoi(port)
	serverAddr := &net.UDPAddr{IP: net.ParseIP(host), Port: udpPort}
	sa := &syscall.SockaddrInet4{
		Port: serverAddr.Port,
	}
	copy(sa.Addr[:], serverAddr.IP.To4())

	// 发送数据
	message := []byte("Hello, RDMA UDP Server!")
	n, err := rsocket.SendTo(fd, message, 0, sa)
	if err != nil {
		log.Fatal("发送数据失败:", err)
	}
	fmt.Printf("发送了 %d 字节的数据\n", n)

	// 接收响应
	buffer := make([]byte, 65507) // UDP最大包大小
	n, err = rsocket.Read(fd, buffer)
	if err != nil {
		log.Fatal("接收数据失败:", err)
	}
	fmt.Printf("收到服务器响应: %s\n", buffer[:n])
}
