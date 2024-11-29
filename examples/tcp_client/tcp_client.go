package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"syscall"

	"github.com/smallnest/rsocket"
)

var (
	host = flag.String("h", "127.0.0.1", "server host")
	port = flag.Int("p", 8000, "server port")
)

func main() {
	flag.Parse()

	// 创建RDMA socket
	fd, err := rsocket.Socket(rsocket.AF_INET, rsocket.SOCK_STREAM, 0)
	if err != nil {
		log.Fatal(err)
	}
	defer rsocket.Close(fd)
	// Create RDMA socket
	// 准备服务器地址
	serverAddr := &net.TCPAddr{IP: net.ParseIP(*host), Port: *port}
	sa := &syscall.SockaddrInet4{
		Port: serverAddr.Port,
	}
	copy(sa.Addr[:], serverAddr.IP.To4())

	// 连接到服务器
	if err := rsocket.Connect(fd, sa); err != nil {
		log.Fatal("连接失败:", err)
	}
	fmt.Println("成功连接到服务器")

	// 发送数据
	message := []byte("Hello, RDMA Server!")
	n, err := rsocket.Write(fd, message)
	if err != nil {
		log.Fatal("发送数据失败:", err)
	}
	fmt.Printf("发送了 %d 字节的数据\n", n)

	// 接收响应
	buffer := make([]byte, 1024)
	n, err = rsocket.Read(fd, buffer)
	if err != nil {
		log.Fatal("接收数据失败:", err)
	}
	fmt.Printf("收到服务器响应: %s\n", buffer[:n])
}
