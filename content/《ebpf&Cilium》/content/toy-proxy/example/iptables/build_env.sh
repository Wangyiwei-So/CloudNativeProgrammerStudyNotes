set -x
CONTAINER1_ID=$(sudo docker run -d nginx:1.18-alpine)
sudo docker exec $CONTAINER1_ID sh -c 'echo "this is container1" > /usr/share/nginx/html/index.html'
CONTAINER2_ID=$(sudo docker run -d nginx:1.18-alpine)
sudo docker exec $CONTAINER2_ID sh -c 'echo "this is container2" > /usr/share/nginx/html/index.html'
CONTAINER3_ID=$(sudo docker run -d nginx:1.18-alpine)
sudo docker exec $CONTAINER3_ID sh -c 'echo "this is container3" > /usr/share/nginx/html/index.html'

CONTAINER1_IP=$(sudo docker inspect $CONTAINER1_ID | jq -r ".[0].NetworkSettings.IPAddress")
CONTAINER2_IP=$(sudo docker inspect $CONTAINER2_ID | jq -r ".[0].NetworkSettings.IPAddress")
CONTAINER3_IP=$(sudo docker inspect $CONTAINER3_ID | jq -r ".[0].NetworkSettings.IPAddress")

CONTAINER1_MAC=$(sudo docker inspect $CONTAINER1_ID | jq -r ".[0].NetworkSettings.MacAddress")
CONTAINER2_MAC=$(sudo docker inspect $CONTAINER2_ID | jq -r ".[0].NetworkSettings.MacAddress")
CONTAINER3_MAC=$(sudo docker inspect $CONTAINER3_ID | jq -r ".[0].NetworkSettings.MacAddress")

FILE_PATH="history.txt"
echo "CONTAINER1_ID=$CONTAINER1_ID" > $FILE_PATH
echo "CONTAINER2_ID=$CONTAINER2_ID" >> $FILE_PATH
echo "CONTAINER3_ID=$CONTAINER3_ID" >> $FILE_PATH

echo "CONTAINER1_IP=$CONTAINER1_IP" >> $FILE_PATH
echo "CONTAINER2_IP=$CONTAINER2_IP" >> $FILE_PATH
echo "CONTAINER3_IP=$CONTAINER3_IP" >> $FILE_PATH

echo "CONTAINER1_MAC=$CONTAINER1_MAC" >> $FILE_PATH
echo "CONTAINER2_MAC=$CONTAINER2_MAC" >> $FILE_PATH
echo "CONTAINER3_MAC=$CONTAINER3_MAC" >> $FILE_PATH
