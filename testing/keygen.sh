IP="$(hostname -I | cut -d' ' -f1)"
sudo sed -i -e "s/XXXXX/$IP/g" ./private/openssl.conf
sudo openssl req -x509 -nodes -days 730 -newkey rsa:2048 -keyout ./private/server.key -out ./private/server.crt -config ./private/openssl.conf -extensions 'v3_req'