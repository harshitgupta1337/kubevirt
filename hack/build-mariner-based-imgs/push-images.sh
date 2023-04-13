#!/bin/bash

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

for ctr in virt-launcher virt-operator virt-api virt-handler virt-controller
do
    docker push harshitg/${ctr}:mariner
done
