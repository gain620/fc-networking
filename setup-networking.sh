#!/bin/bash

# Run this once on the host, or upon reboot

TAP_DEV=fctap0
HOST_IFNAME=enP2p4s0

TAP_IP="172.16.0.1"
MASK_SHORT="/24"

# Add a tap device to act as a bridge between the microVM
# and the host.
sudo ip tuntap add dev $TAP_DEV mode tap

# The subnet is 172.16.0.0/24 and so the 
# host will be 172.16.0.1 and the microVM is going to be set to 
# 172.16.0.2
sudo ip addr add 172.16.0.1/24 dev $TAP_DEV
sudo ip link set $TAP_DEV up
ip addr show dev $TAP_DEV

# Set up IP forwarding and masquerading

# Change IFNAME to match your main ethernet adapter, the one that
# accesses the Internet - check "ip addr" or "ifconfig" if you don't 
# know which one to use.
IFNAME=$HOST_IFNAME

# Enable IP forwarding
sudo sh -c "echo 1 > /proc/sys/net/ipv4/ip_forward"

# Enable masquerading / NAT - https://tldp.org/HOWTO/IP-Masquerade-HOWTO/ipmasq-background2.5.html
sudo iptables -t nat -A POSTROUTING -o $IFNAME -j MASQUERADE
sudo iptables -A FORWARD -m conntrack --ctstate RELATED,ESTABLISHED -j ACCEPT
sudo iptables -A FORWARD -i $TAP_DEV -o $IFNAME -j ACCEPT


#ip addr add 172.16.0.5/24 dev eth0
#ip link set eth0 up
#ip route add default via 172.16.0.1 dev eth0

# Set up nameserver
#echo "nameserver 1.1.1.1" > /etc/resolv.conf

