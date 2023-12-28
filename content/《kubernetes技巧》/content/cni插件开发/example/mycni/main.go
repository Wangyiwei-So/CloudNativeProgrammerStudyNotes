package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"net"
	"os"

	"github.com/containernetworking/cni/pkg/skel"
	"github.com/containernetworking/cni/pkg/types"
	types4 "github.com/containernetworking/cni/pkg/types/040"
	"github.com/containernetworking/cni/pkg/version"
	"github.com/containernetworking/plugins/pkg/ipam"
	"github.com/vishvananda/netlink"
	"github.com/vishvananda/netns"
)

const (
	BR_ADDR = "10.16.0.1/16"
)

func log(format string, a ...any) { //由于不能往stdout中打印，所以搞了个往stderr中打印的函数，用fmt.Errorf一样
	fmt.Fprintf(os.Stderr, format, a...)
}

func addHandler(args *skel.CmdArgs) error { //cni add的处理函数
	cfg, err := ConfigFromStdin(args.StdinData) //从stdin中获取配置, 可选
	if err != nil {
		return err
	}

	ret := types4.Result{ //注意看引用的类型040就是cniversion是0.4.0的，直接使用库给的返回result的方法
		CNIVersion: cfg.CNIVersion,
	}

	if cfg.IPAM.Type != "" {
		r, err := ipam.ExecAdd(cfg.IPAM.Type, args.StdinData) //执行host-local
		if err != nil {
			return err
		}
		ipamRet, err := types4.NewResultFromResult(r) //解析result
		if err != nil {
			return err
		}
		ret.IPs = ipamRet.IPs //传给上下文
		ret.DNS = ipamRet.DNS
	}

	br, err := CreateOrUpdateBridge(cfg.Bridge, BR_ADDR) //这里地址不要和docker0重了
	if err != nil {
		return err
	}
	log("得到了hostlocal分配的ip: %s, ns是%s", ret.IPs[0].Address.String(), args.Netns)
	err = CreateVeth(args.Netns, ret.IPs[0].Address.String(), br)
	if err != nil {
		return err
	}

	return ret.Print()
}

type Config struct {
	types.NetConf
	Bridge string `json:"bridge"`
}

func ConfigFromStdin(data []byte) (*Config, error) {
	cfg := &Config{}
	err := json.Unmarshal(data, cfg)
	if err != nil {
		return nil, err
	} else {
		return cfg, nil
	}
}

// 创建或更新网桥
func CreateOrUpdateBridge(br string, br_addr string) (*netlink.Bridge, error) {
	link, err := netlink.LinkByName(br)
	if err != nil {
		if _, ok := err.(netlink.LinkNotFoundError); ok {
			br := &netlink.Bridge{
				LinkAttrs: netlink.LinkAttrs{
					Name: br,
					MTU:  1500,
				},
			}
			if err := netlink.LinkAdd(br); err != nil {
				return nil, err
			}
			var addr *netlink.Addr
			if addr, err = netlink.ParseAddr(br_addr); err != nil {
				return nil, err
			}
			if err = netlink.AddrAdd(br, addr); err != nil { //设置ip
				return nil, err
			}
			if err = netlink.LinkSetUp(br); err != nil { //把网桥设为UP
				return nil, err
			}
			return br, nil
		} else {
			return nil, err
		}
	}
	if br, ok := link.(*netlink.Bridge); ok {
		return br, nil
	}
	return nil, fmt.Errorf("错误的网桥对象")
}

// 创建veth设备
func CreateVeth(nspath string, addrstr string, br *netlink.Bridge) error {
	var veth_host, veth_container = RandomVethName(), RandomVethName()
	vethpeer := &netlink.Veth{
		LinkAttrs: netlink.LinkAttrs{
			Name: veth_host,
			MTU:  1500,
		},
		PeerName: veth_container,
	}
	err := netlink.LinkAdd(vethpeer)
	if err != nil {
		return err
	}

	veth_host_interface, err := netlink.LinkByName(veth_host)
	if err != nil {
		return err
	}
	err = netlink.LinkSetMaster(veth_host_interface, br)
	if err != nil {
		return err
	}
	err = netlink.LinkSetUp(veth_host_interface)
	if err != nil {
		return err
	}

	ns, err := netns.GetFromPath(nspath)
	if err != nil {
		return err
	}
	defer ns.Close()
	veth_container_interface, err := netlink.LinkByName(veth_container)
	if err != nil {
		return err
	}
	err = netlink.LinkSetNsFd(veth_container_interface, int(ns))
	if err != nil {
		return err
	}
	err = netns.Set(ns) //进入这个ns进行操作
	if err != nil {
		return err
	}
	veth_container_interface, err = netlink.LinkByName(veth_container)
	if err != nil {
		return err
	}
	addr, _ := netlink.ParseAddr(addrstr)
	err = netlink.AddrAdd(veth_container_interface, addr)
	if err != nil {
		return err
	}
	err = netlink.LinkSetName(veth_container_interface, "eth0") //把容器内的veth名字设为eth0
	if err != nil {
		return err
	}
	err = netlink.LinkSetUp(veth_container_interface)
	if err != nil {
		return err
	}
	return AddRoute()
}

func RandomVethName() string {
	entropy := make([]byte, 4)
	rand.Read(entropy)
	return fmt.Sprintf("wywveth%x", entropy)
}

func AddRoute() error {
	route := &netlink.Route{
		Dst: &net.IPNet{
			IP:   net.IPv4(0, 0, 0, 0),
			Mask: net.IPv4Mask(0, 0, 0, 0),
		},
		Gw: net.IPv4(10, 16, 0, 1),
	}
	return netlink.RouteAdd(route)
}

func main() {
	skel.PluginMain(addHandler, nil, nil, version.All, "")
}
