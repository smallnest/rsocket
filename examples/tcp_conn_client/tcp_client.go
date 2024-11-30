package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/smallnest/rsocket"
)

var (
	serverAddr = flag.String("s", "127.0.0.1:8000", "server address")
)

func main() {
	flag.Parse()

	conn, err := rsocket.DialTCP(*serverAddr)
	if err != nil {
		log.Fatal(err)
	}

	// 发送数据
	message := []byte("Hello, RDMA Server!")
	_, err = conn.Write(message)
	if err != nil {
		log.Fatal("发送数据失败:", err)
	}

	// 接收响应
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		log.Fatal("接收数据失败:", err)
	}
	fmt.Printf("收到服务器响应: %s\n", buffer[:n])
}
