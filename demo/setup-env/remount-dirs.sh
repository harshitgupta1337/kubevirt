#!/bin/bash

# Create directories in /mnt drive
mkdir -p /mnt/opt/local-path-provisioner
mkdir -p /mnt/run/containerd
mkdir -p /mnt/var/lib/docker
mkdir -p /mnt/var/lib/containerd

# Mount K8s and container runtime storage directories to /mnt that has more space
mount --bind /mnt/opt/local-path-provisioner /opt/local-path-provisioner
mount --bind /mnt/run/containerd /run/containerd/
mount --bind /mnt/var/lib/docker /var/lib/docker
mount --bind /mnt/var/lib/containerd /var/lib/containerd

# Restart the system daemons
systemctl daemon-reload
systemctl restart containerd
systemctl restart docker
