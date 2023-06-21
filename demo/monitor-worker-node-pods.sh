#!/bin/bash
wkr_node=$(kubectl get nodes -l '!node-role.kubernetes.io/control-plane' -o jsonpath='{.items[*].metadata.name}')

watch -n1 "kubectl get pods -A --field-selector spec.nodeName=$wkr_node,metadata.namespace!=kube-system,metadata.namespace!=calico-system,metadata.namespace!=calico-apiserver,metadata.namespace!=local-path-storage | grep -v tigera | grep -v cdi-apiserver | grep -v cdi-deployment | grep -v cdi-operator | grep -v cdi-uploadproxy"
