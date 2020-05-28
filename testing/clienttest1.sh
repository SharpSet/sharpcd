sudo apt-get install lsof -y

go build -o "sharpcd" ./src
sudo ./sharpcd server &

sleep 2
./sharpcd

sudo kill $(sudo lsof -t -i:5666)