#include <linux/bpf.h> // struct bpf_sock_addr
#include <bpf/bpf_helpers.h>
#include <bpf/bpf_tracing.h>

#define SYS_REJECT	0
#define SYS_PROCEED	1

static int 
__sock4_xlate_fwd(struct bpf_sock_addr *ctx){
    const __be32 cluster_ip = 0x846F070A; // 10.7.111.132
    const __be32 pod_ip = 0x0529050A;     // 10.5.41.5
    if (ctx->user_ip4 != cluster_ip) {
        return 0;
    }
    ctx->user_ip4 = pod_ip;
    return 0;
}

SEC("connect4")
int sock4_connect(struct bpf_sock_addr* ctx){
    __sock4_xlate_fwd(ctx);
    return SYS_PROCEED;
}