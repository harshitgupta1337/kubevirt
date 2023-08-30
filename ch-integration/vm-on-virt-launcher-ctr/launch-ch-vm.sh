#!/bin/bash

set -e

./cloud-hypervisor-static \
    --kernel ./hypervisor-fw \
    --disk path=focal-server-cloudimg-amd64.raw --disk path=ubuntu-cloudinit.img \
    --cpus boot=4 \
    --memory size=1024M \
    --net "tap=,mac=,ip=,mask="
