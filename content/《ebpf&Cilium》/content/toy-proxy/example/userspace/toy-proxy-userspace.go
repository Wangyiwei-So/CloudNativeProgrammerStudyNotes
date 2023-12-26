package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os/exec"
	"strconv"
	"strings"
	"sync"
)

func main() {
	clusterIP := "10.7.111.132"
	podIP := "10.5.41.204"
	port := 80
	proto := "tcp"
	addRedirectRules(clusterIP, port, proto)
	createProxy(podIP, port, proto)
}

func addRedirectRules(clusterIP string, port int, proto string) error {
	p := strconv.Itoa(port)
	cmd := exec.Command("iptables", "-t", "nat", "-A", "OUTPUT", "-p", "tcp",
		"-d", clusterIP, "--dport", p, "-j", "REDIRECT", "--to-port", p)
	/*针对iptables命令的解释
	-t nat: 指定操作的表为NAT表
	-A OUTPUT: 将规则添加到OUTPUT链
	-p tcp: 指定匹配的协议为TCP
	-d clusterIP: 指定目标地址为10.7.111.132
	--dport p: 指定目标端口为80
	-j REDIRECT: 如果规则匹配，将流量重定向
	--to-port p: 重定向到本地端口80
	*/
	return cmd.Run()
}

func createProxy(podIP string, port int, proto string) {
	host := ""
	listener, err := net.Listen(proto, net.JoinHostPort(host, strconv.Itoa(port)))
	if err != nil {
		log.Fatalln("listen error: ", err)
	}
	for {
		inConn, err := listener.Accept()
		if err != nil {
			log.Fatalln("connect error: ", err)
		}
		outConn, err := net.Dial(proto, net.JoinHostPort(podIP, strconv.Itoa(port)))
		if err != nil {
			log.Fatalln("dial error: ", err)
		}
		go func(in, out *net.TCPConn) {
			var wg sync.WaitGroup
			wg.Add(2)
			fmt.Printf("Proxying %v <-> %v <-> %v <-> %v\n",
				in.RemoteAddr(), in.LocalAddr(), out.LocalAddr(), out.RemoteAddr())
			go copyBytes(in, out, &wg)
			go copyBytes(out, in, &wg)
			wg.Wait()
		}(inConn.(*net.TCPConn), outConn.(*net.TCPConn))
	}
}

func copyBytes(dst, src *net.TCPConn, wg *sync.WaitGroup) {
	defer wg.Done()
	if _, err := io.Copy(dst, src); err != nil {
		if !strings.HasSuffix(err.Error(), "use of closed network connection") {
			fmt.Printf("io.Copy error: %v", err)
		}
	}
	dst.Close()
	src.Close()
}
