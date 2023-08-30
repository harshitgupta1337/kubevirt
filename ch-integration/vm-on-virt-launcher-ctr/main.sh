#!/bin/bash

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

docker build -t virt-launcher:test \
    -f $SCRIPT_DIR/virt-launcher-test-image/Dockerfile $SCRIPT_DIR/virt-launcher-test-image/

$SCRIPT_DIR/create-prereq-files.sh

ctr=$(docker run -d --privileged --name virt-launcher-test virt-launcher:test)

# Copy the needed files into the container
$SCRIPT_DIR/copy-files-to-docker.sh $ctr

# Now start libvirtd service
echo "Starting libvirtd"
docker exec -d $ctr /libvirt/build/src/libvirtd

echo "Started libvirtd. Now waiting 5 secs before creating VM."
sleep 5

# Now run the VM in the container
docker exec -it $ctr virsh create /vm.xml

# Now display the VM
docker exec -it $ctr virsh console clh_dom0_test
