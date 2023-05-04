#!/bin/bash

kubectl apply -f /home/ubuntu/git/kubevirt/demo/kubevirt-operator.yaml
kubectl apply -f /home/ubuntu/git/kubevirt/demo/kubevirt-cr.yaml
