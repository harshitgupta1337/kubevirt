#!/bin/bash

ctr_ip=$(docker exec -it kind-control-plane ip a show dev eth0 | grep "inet\>" | awk '{print $2}' | awk -F "/" '{print $1}')
ssh-keygen -R $ctr_ip

docker exec -it kind-control-plane bash -c "apt -y update; apt -y install ssh; service ssh start"

SSH_KEY=$(cat ~/.ssh/id_rsa.pub )
docker exec -it kind-control-plane bash -c "mkdir -p /root/.ssh/; echo \"$SSH_KEY\" >> /root/.ssh/authorized_keys"
