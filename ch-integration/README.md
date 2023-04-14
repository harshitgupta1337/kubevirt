# Integrating Kubevirt with Cloud-Hypervisor

## Running a small CH VM on a virt-launcher container
Run the script `ch-integration/vm-on-virt-launcher-ctr/main.sh`.

The script will create the prereq files for the VM, run a virt-launcher container and start a CH VM inside the container.

After the VM is up, a console will appear to login. Username is `cloud` and password is `cloud123`.

After testing is done, use `docker rm -f <container_name` command to delete the Docker container.
