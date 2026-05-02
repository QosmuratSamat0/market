#!/bin/bash
set -euo pipefail

# Redirect output to a log for cloud-init/user-data debugging.
exec > >(tee /var/log/user-data.log | logger -t user-data -s 2>/dev/console) 2>&1

echo "Provisioning started..."

export DEBIAN_FRONTEND=noninteractive

# Setup 2GB swap for the 1GB RAM free-tier instance.
if ! swapon --show=NAME | grep -q '^/swapfile$'; then
  if [ ! -f /swapfile ]; then
    fallocate -l 2G /swapfile || dd if=/dev/zero of=/swapfile bs=1M count=2048
    chmod 600 /swapfile
    mkswap /swapfile
  fi

  swapon /swapfile
fi

if ! grep -q '^/swapfile none swap sw 0 0$' /etc/fstab; then
  echo '/swapfile none swap sw 0 0' >> /etc/fstab
fi

# Install Docker, Compose, Git, and firewall persistence tools.
apt-get update
apt-get upgrade -y
apt-get install -y docker.io docker-compose-v2 git iptables-persistent netfilter-persistent

systemctl enable --now docker

# Open application and administration ports on the host firewall.
for port in 22 80 81 443 3000 3001 9090 9091; do
  iptables -C INPUT -m state --state NEW -p tcp --dport "$port" -j ACCEPT 2>/dev/null \
    || iptables -I INPUT 6 -m state --state NEW -p tcp --dport "$port" -j ACCEPT
done

netfilter-persistent save || iptables-save > /etc/iptables/rules.v4

# Create the deployment user and grant Docker/sudo access.
if ! id devops >/dev/null 2>&1; then
  useradd -m -s /bin/bash devops
fi

usermod -aG sudo,docker devops
echo "devops ALL=(ALL) NOPASSWD:ALL" > /etc/sudoers.d/90-devops
chmod 440 /etc/sudoers.d/90-devops

# Copy the initial SSH access from the default Ubuntu user.
mkdir -p /home/devops/.ssh
if [ -f /home/ubuntu/.ssh/authorized_keys ]; then
  cp /home/ubuntu/.ssh/authorized_keys /home/devops/.ssh/authorized_keys
fi
chown -R devops:devops /home/devops/.ssh
chmod 700 /home/devops/.ssh
chmod 600 /home/devops/.ssh/authorized_keys

# Create deployment directories.
mkdir -p /home/devops/deploy/app /home/devops/deploy/monitoring /home/devops/deploy/npm
chown -R devops:devops /home/devops/deploy

echo "Provisioning finished."
