#!/bin/bash
set -x

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
outdir=$SCRIPT_DIR/out

img=$outdir/ubuntu-cloudinit.img

rm -f $img
mkdosfs -n CIDATA -C $outdir/ubuntu-cloudinit.img 8192
mcopy -oi $img -s $SCRIPT_DIR/cloud-init/ubuntu/user-data ::
mcopy -oi $img -s $SCRIPT_DIR/cloud-init/ubuntu/meta-data ::
mcopy -oi $img -s $SCRIPT_DIR/cloud-init/ubuntu/network-config ::

