package bpf

//go:generate go run github.com/cilium/ebpf/cmd/bpf2go -target amd64 -cc clang toy_proxy__bpf toy-proxy-bpf.c -- -I $BPF_HEADERS
