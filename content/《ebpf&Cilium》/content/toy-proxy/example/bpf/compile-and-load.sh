set -x

NIC=docker0

# compile c code into bpf code
# clang -O2 -Wall -c toy-proxy-bpf.c -target bpf -o toy-proxy-bpf.o

# add tc queuing discipline (egress and ingress buffer)
sudo tc qdisc del dev $NIC clsact 2>&1 >/dev/null
sudo tc qdisc add dev $NIC clsact

# load bpf code into the tc egress and ingress hook respectively
sudo tc filter add dev $NIC egress bpf da obj toy_proxy__bpf_bpfel_x86.o sec egress
sudo tc filter add dev $NIC ingress bpf da obj toy_proxy__bpf_bpfel_x86.o sec ingress

# show info
sudo tc filter show dev $NIC egress
sudo tc filter show dev $NIC ingress