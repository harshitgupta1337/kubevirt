# Integration of Cloud-Hypervisor/MSHV stack

## Source repo and branch

Fork from the upstream KubeVirt project. To download the latest source code, use the following command.

`git clone -b guptaharshit/hermes-ch-main git@github.com:harshitgupta1337/kubevirt.git`

## Building Mariner-based KubeVirt component images

First, generate the boiler-plate code using the following commands.

```

make generate

make

```

Build the Mariner-based images for KubeVirt components, such as `virt-operator`, `virt-launcher`, etc. following the build tooling used by Nexus.

```bash

# Set the Docker registry, username and image tag
export DOCKER_PREFIX=docker.io/harshitg   # Please replace with own preferred values
export DOCKER_TAG=mariner                 # Please replace with own preferred values

# This var can be set to true to speed up development
# If set to false, KubeVirt Operator manifest will refer to images by SHA
export KUBEVIRT_ONLY_USE_TAGS=true

# Build Mariner-based images
./hack/build-mariner-based-imgs/build-images.sh ch

# Push the Mariner-based images to container registry
./hack/build-mariner-based-imgs/push-images.sh

```

## Generate KubeVirt manifests

Now, we can generate the manifest YAMLs to install a KubeVirt deployment based on the images we just built in the previous step.

```bash
# Set the Docker registry, username and image tag
export DOCKER_PREFIX=docker.io/harshitg   # Please replace with own preferred values
export DOCKER_TAG=mariner                 # Please replace with own preferred values

# This var can be set to true to speed up development
# If set to false, KubeVirt Operator manifest will refer to images by SHA
export KUBEVIRT_ONLY_USE_TAGS=true

./hack/build-mariner-based-imgs/make-mariner-manifests.sh

```

The above script generates two manifest files, which can be installed in the K8s cluster using the `kubectl apply -f <file>` command:

- YAML containing the CRDs used by KubeVirt and the KubeVirt Operator deployment.
  
  `./_out/manifests/release/kubevirt-operator.yaml`

- YAML containing an instance of KubeVirt. 

  `./_out/manifests/release/kubevirt-cr.yaml` 

## Testing KubeVirt

### Prerequisites

- Install CDI on K8s cluster

### Stage 0

Relevant files are located in `./demo/` directory. VM spec for Stage 0 VM can be found at `./demo/stage-0-vm.yaml`.

### Stage 1

Relevant files are located in `./stage1/` directory.