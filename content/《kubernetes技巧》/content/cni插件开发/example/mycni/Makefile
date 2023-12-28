all:
	go build -o bin/mycni-bin main.go

test:
	sudo ip netns add testing
	sudo NETCONFPATH=./conf CNI_PATH=./bin cnitool add mycni /var/run/netns/testing

clear:
	sudo ip link delete wyw0
	sudo ip netns delete testing
	sudo rm -rf /tmp/cni-host/*

hostlocal:
	sudo CNI_COMMAND=ADD CNI_CONTAINERID=cc1 CNI_NETNS=/var/run/netns/testing CNI_IFNAME=eth10 CNI_PATH=./bin ./bin/host-local < ./testconf/10-test.conf