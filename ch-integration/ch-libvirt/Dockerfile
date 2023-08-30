FROM mcr.microsoft.com/cbl-mariner/base/core:2.0

RUN tdnf -y update \
    && tdnf -y install dnf dnf-plugins-core \
    && dnf -y builddep libvirt \
    && tdnf -y install libvirt \
    #&& systemctl stop libvirtd.service \
    && systemctl disable libvirtd.service \
    && tdnf install -y git kernel-headers curl-devel

RUN git config --global http.sslVerify false \
    && git clone -b mshv https://github.com/praveen-pk/libvirt \
    && cd libvirt/ \
    && git submodule update --init --recursive subprojects/keycodemapdb/

RUN cd libvirt \
    && meson build --prefix=/ --libdir=usr/lib -Dapparmor=disabled -Ddocs=disabled -Ddriver_openvz=disabled -Ddriver_vbox=disabled -Ddriver_esx=disabled -Dstorage_iscsi_direct=disabled -Ddriver_qemu=disabled -Ddriver_ch=enabled -Ddriver_qemu=disabled -Dfirewalld=disabled -Dpolkit=disabled \
    && ninja -C build

RUN tdnf -y install vim wget

RUN wget --no-check-certificate https://github.com/cloud-hypervisor/cloud-hypervisor/releases/download/v31.0/cloud-hypervisor-static \
    && wget --no-check-certificate https://github.com/cloud-hypervisor/cloud-hypervisor/releases/download/v31.0/cloud-hypervisor \
    && chmod +x cloud-hypervisor-static \
    && chmod +x cloud-hypervisor \
    && mv cloud-hypervisor-static /usr/bin/ \
    && mv cloud-hypervisor /usr/bin/
