## Prerequisites

IMPORTANT: Do this only when the machine is rebooted, or if a new machine is being setup.

Execute the script `./remount-dirs.sh` to mount K8s and container runtime storage directories on `/mnt`.

IMPORTANT: Remember to keep running `sudo docker system prune` regularly to cleanup disk space.

## Create a K8s cluster

```
git clone https://mariner-org@dev.azure.com/mariner-org/ECF/_git/platform-helm-charts
cd platform-helm-charts/

# Build the Helm charts
make build

# Remove any existing KinD cluster
./tests/scripts/cleanup.sh

# Create a KinD (K8s in Docker) cluster
./tests/scripts/create-kind-cluster.sh

# Install the Calico Container Networking Inteface (CNI) using Tigera Operator (operator for Calico)
helm install tigera-operator ./output/tigera-operator-0.4.0.tgz \
	--set-string tigera-operator.installation.calicoNetwork.nodeAddressAutodetectionV4.interface=eth0

# Wait for the node to be ready
kubectl wait node/kind-control-plane --for=condition=Ready

# Install CDI charts
./tests/cdi/install-cdi-charts.sh

```

## Build KubeVirt Mariner-based images

### Get KubeVirt
```
git clone -b guptaharshit/hermes-ch-main https://github.com/harshitgupta1337/kubevirt.git
cd kubevirt/
```

### Build Mariner-based images
```
./hack/build-mariner-based-imgs/build-images.sh ch
```

### Push images to Docker registry [This step is optional]
```
./hack/build-mariner-based-imgs/push-images.sh
```

### Load the images into the KinD cluster provisioned above
```
./demo/setup-env/load-kubevirt-mariner-imgs-into-kind.sh
```

### Generate manifest YAMLs
```
./hack/build-mariner-based-imgs/make-mariner-manifests.sh
```
The result of this command is 2 YAML files 
* `_out/manifests/release/kubevirt-cr.yaml`, and
* `_out/manifests/release/kubevirt-operator.yaml`

## Test KubeVirt + CloudHypervisor

### Install KubeVirt on KinD cluster

```
kubectl apply -f _out/manifests/release/kubevirt-operator.yaml

kubectl apply -f _out/manifests/release/kubevirt-cr.yaml
```

## Create the test VM

```
kubectl apply -f ./ch-integration/ch-libvirt/test-ch-vm.yaml
```
