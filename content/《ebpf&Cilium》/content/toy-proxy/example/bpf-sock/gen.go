package bpfsock

//go:generate go run github.com/cilium/ebpf/cmd/bpf2go -target amd64 -cc clang toy_proxy_bpf_sock toy-proxy-bpf-sock.c -- -I $BPF_HEADERS
