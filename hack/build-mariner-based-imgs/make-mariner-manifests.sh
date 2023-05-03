#!/bin/bash

SCRIPT_DIR=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" &>/dev/null && pwd)

cd $SCRIPT_DIR/../../

export KUBEVIRT_ONLY_USE_TAGS=true
export DOCKER_PREFIX=docker.io/harshitg
export DOCKER_TAG=mariner
export FEATURE_GATES="Root"

make manifests

cd -
