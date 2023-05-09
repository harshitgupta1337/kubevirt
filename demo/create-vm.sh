#!/bin/bash
SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

kubectl apply -f $SCRIPT_DIR/../ch-integration/ch-libvirt/test-ch-vm-cloudinit.yaml
