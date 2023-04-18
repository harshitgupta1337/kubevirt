#!/bin/bash

if [ "$#" -ne 1 ]; then
    echo "Illegal number of arguments."
    echo "Correct usage: hack/build-mariner-based-imgs/build-images.sh <vmm>"
    echo "Exiting..."
    exit 1
fi

vmm=$1
if [ $vmm == "ch" ]; then
    LAUNCHER_BASE_IMG="docker.io/harshitg/virt-launcher-base:mariner"
elif [ $vmm == "qemu" ]; then
    LAUNCHER_BASE_IMG="mcr.microsoft.com/cbl-mariner/base/core:2.0"
else
    echo "Invalid args: VMM should be one of \"ch\" or \"qemu\". Exiting..."
    exit 1
fi

SCRIPT_DIR=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" &>/dev/null && pwd)

DOCKER_BUILDKIT=1 docker build -t afo-builder -f $SCRIPT_DIR/Dockerfile-builder $SCRIPT_DIR

# Building generic containers that are not dependent on the VMM
for ctr in virt-operator virt-api virt-handler virt-controller; do
    DOCKER_BUILDKIT=1 docker build --build-arg BUILDER_IMAGE=afo-builder:latest \
        -t harshitg/${ctr}:mariner -f $SCRIPT_DIR/Dockerfile-${ctr} .
done

# Building the virt-launcher container separately because
# It might need a different base image when building for CloudHypervisor
for ctr in virt-launcher; do
    DOCKER_BUILDKIT=1 docker build --build-arg BUILDER_IMAGE=afo-builder:latest \
        -t harshitg/${ctr}:mariner -f $SCRIPT_DIR/Dockerfile-${ctr}-${vmm} .
done
