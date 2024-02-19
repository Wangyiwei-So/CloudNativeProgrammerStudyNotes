tools:
ifeq ($(wildcard ./ecli),)
	wget https://aka.pw/bpf-ecli -O ecli && chmod +x ./ecli
endif
ifeq ($(wildcard ./ecc),)
	wget https://github.com/eunomia-bpf/eunomia-bpf/releases/latest/download/ecc && chmod +x ./ecc
endif

vmlinux:
	bpftool btf dump file /sys/kernel/btf/vmlinux format c > vmlinux.h