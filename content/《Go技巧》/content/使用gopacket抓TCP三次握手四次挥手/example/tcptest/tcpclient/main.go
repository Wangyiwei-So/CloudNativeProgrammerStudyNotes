package main

import (
	"fmt"
	"log"
	"net"
	"tcptest/util"
)

func main() {
	conn, err := net.Dial("tcp", fmt.Sprintf("localhost:%d", util.TCP_SERVER_PORT))
	if err != nil {
		log.Fatal(err)
	}
	conn.Write([]byte("hello world"))
	defer conn.Close()
}
