#!/bin/bash

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

kubectl apply -f $SCRIPT_DIR/../demo/kubevirt-operator.yaml
kubectl apply -f $SCRIPT_DIR/../demo/kubevirt-cr.yaml
