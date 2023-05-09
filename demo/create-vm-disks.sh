#!/bin/bash
SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

kubectl apply -f $SCRIPT_DIR/../ch-integration/ch-libvirt/disk-vol/os-dv.yaml

kubectl apply -f $SCRIPT_DIR/../ch-integration/ch-libvirt/cloudinit/cloudinit-dv.yaml
