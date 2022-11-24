#!/bin/bash
set -euo pipefail

useradd mdns-ingress-publisher -d /home/mdns-ingress-publisher

mkdir -p /home/mdns-ingress-publisher/.kube
microk8s config >/home/mdns-ingress-publisher/.kube/config

cp mdns-ingress-publisher.service /etc/systemd/system
cp mdns-ingress-publisher /usr/local/bin

chmod 755 /etc/systemd/system/mdns-ingress-publisher.service
chown root:root /etc/systemd/system/mdns-ingress-publisher.service
systemctl enable mdns-ingress-publisher
