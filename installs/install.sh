#/bin/bash

echo "Installing required modules"
sudo apt-get install -y lsof

# Uninstall any previous versions.
echo "Checking for any previous version..."
sudo kill $(sudo lsof -t -i:5666) > /dev/null 2>&1 || true

ver=$(echo $(sharpcd version) | sed "s/^.*Version: \([0-9.]*\).*/\1/")
vernum=$(echo "$ver" | sed -r 's/[.0]+//g')

if [[ $vernum =~ ^[0-9]+$ ]];
then
  if [[ $vernum < 0 ]];
  then
    echo "Breaking changes: Removing old sharpcd-data"
    read -r -p "Are you sure? [y/N] " response
    if [[ "$response" =~ ^([yY][eE][sS]|[yY])$ ]]
    then
        echo "Deleting old Data..."
        sudo rm -r /usr/local/bin/sharpcd-data
    else
        exit
    fi
  fi
fi

sudo rm -r /usr/local/bin/sharpcd
sudo rm /etc/systemd/system/sharpcd.service
sudo systemctl daemon-reload

# Download and unpack
wget https://github.com/Sharpz7/sharpcd/releases/download/XXXXX/linux.tar.gz
sudo tar -C /tmp/sharpcd -zxvf linux.tar.gz
sudo cp -n -R /tmp/sharpcd/* /usr/local/bin/
rm -r linux.tar.gz

# Create SharpCD User
if !(grep -c '^sharpcd:' /etc/passwd) then
    sudo chown sharpcd:sharpcd /usr/local/bin/docker-compose
    sudo useradd sharpcd
    sudo mkdir /home/sharpcd
    sudo chown -R sharpcd:sharpcd /home/sharpcd
    sudo groupadd docker
    sudo usermod -aG docker sharpcd
    sudo systemctl restart docker
fi

# Gen Keys
IP="$(hostname -I | cut -d' ' -f1)"
sudo sed -i -e "s/XXXXX/$IP/g" /usr/local/bin/sharpcd-data/private/openssl.conf
sudo openssl req -x509 -nodes -days 730 -newkey rsa:2048 -keyout /usr/local/bin/sharpcd-data/private/server.key -out /usr/local/bin/sharpcd-data/private/server.crt -config /usr/local/bin/sharpcd-data/private/openssl.conf -extensions 'v3_req'

# Permissions
sudo chmod +x /usr/local/bin/sharpcd
sudo chmod 755 /usr/local/bin/sharpcd
sudo chown -R sharpcd:sharpcd /usr/local/bin/sharpcd-data


# Create system service
test=$(cat <<-END
[Unit]
Description=SharpCD Service.
[Service]
Type=simple
User=sharpcd
Restart=always
WorkingDirectory=/usr/local/bin
ExecStart=/usr/local/bin/sharpcd server
[Install]
WantedBy=multi-user.target
END
)

# Create and enable service.
sudo touch /etc/systemd/system/sharpcd.service
sudo echo "$test" | sudo tee -a /etc/systemd/system/sharpcd.service
sudo systemctl enable sharpcd

# Initial Run
sudo sharpcd --secret Pass123 setsecret
sudo systemctl restart sharpcd

echo ""
echo "SHARPCD IS NOW RUNNING"
echo "======================="
echo "Your password is (Pass123). it is highly recommended you change this."
echo "The sharpcd server has now started and will startup on a reboot"
echo "Do sharpcd -h for more info!"
sharpcd --help
echo ""
echo "PLEASE NOTE"
echo "======================"
echo "SharpCD will be open on port 5666"
echo "Use the follow IP table commands to block port 5666 from outside localhost"
echo "sudo iptables -D INPUT -p tcp --dport 5666 -s localhost -j ACCEPT"
echo "sudo iptables -D INPUT -p tcp --dport 5666 -j DROP"
echo "sudo iptables -A INPUT -p tcp --dport 5666 -s localhost -j ACCEPT"
echo "sudo iptables -A INPUT -p tcp --dport 5666 -j DROP"
