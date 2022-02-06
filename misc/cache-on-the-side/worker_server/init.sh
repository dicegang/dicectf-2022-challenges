sudo apt-get update
sudo apt-get install -y make unzip gcc python3-pip
sudo pip3 install google-cloud-pubsub google-cloud-storage command_runner

sudo apt-get install -y ca-certificates curl gnupg lsb-release
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg
echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
sudo apt-get update
sudo apt-get install -y docker-ce docker-ce-cli containerd.io

sudo mv etc_systemd_system/* /etc/systemd/system/
sudo mv victim/victim /srv/
sudo mkdir /srv/attack_srv
sudo mv attack_srv/* /srv/attack_srv

cd /srv/attack_srv
sudo docker build -t attack .

sudo systemctl enable victim.service
sudo systemctl enable attack_manager.service
sudo systemctl start  victim.service
sudo systemctl start  attack_manager.service
