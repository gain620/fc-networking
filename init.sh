#!/bin/bash

# Update and install necessary tools
sudo apt update -y
sudo apt install acl git build-essential -y

sudo setfacl -m u:${USER}:rw /dev/kvm

[ $(stat -c "%G" /dev/kvm) = kvm ] && sudo usermod -aG kvm ${USER} \
&& echo "Access granted."

[ -r /dev/kvm ] && [ -w /dev/kvm ] && echo "OK" || echo "FAIL"

# Install Firecracker binary
ARCH="$(uname -m)"
release_url="https://github.com/firecracker-microvm/firecracker/releases"
latest=$(basename $(curl -fsSLI -o /dev/null -w  %{url_effective} ${release_url}/latest))
curl -L ${release_url}/download/${latest}/firecracker-${latest}-${ARCH}.tgz \
| tar -xz

# Rename the binary to "firecracker"
mv release-${latest}-$(uname -m)/firecracker-${latest}-${ARCH} firecracker
sudo mv firecracker /usr/local/bin/

# Clean up the old firecracker API unix socket
API_SOCKET="/tmp/firecracker.socket"
sudo rm -f $API_SOCKET

ARCH=linux-arm64
#ARCH=linux-amd64
GO_VERSION=1.22.1

# Update and install necessary tools
sudo apt-get update -y
sudo apt-get install -y tar git jq acl

# Install docker-ce
sudo apt-get install ca-certificates curl
sudo install -m 0755 -d /etc/apt/keyrings
sudo curl -fsSL https://download.docker.com/linux/debian/gpg -o /etc/apt/keyrings/docker.asc
sudo chmod a+r /etc/apt/keyrings/docker.asc

# Add the repository to Apt sources:
echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.asc] https://download.docker.com/linux/debian \
  $(. /etc/os-release && echo "$VERSION_CODENAME") stable" | \
  sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
sudo apt-get update -y

sudo apt-get install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin

sudo groupadd docker
sudo usermod -aG docker $USER

# Download the latest version of Go for corresponding architecture
GO_URL="https://golang.org/dl/go${GO_VERSION}.${ARCH}.tar.gz"
curl -L "$GO_URL" -o go${GO_VERSION}.${ARCH}.tar.gz

# Extract the archive and install Go
sudo tar -C /usr/local -xvzf go*.${ARCH}.tar.gz
rm go*.${ARCH}.tar.gz

# Set up Go environment variables
sudo echo 'export PATH=$PATH:/usr/local/go/bin' >> $HOME/.profile
sudo echo 'export PATH=$PATH:/usr/local/go/bin' >> $HOME/.bashrc

# Apply the environment variables
source $HOME/.profile
source $HOME/.bashrc