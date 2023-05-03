#!/bin/bash

kubectl apply -f /home/ubuntu/git/kubevirt/_out/manifests/release/kubevirt-operator.yaml
kubectl apply -f /home/ubuntu/git/kubevirt/_out/manifests/release/kubevirt-cr.yaml
