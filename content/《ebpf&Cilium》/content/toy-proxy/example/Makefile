export CLUSTER_IP = ""
export PORT = ""
export POD_IP = ""
export PROTO = tcp
userspace:
	sudo iptables -t nat -A OUTPUT -p $PROTO -d $CLUSTER_IP --dport $PORT -j REDIRECT --to-port $PORT

iptable:
	sudo iptables -t nat -A OUTPUT -p $PROTO -d $CLUSTER_IP --dport $PORT -j DNAT --to-destination $POD_IP:$PORT

generate:
	go generate ./... 