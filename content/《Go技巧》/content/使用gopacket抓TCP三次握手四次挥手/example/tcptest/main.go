package main

import (
	"fmt"
	"log"
	"tcptest/util"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

const (
	NET_INTERFACE = "lo"
)

func main() {
	// 打开要监控的网络设备，这里使用了环回设备lo，得到handle
	handle, err := pcap.OpenLive(NET_INTERFACE, 2048, false, time.Second*5)
	if err != nil {
		log.Fatalln(err)
	}
	defer handle.Close()
	// 从handle创建包的数据源
	source := gopacket.NewPacketSource(handle, handle.LinkType())
	log.Printf("开始监听%s网络接口的tcp报文", NET_INTERFACE)
	for pkg := range source.Packets() { //从channel中取出包
		if layer4 := pkg.TransportLayer(); layer4 != nil { //判断包的网络层协议有没有，比如ARP协议就没有网络层
			if tcp, ok := layer4.(*layers.TCP); ok { //断言判断包的网络层协议是不是TCP，比如UDP就不是
				if layer3 := pkg.NetworkLayer(); layer3 != nil {
					if ip, ok := layer3.(*layers.IPv4); ok {
						if tcp.DstPort == util.TCP_SERVER_PORT || tcp.SrcPort == util.TCP_SERVER_PORT {
							fmt.Printf("tcp报文 | %s:%d-->%s:%d | SYN: %v, ACK: %v, FIN: %v | 序列号: %d, 确认应答号: %d | payload: %s\n",
								ip.SrcIP.String(),
								tcp.SrcPort,
								ip.DstIP.String(),
								tcp.DstPort,
								tcp.SYN,
								tcp.ACK,
								tcp.FIN,
								tcp.Seq,
								tcp.Ack,
								string(tcp.Payload),
							)
						}
					}
				}

			}
		}
	}
}
