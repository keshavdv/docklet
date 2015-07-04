# -*- mode: ruby -*-
# vi: set ft=ruby :

$script = <<SCRIPT
export ETCD=v2.0.13

# Installing docker
echo 'Installing docker...'
if [ ! -f /usr/bin/docker ]; then
  apt-get update && apt-get install -y wget linux-image-generic-lts-trusty
  wget -qO- https://get.docker.com/ | sh
fi

# Install etcd
echo 'Installing etcd...'
if [ ! -f /usr/bin/etcd ]; then
  cd /usr/local/src/
  sudo wget https://github.com/coreos/etcd/releases/download/$ETCD/etcd-$ETCD-linux-amd64.tar.gz
  sudo tar -xzf etcd-$ETCD-linux-amd64.tar.gz && cd etcd-$ETCD-linux-amd64/ && sudo cp etc* /usr/bin
  sudo ETCD_BIND_ADDR=0.0.0.0 killall etcd && etcd -bind-addr=0.0.0.0 &
fi

# Installing haproxy
echo 'Installing haproxy...'
if [ ! -f /usr/sbin/haproxy ]; then
  add-apt-repository ppa:vbernat/haproxy-1.5
  apt-get update && apt-get install -y haproxy
fi

# Installing confd
echo 'Installing confd...'
if [ ! -f /usr/bin/confd ]; then
  sudo wget https://github.com/kelseyhightower/confd/releases/download/v0.10.0/confd-0.10.0-linux-amd64 -O /usr/bin/confd
  sudo chmod +x /usr/bin/confd
fi
SCRIPT

Vagrant.configure(2) do |config|

  config.vm.box = "ubuntu/trusty64"
  config.vm.hostname = "docklet"

  config.vm.provision "shell", inline: $script

end
