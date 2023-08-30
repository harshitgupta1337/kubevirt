#!/bin/bash


sudo tdnf -y install mariner-repos-extended.noarch
sudo tdnf -y install mariner-repos-preview
sudo tdnf -y update

sudo iptables -P INPUT ACCEPT
sudo iptables -P OUTPUT ACCEPT
sudo iptables -P FORWARD ACCEPT

sudo tdnf install -y ms-kubernetes-*

sudo tdnf install -y conntrack

sudo modprobe br_netfilter

sudo kubeadm init --config kubeadm-config.yaml

mkdir .kube
sudo cp /etc/kubernetes/admin.conf .kube/config
sudo chown -R $USER .kube

kubectl get pods -A

kubectl apply -f https://raw.githubusercontent.com/flannel-io/flannel/v0.19.2/Documentation/kube-flannel.yml

kubectl taint node --all node-role.kubernetes.io/control-plane-
