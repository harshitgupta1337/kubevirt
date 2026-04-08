#!/bin/bash

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" >/dev/null 2>&1 && pwd)"

# By default use the upstream QEMU repo
# and the v10.2.0 tag because it contains MSHV support
QEMU_REPO=${QEMU_REPO:-https://github.com/qemu/qemu.git}
QEMU_TAG=${QEMU_TAG:-v10.2.0}

if [ -z "$QEMU_VERSION" ]; then
    # Derive QEMU_VERSION from QEMU_TAG
    QEMU_VERSION=${QEMU_TAG#v}
    QEMU_VERSION=${QEMU_VERSION//-/.}
fi

# Fetch QEMU source code from upstream
rm -rf ./qemu-rpm-build ./qemu-rpms-out
mkdir -p ./qemu-rpm-build
pushd ./qemu-rpm-build

# Fetch the QEMU source code from the specified repository and tag
git clone --depth=1 -b ${QEMU_TAG} ${QEMU_REPO} qemu-${QEMU_VERSION}
pushd qemu-${QEMU_VERSION}/
meson subprojects download
popd 

# Create tarball of QEMU source code
tar -cf qemu-${QEMU_VERSION}.tar.xz \
    qemu-${QEMU_VERSION}
rm -rf qemu-${QEMU_VERSION}

# Copy spec file and related files
cp $SCRIPT_DIR/qemu-kvm/* .

# Remove any existing container with the same name
docker rm -f qemu-build

# Launch container for building the RPM
# mounting the current directory as a volume
docker run -td \
    --name qemu-build \
    -v $(pwd):/qemu-src \
    registry.gitlab.com/libvirt/libvirt/ci-centos-stream-9

# Build QEMU RPM
docker exec -w /qemu-src qemu-build bash -c "
  set -ex
  mkdir -p ~/rpmbuild/{BUILD,RPMS,SOURCES,SPECS,SRPMS}
  cp qemu-kvm.spec ~/rpmbuild/SPECS
  cp * ~/rpmbuild/SOURCES/
  cp qemu-${QEMU_VERSION}.tar.xz ~/rpmbuild/SOURCES/
  cd ~/rpmbuild/SPECS
  dnf update -y
  dnf -y install createrepo
  dnf builddep -y qemu-kvm.spec
  rpmbuild -ba qemu-kvm.spec
  cd ~/rpmbuild/RPMS
  createrepo --general-compress-type=gz --checksum=sha256 x86_64
"

popd

docker cp qemu-build:/root/rpmbuild/RPMS ./qemu-rpms-out

qemuRelease=$(cat $SCRIPT_DIR/qemu-kvm/qemu-kvm.spec | grep "Release: " | cut -d ' ' -f 2 | cut -d '%' -f 1)

cat >./qemu-rpms-out/build-info.json <<EOF
{
  "qemu_version": "17:${QEMU_VERSION}-${qemuRelease}.el9"
}
EOF
