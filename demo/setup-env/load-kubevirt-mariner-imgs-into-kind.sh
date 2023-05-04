#!/bin/bash
prefix=harshitg
tag=mariner

for comp in virt-launcher virt-controller virt-handler virt-api virt-operator; do
    kind load docker-image $prefix/$comp:$tag
done
