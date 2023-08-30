#!/bin/bash
SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
files_dir=$SCRIPT_DIR/out

ctr=$1

docker cp $files_dir/focal-server-cloudimg-amd64.raw $ctr:/
docker cp $files_dir/hypervisor-fw $ctr:/
docker cp $files_dir/cloud-hypervisor-static $ctr:/
docker cp $files_dir/ubuntu-cloudinit.img $ctr:/
docker cp $SCRIPT_DIR/launch-ch-vm.sh $ctr:/
docker cp $SCRIPT_DIR/vm.xml $ctr:/
