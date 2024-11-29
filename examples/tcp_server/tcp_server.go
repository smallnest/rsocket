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
	host = flag.String("h", "0.0.0.0", "listen server host")
	port = flag.Int("p", 8000, "listen server port")
)

func main() {
	flag.Parse()

	// 创建RDMA socket
	fd, err := rsocket.Socket(rsocket.AF_INET, rsocket.SOCK_STREAM, 0)
	if err != nil {
		log.Fatal(err)
	}
	defer rsocket.Close(fd)

	// 绑定到地址
	addr := &net.TCPAddr{IP: net.ParseIP(*host), Port: *port}
	sa := &syscall.SockaddrInet4{
		Port: addr.Port,
	}
	copy(sa.Addr[:], addr.IP.To4())

	if err := rsocket.Bind(fd, sa); err != nil {
		log.Fatal(err)
	}

	// 监听连接
	if err := rsocket.Listen(fd, 128); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("服务器正在监听 :%d\n", addr.Port)

	for {
		// 接受新的连接
		clientFd, clientAddr, err := rsocket.Accept(fd)
		if err != nil {
			log.Printf("接受连接失败: %v\n", err)
			continue
		}

		// 处理新的连接
		go handleConnection(clientFd, clientAddr)
	}
}

func handleConnection(clientFd int, clientAddr syscall.Sockaddr) {
	defer rsocket.Close(clientFd)

	// 将clientAddr转换为更易读的格式
	addr, ok := clientAddr.(*syscall.SockaddrInet4)
	if ok {
		ip := net.IPv4(addr.Addr[0], addr.Addr[1], addr.Addr[2], addr.Addr[3])
		fmt.Printf("新的客户端连接: %s:%d\n", ip.String(), addr.Port)
	}

	// 读取客户端数据
	buffer := make([]byte, 1024)
	n, err := rsocket.Read(clientFd, buffer)
	if err != nil {
		log.Printf("读取数据失败: %v\n", err)
		return
	}
	fmt.Printf("收到客户端消息: %s\n", buffer[:n])

	// 发送响应
	response := []byte("Server received your message!")
	n, err = rsocket.Write(clientFd, response)
	if err != nil {
		log.Printf("发送响应失败: %v\n", err)
		return
	}
	fmt.Printf("发送了 %d 字节的响应\n", n)
}
