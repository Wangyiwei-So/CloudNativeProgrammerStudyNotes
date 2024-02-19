set -x
FILE_PATH="history.txt"

source $FILE_PATH
 
sudo docker rm -f $CONTAINER1_ID $CONTAINER2_ID $CONTAINER3_ID

sudo iptables -t nat -F $KUBE_SVCS
sudo iptables -t nat -F $SVC_WEBAPP
sudo iptables -t nat -F $WEBAPP_EP1
sudo iptables -t nat -F $WEBAPP_EP2

KUBE_SVC_IDS=$(sudo iptables -t nat -L -n --line-numbers | grep MY_KUBE-SERVICES | grep -E '^\s*[0-9]+\s+' | awk '{print $1}')
for CHAIN_ID in $KUBE_SVC_IDS; do
    sudo iptables -t nat -D OUTPUT $CHAIN_ID
done

sudo iptables -t nat -X $KUBE_SVCS
sudo iptables -t nat -X $SVC_WEBAPP
sudo iptables -t nat -X $WEBAPP_EP1
sudo iptables -t nat -X $WEBAPP_EP2
