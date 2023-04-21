#!/bin/bash

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

docker build -t virt-launcher:test \
    -f $SCRIPT_DIR/virt-launcher-test-image/Dockerfile $SCRIPT_DIR/virt-launcher-test-image/

$SCRIPT_DIR/create-prereq-files.sh

ctr=$(docker run -d --privileged virt-launcher:test)

# Copy the needed files into the container
$SCRIPT_DIR/copy-files-to-docker.sh $ctr

# Now run the VM in the container
docker exec -it $ctr virsh create /vm.xml
