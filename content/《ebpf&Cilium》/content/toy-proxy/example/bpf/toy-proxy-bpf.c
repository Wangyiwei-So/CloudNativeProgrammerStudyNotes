//go: build ignore

#include <linux/bpf.h>     // struct __sk_buff
#include <linux/pkt_cls.h> // TC_ACT_OK
#include <linux/ip.h>      // struct iphdr
#include <linux/tcp.h>     // struct tcphdr
#include <stdint.h>        // uint32_t
#include <stddef.h>        // offsetof()
#include <bpf/bpf_helpers.h>
#include <bpf/bpf_tracing.h>

char LICENSE[] SEC("license") = "Dual BSD/GPL";

#define ETH_HLEN 14
#define IPPROTO_TCP 6

SEC("egress")
int tc_egress(struct __sk_buff *skb)
{
    const __be32 cluster_ip = 0x846F070A;
    const __be32 pod_ip = 0x0529050A;

    const int l3_off = ETH_HLEN;
    const int l4_off = l3_off + 20;
    __be32 sum;
    void *data = (void *)(long)skb->data;
    void *data_end = (void *)(long)skb->data_end;
    if (data_end < data + l4_off) { // not our packet
        return TC_ACT_OK;
    }

    struct iphdr *ip4 = (struct iphdr *)(data + l3_off);
    if (ip4->daddr != cluster_ip || ip4->protocol != IPPROTO_TCP /* || tcp->dport == 80 */) { //只拦截目标地址是cluster_ip的TCP流量
        return TC_ACT_OK;
    }

    // DNAT: cluster_ip -> pod_ip, then update L3 and L4 checksum
    sum = bpf_csum_diff((void *)&ip4->daddr, 4, (void *)&pod_ip, 4, 0);
    bpf_skb_store_bytes(skb, l3_off + offsetof(struct iphdr, daddr), (void *)&pod_ip, 4, 0); //把ClusterIP改为PodIP
    bpf_l3_csum_replace(skb, l3_off + offsetof(struct iphdr, check), 0, sum, 0);
	bpf_l4_csum_replace(skb, l4_off + offsetof(struct tcphdr, check), 0, sum, BPF_F_PSEUDO_HDR);
    return TC_ACT_OK;  
}

SEC("ingress")
int tc_ingress(struct __sk_buff *skb){
    const __be32 cluster_ip = 0x846F070A; // 10.7.111.132
    const __be32 pod_ip = 0x0529050A;     // 10.5.41.5

    const int l3_off = ETH_HLEN;    // IP header offset
    const int l4_off = l3_off + 20; // TCP header offset: l3_off + sizeof(struct iphdr)
    __be32 sum;                     // IP checksum

    void *data = (void *)(long)skb->data;
    void *data_end = (void *)(long)skb->data_end;
    if (data_end < data + l4_off) { // not our packet
        return TC_ACT_OK;
    }

    struct iphdr *ip4 = (struct iphdr *)(data + l3_off);
    if (ip4->saddr != pod_ip || ip4->protocol != IPPROTO_TCP /* || tcp->dport == 80 */) {
        return TC_ACT_OK;
    }

    // SNAT: pod_ip -> cluster_ip, then update L3 and L4 header
    sum = bpf_csum_diff((void *)&ip4->saddr, 4, (void *)&cluster_ip, 4, 0);
    bpf_skb_store_bytes(skb, l3_off + offsetof(struct iphdr, saddr), (void *)&cluster_ip, 4, 0);
    bpf_l3_csum_replace(skb, l3_off + offsetof(struct iphdr, check), 0, sum, 0);
	bpf_l4_csum_replace(skb, l4_off + offsetof(struct tcphdr, check), 0, sum, BPF_F_PSEUDO_HDR);

    return TC_ACT_OK;
}