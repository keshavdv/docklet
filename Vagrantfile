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
  sudo wget -q https://github.com/coreos/etcd/releases/download/$ETCD/etcd-$ETCD-linux-amd64.tar.gz
  sudo tar -xzf etcd-$ETCD-linux-amd64.tar.gz && cd etcd-$ETCD-linux-amd64/ && sudo cp etc* /usr/bin
  sudo killall etcd 2> /dev/null
  sudo etcd --listen-client-urls 'http://0.0.0.0:2379,http://0.0.0.0:4001' --advertise-client-urls 'http://0.0.0.0:2379,http://0.0.0.0:4001' 2>&1 &
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
  sudo wget -q -O /usr/bin/confd https://github.com/kelseyhightower/confd/releases/download/v0.10.0/confd-0.10.0-linux-amd64 
  sudo chmod +x /usr/bin/confd && sudo chmod 755 /usr/bin/confd
fi

# Configuring confd
echo 'Configuring confd...'
sudo mkdir -p /etc/confd/conf.d
sudo mkdir -p /etc/confd/templates
sudo bash -c 'cat << EOF > /etc/confd/confd.toml
confdir = "/etc/confd"
interval = 20
backend = "etcd"
prefix = "/"
scheme = "http"
verbose = true
EOF'

sudo bash -c 'cat << EOF > /etc/confd/conf.d/haproxy.toml
[template]
src = "haproxy.cfg.tmpl"
dest = "/etc/haproxy/haproxy.cfg"
keys = [
        "/docklet/"
]
reload_cmd = "echo restarting && /usr/sbin/service haproxy reload"
EOF'

sudo bash -c 'cat << EOF > /etc/confd/templates/haproxy.cfg.tmpl
defaults
  log     global
  mode    http

listen stats :1936
    mode http
    stats enable
    stats hide-version
    stats realm Haproxy\ Statistics
    stats uri /
    stats auth admin:admin
EOF'

# Starting confd
echo 'Starting confd...'
sudo etcdctl setdir /docklet
sudo bash -c 'confd > /var/log/confd-docklet.log 2>&1 &'

export IP=$(/sbin/ifconfig eth1 | grep 'inet addr:' | cut -d: -f2 | awk '{ print $1}')
echo "Done (available at $IP). Enjoy =)"

SCRIPT

Vagrant.configure(2) do |config|

  config.vm.box = "ubuntu/trusty64"
  config.vm.hostname = "docklet"

  config.vm.network "public_network"
  config.vm.provision "shell", inline: $script

end
