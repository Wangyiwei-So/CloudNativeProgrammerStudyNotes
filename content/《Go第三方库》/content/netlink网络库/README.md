# netlink网络库

`netlink`包是用于与 Linux 系统进行网络操作的库。它提供了一种与 Linux 内核通信的方式，可以进行一些网络配置，如路由表、网络设备等的管理。

# 1. 示例

## 1.1 获取网络设备

```go
// 获取所有设备，相当于运行 `ip link show`
func PrintLinks() error {
	lls, err := netlink.LinkList()
	if err != nil {
		return err
	}
	for _, l := range lls {
		fmt.Printf("Type: %s, Name: %s, Flags: %s, OperState: %s\n", l.Type(), l.Attrs().Name, l.Attrs().Flags.String(), l.Attrs().OperState.String())
	}
	return nil
}
```

这里有一个比较重要的结构体`LinkAttrs`，代表了`ip link show`得到的条目

```go
// LinkAttrs represents data shared by most link types
type LinkAttrs struct {
	Index        int //网络接口的索引号
	MTU          int //最大传输单元，表示网络接口可以传输的最大数据包大小
	TxQLen       int //传输队列长度，表示网络接口的传输队列长度
	Name         string //网络接口名称
	HardwareAddr net.HardwareAddr //MAC地址
	Flags        net.Flags //网络接口的状态UP、DOWN、BROADCAST等, 最终是类似这样的up|broadcast|multicast
	RawFlags     uint32 //原始形式表示的网络接口状态
	ParentIndex  int         // 父网络接口索引号
	MasterIndex  int         // 如果该对象是网桥的成员，则网桥的索引号
	Namespace    interface{} // 网络命名空间的标识可以是nil | NsPid | NsFd
	Alias        string //别名
	Statistics   *LinkStatistics //有关网络接口统计信息
	Promisc      int //表示网络接口是否使用混杂模式
	Allmulti     int //Allmulti模式标志
	Multi        int //Multi模式标志
	Xdp          *LinkXdp //包含与eBPF相关的信息的结构体指针
	EncapType    string //封装类型
	Protinfo     *Protinfo //有关网络接口协议信息的结构体
	OperState    LinkOperState //网络接口的操作状态
	PhysSwitchID int //物理交换机标识
	NetNsID      int //网络命名空间的ID
	NumTxQueues  int //发送队列的数量
	NumRxQueues  int //接收队列的数量
	GSOMaxSize   uint32 //GSO（Generic Segmentation Offload）最大大小
	GSOMaxSegs   uint32 //GSO最大分段数
	Vfs          []VfInfo // virtual functions available on link
	Group        uint32 //网络接口所属的组
	Slave        LinkSlave //标识网络接口是否是另一个接口的从属
}
```

以下是我没看懂的字段

- `MasterIndex`：通常，只有作为`网桥的成员`的接口才会有主网桥。当一个网络接口被添加到网桥中时，它的`MasterIndex`字段将指示该接口所属的网桥的索引号。如果网络接口不是任何网桥的成员，那么`MasterIndex`将为0
- `Promisc`：混杂模式，通常网络接口只会接受目标地址是自己的数据帧。在混杂模式下，网络接口将接收通过该网络传输的所有数据帧，通常用于网络监测、分析和数据包捕获等应用场景，例如Wireshark
- `Allmulti`: 是网络接口的一种工作模式，通常网络接口只会接收目标地址是自己、或者广播地址的数据帧。如果启用`Allmulti`模式后，网络接口将接收到经过它所在网络的**所有**多播组的数据帧
- `Multi`: 网络接口仅接收目标地址为**特定**多播组地址的数据帧
- `EncapType`: 不同类型的网络接口可能支持不同的封装类型。例如在虚拟化环境中，网络接口可能支持不同的封装协议，VLAN（虚拟局域网）VXLAN（虚拟扩展局域网）Geneve（通用网络封装）GRE（通用路由封装）
- `OperState` : UP和DOWN
- `GSO`: 是一种网络协议栈的功能，用于提高大型数据包的传输效率。它的主要目标是将大型数据包分割成较小的片段，以便更容易在网络上传输。这种技术通常用于提高网络性能和降低传输延迟。

```go
// 按名字获取，相当于运行`ip link show xxx`
func PrintLink(name string) error {
	l, err := netlink.LinkByName(name)
	if err != nil {
		return err
	}
	fmt.Printf("Type: %s, Name: %s, Flags: %s, OperState: %s\n", l.Type(), l.Attrs().Name, l.Attrs().Flags.String(), l.Attrs().OperState.String())
	return nil
}
```

## 1.2 创建网桥

```go
// 创建或更新网桥
func CreateOrUpdateBridge(br string, br_addr string) (*netlink.Bridge, error) {
	link, err := netlink.LinkByName(br)
	if err != nil {
		if _, ok := err.(netlink.LinkNotFoundError); ok {
			//初始化各网桥对象
			br := &netlink.Bridge{
				LinkAttrs: netlink.LinkAttrs{
					Name: br,
					MTU:  1500,
				},
			}
			// 解析地址的方法
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
```

## 1.3 创建veth设备

```go
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
```

## 1.4 操作路由表

```go
// 相当于执行`ip route add $route`
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
```

todo
