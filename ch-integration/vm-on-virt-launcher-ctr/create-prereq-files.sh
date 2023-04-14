#!/bin/bash

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

outdir=$SCRIPT_DIR/out
rm -rf $outdir
mkdir -p $outdir

# Download the rootfs image
wget https://cloud-images.ubuntu.com/focal/current/focal-server-cloudimg-amd64.img -O $outdir/focal-server-cloudimg-amd64.img
qemu-img convert -p -f qcow2 -O raw $outdir/focal-server-cloudimg-amd64.img $outdir/focal-server-cloudimg-amd64.raw

# Download the firmware
wget https://github.com/cloud-hypervisor/rust-hypervisor-firmware/releases/download/0.4.2/hypervisor-fw -O $outdir/hypervisor-fw

# Create cloudinit disk
$SCRIPT_DIR/create-cloud-init.sh

# Download the cloud-hypervisor binary
wget https://github.com/cloud-hypervisor/cloud-hypervisor/releases/download/v31.0/cloud-hypervisor-static -O $outdir/cloud-hypervisor-static
chmod +x $outdir/cloud-hypervisor-static
