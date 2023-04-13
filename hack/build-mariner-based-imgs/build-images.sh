#!/bin/bash

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

DOCKER_BUILDKIT=1 docker build -t afo-builder -f $SCRIPT_DIR/Dockerfile-builder $SCRIPT_DIR

for ctr in virt-launcher virt-operator virt-api virt-handler virt-controller
do
    DOCKER_BUILDKIT=1 docker build --build-arg BUILDER_IMAGE=afo-builder:latest \
        -t harshitg/${ctr}:mariner -f $SCRIPT_DIR/Dockerfile-${ctr} .
done
