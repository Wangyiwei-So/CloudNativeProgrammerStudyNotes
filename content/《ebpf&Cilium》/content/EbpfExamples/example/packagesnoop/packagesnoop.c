#include <vmlinux.h>
#include <bpf/bpf_helpers.h>
#include <bpf/bpf_tracing.h>
// #include <linux/sched.h>
char LICENSE[] SEC("license") = "Dual BSD/GPL";

const char* get_protocol_name(__u8 protocol){
    switch (protocol){
        case 1: //头文件#include <netinet/in.h>
            return "ICMP";
        case 6:
            return "TCP";
        case 17:
            return "UDP";
        case 2:
            return "IGMP";
        case 132:
            return "SCTP";
        default:
            return "Unknown";
    }
}

#define TASK_COMM_LEN 16

SEC("xdp")
int xdp_program(struct xdp_md* ctx) {
    void *data = (void *)(long)ctx->data;
    void *data_end = (void *)(long)ctx->data_end;
    int pkt_sz = data_end - data;

    //链路层
    struct ethhdr *eth = data; //ethhdr是链路层(Ethernet)头结构体
    //ethhdr在#include <linux/if_ether.h>头文件
    if ((void*)eth + sizeof(*eth) > data_end){ //eth的首地址+ethhdr的结构体大小得到尾地址，不会比数据包的尾地址大的
        bpf_printk("Invalid ethernet header\n");
        return XDP_PASS;
    }

    //ip层
    struct iphdr *ip = data+sizeof(*eth); //iphdr是传输层ip头
    //在<linux/ip.h>头文件
    if((void*)ip + sizeof(*ip) > data_end) {
        bpf_printk("Invalid ip header\n");
        return XDP_PASS;
    }
    // char comm[TASK_COMM_LEN];
    // bpf_get_current_comm(&comm, sizeof(comm));
    // u64 tidpid = bpf_get_current_pid_tgid();
    bpf_printk("sent package size is %d, protocol is %s" , pkt_sz, get_protocol_name(ip->protocol));
    return XDP_PASS;
}