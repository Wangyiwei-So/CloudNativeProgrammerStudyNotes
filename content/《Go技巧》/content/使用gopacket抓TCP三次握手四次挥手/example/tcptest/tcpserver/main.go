package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"tcptest/util"
)

func handleConn(conn net.Conn) {
	defer conn.Close()
	for {
		reader := bufio.NewReader(conn)
		var buf [128]byte
		n, err := reader.Read(buf[:])
		if err != nil {
			if err == io.EOF {
				log.Println("断开连接")
				break
			}
			log.Println("数据读取错误", err)
			break
		}
		recvStr := string(buf[:n])
		fmt.Println("收到数据", recvStr)
		// conn.Write([]byte("pang"))
	}
}

func main() {
	listener, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", util.TCP_SERVER_PORT))
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("启动tcpserver,监听", util.TCP_SERVER_PORT)
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("建立连接错误", err)
			continue
		}
		go handleConn(conn)
	}
}
