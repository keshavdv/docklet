# -*- mode: ruby -*-
# vi: set ft=ruby :

$script = <<SCRIPT
export ETCD=v2.0.13

# Installing docker
echo 'Installing docker...'
apt-get update && apt-get install -y wget linux-image-generic-lts-trusty
wget -qO- https://get.docker.com/ | sh

# Install etcd
echo 'Installing etcd...'
cd /usr/local/src/
sudo wget https://github.com/coreos/etcd/releases/download/$ETCD/etcd-$ETCD-linux-amd64.tar.gz
sudo tar -xzf etcd-$ETCD-linux-amd64.tar.gz && cd etcd-$ETCD-linux-amd64/ && sudo cp etc* /usr/bin
sudo etcd --listen-client-urls 'http://0.0.0.0:2379,http://0.0.0.0:4001' &

# Installing vulcand
echo 'Installing vulcand...'
# download vulcand from the trusted build
docker pull mailgun/vulcand:v0.8.0-beta.2

# launch vulcand in a container
docker run -p 8182:8182 -p 8181:8181 mailgun/vulcand:v0.8.0-beta.2 /go/bin/vulcand -apiInterface=0.0.0.0 --etcd=http://127.0.0.1:4001
SCRIPT

Vagrant.configure(2) do |config|

  config.vm.box = "ubuntu/trusty64"
  config.vm.hostname = "docklet"

  config.vm.provision "shell", inline: $script

end
