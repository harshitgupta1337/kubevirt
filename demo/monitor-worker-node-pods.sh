#!/bin/bash
wkr_node=$(kubectl get nodes -l '!node-role.kubernetes.io/control-plane' -o jsonpath='{.items[*].metadata.name}')

watch -n1 kubectl get pods -A --field-selector spec.nodeName=$wkr_node
