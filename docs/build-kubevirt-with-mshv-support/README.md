
## Build QEMU RPMs with MSHV accelerator enabled

The QEMU RPMs in CentOS Stream 9 repo, from where upstream KubeVirt's build scripts pull the QEMU (and other) RPMs, do not contain all the code patches needed for deployment on HyperV-Direct (L1VH). Therefore, we need to build custom QEMU RPMs.

### Fetch the `qemu-kvm` SPEC

NOTE: The `qemu-kvm` SPEC should be cloned and be placed in the same directory as the `dockerized-build-qemu-rpm.sh` script.

```bash
git clone -b kubevirt-mshv https://github.com/harshitgupta1337/qemu-kvm.git
```

### Build QEMU RPMs

Build the desired version of QEMU. 

NOTE: Please keep in mind that the SPEC that was downloaded in the previous step hardcodes the version and release of the QEMU RPM to `10.2.0-23`. So, the behavior would be undefined if you tried to build a different tag.

```bash
QEMU_TAG=v10.2.0 ./dockerized-build-qemu-rpm.sh
```

The above step creates a directory called `qemu-rpms-out` with the RPMs and other metadata.

### Build the HTTP Server image with RPMs

```bash
./build-rpms-http-image.sh -r docker.io/<username> -i qemu-rpms -t kubevirt-mshv -d ./qemu-rpms-out/
```

## Build KubeVirt with the Custom QEMU RPMs

### Clone KubeVirt

```bash
git clone https://github.com/kubevirt/kubevirt.git
```

### Configure KubeVirt to pull Custom RPMs

```bash
export DOCKER_TAG=<kubevirt-image-tag>

cd ./kubevirt/  # enter the cloned KubeVirt source code dir

# Setup Libvirt and QEMU RPMs to be pulled from HTTP server
# image built in a previous step
../setup-libvirt-qemu-rpms.sh \
    -q docker.io/<username>/qemu-rpms:kubevirt-mshv
    
# Make and push images to container registry
make bazel-build-images && make bazel-push-images
```
