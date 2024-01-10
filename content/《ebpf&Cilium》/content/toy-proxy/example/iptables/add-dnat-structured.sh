
set -x
FILE_PATH="history.txt"
source $FILE_PATH

KUBE_SVCS="MY_KUBE-SERVICES"        # chain that serves as kubernetes service portal
SVC_WEBAPP="MY_KUBE-SVC-WEBAPP"     # chain that serves as DNAT entrypoint for webapp
WEBAPP_EP1="MY_KUBE-SEP-WEBAPP1"    # chain that performs dnat to pod1
WEBAPP_EP2="MY_KUBE-SEP-WEBAPP2"    # chain that performs dnat to pod2
PORT="80"
CLUSTER_IP="172.15.0.2"
PROTO="tcp"

echo "KUBE_SVCS=$KUBE_SVCS" >> $FILE_PATH
echo "SVC_WEBAPP=$SVC_WEBAPP" >> $FILE_PATH
echo "WEBAPP_EP1=$WEBAPP_EP1" >> $FILE_PATH
echo "WEBAPP_EP2=$WEBAPP_EP2" >> $FILE_PATH

#OUTPUT -> KUBE-SERVICES
sudo iptables -t nat -N $KUBE_SVCS
sudo iptables -t nat -A OUTPUT -p all -s 0.0.0.0/0 -d 0.0.0.0/0 -j $KUBE_SVCS

# KUBE-SERVICES -> KUBE-SVC-WEBAPP
sudo iptables -t nat -N $SVC_WEBAPP
sudo iptables -t nat -A $KUBE_SVCS -p $PROTO -s 0.0.0.0/0 -d $CLUSTER_IP --dport $PORT -j $SVC_WEBAPP

# KUBE-SVC-WEBAPP -> KUBE-SEP-WEBAPP*
sudo iptables -t nat -N $WEBAPP_EP1
sudo iptables -t nat -N $WEBAPP_EP2
sudo iptables -t nat -A $SVC_WEBAPP -p $PROTO -s 0.0.0.0/0 -d 0.0.0.0/0 -m statistic --mode random --probability 0.5  -j $WEBAPP_EP1
sudo iptables -t nat -A $SVC_WEBAPP -p $PROTO -s 0.0.0.0/0 -d 0.0.0.0/0 -m statistic --mode random --probability 1.0  -j $WEBAPP_EP2
sudo iptables -t nat -A $WEBAPP_EP1 -p $PROTO -s 0.0.0.0/0 -d 0.0.0.0/0 --dport $PORT -j DNAT --to-destination $CONTAINER2_IP:$PORT
sudo iptables -t nat -A $WEBAPP_EP2 -p $PROTO -s 0.0.0.0/0 -d 0.0.0.0/0 --dport $PORT -j DNAT --to-destination $CONTAINER3_IP:$PORT
