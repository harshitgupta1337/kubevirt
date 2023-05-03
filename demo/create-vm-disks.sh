#!/bin/bash

kubectl apply -f ~/git/kubevirt/ch-integration/ch-libvirt/disk-vol/os-dv.yaml

kubectl apply -f ~/git/kubevirt/ch-integration/ch-libvirt/cloudinit/cloudinit-dv.yaml
