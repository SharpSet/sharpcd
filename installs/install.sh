#/bin/bash

if [[ $1 == "client" ]];
then
  # Download and unpack
  wget https://github.com/SharpSet/sharpcd/releases/download/XXXXX/linux.tar.gz
  sudo mkdir -p /tmp/sharpcd
  sudo tar -C /tmp/sharpcd -zxvf linux.tar.gz
  sudo cp /tmp/sharpcd/sharpcd /usr/local/bin/sharpcd
  rm -r linux.tar.gz

  sharpcd --help
  exit 0
fi

echo "Installing required modules"
sudo apt-get install -y lsof

daemon_installed=$(sudo lsof -t -i:5666)

# Uninstall any previous versions.
echo "Checking for any previous version..."
sudo kill $(sudo lsof -t -i:5666) > /dev/null 2>&1 || true

ver=$(echo $(sharpcd version) | sed "s/^.*Version: \([0-9.]*\).*/\1/")
vernum=$(echo "$ver" | sed -r 's/[.0]+//g')

# set breaking_version to $vernum < X
# Note that 4.4 => 44
breaking_version=$([[ $vernum < 0 ]] && echo "true" || echo "false")

if [[ $vernum =~ ^[0-9]+$ ]]; then
  # Set to 0 as there is no breaking version.
  if [[ $breaking_version == "true" ]]; then
    echo ""
    echo "BREAKING CHANGES: Removing old sharpcd-data"
    read -r -p "Are you sure? [y/N] " response
    if [[ "$response" =~ ^([yY][eE][sS]|[yY])$ ]]; then
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
wget https://github.com/SharpSet/sharpcd/releases/download/XXXXX/linux.tar.gz
sudo mkdir -p /tmp/sharpcd
sudo tar -C /tmp/sharpcd -zxvf linux.tar.gz
sudo cp -n -R /tmp/sharpcd/* /usr/local/bin/
rm -r linux.tar.gz

# Create SharpCD User
if !(grep -c '^sharpcd:' /etc/passwd) then
    sudo useradd sharpcd
    sudo mkdir /home/sharpcd
    sudo chown -R sharpcd:sharpcd /home/sharpcd
    sudo chown sharpcd:sharpcd /usr/local/bin/docker-compose
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
sudo systemctl restart sharpcd

echo ""
echo ""
echo "SHARPCD IS NOW RUNNING"
echo "======================="
echo ""
echo "SharpCD will be open on port 5666"
echo "Use the follow IP table commands to block port 5666 from outside localhost:"
echo ""
echo "sudo iptables -D INPUT -p tcp --dport 5666 -s localhost -j ACCEPT"
echo "sudo iptables -D INPUT -p tcp --dport 5666 -j DROP"
echo "sudo iptables -A INPUT -p tcp --dport 5666 -s localhost -j ACCEPT"
echo "sudo iptables -A INPUT -p tcp --dport 5666 -j DROP"

echo ""

read -p "Press enter to continue"
echo ""

# if $daemon_installed is false or breaking_version
if [[ $daemon_installed == "" || $breaking_version == "true" ]];
then
  echo "Please choose a password for your SharpCD server"
  sudo sharpcd setsecret

  echo "Please choose a valid github repository for your SharpCD server"
  echo "This is the only repository that SharpCD will accept compose files from"
  echo "Should be in the form https://raw.githubusercontent.com/{USER}/"
  echo ""

  read -p "Enter your repository name: " filter_loc
  sudo sharpcd addfilter $filter_loc

  echo "To get a github token, visit https://github.com/settings/tokens"
  read -p "Enter your GitHub token: " token
  sudo sharpcd changetoken $token
fi



echo ""
echo "SHARPCD IS NOW READY"
echo "====================="
echo "To learn more use sharpcd --help"
echo ""
