IP="$(hostname -I | cut -d' ' -f1)"
sudo sed -i -e "s/XXXXX/$IP/g" ./internal/private/openssl.conf
sudo openssl req -x509 -nodes -days 730 -newkey rsa:2048 -keyout ./internal/private/server.key -out ./internal/private/server.crt -config ./internal/private/openssl.conf -extensions 'v3_req'