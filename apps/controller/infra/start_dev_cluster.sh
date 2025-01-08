#!/bin/bash

which k3d
if [ $? != 0 ]; then
    echo "K3d is not installed on your system.";
    exit 1
fi

which helm # dependency of 'install_nginx_ingress_controller.sh'
if [ $? != 0 ]; then
    echo "Helm is not installed on your system.";
    exit 1
fi

which docker
if [ $? != 0 ]; then
    echo "Docker is not installed on your system.";
    exit 1
fi

echo "\nCreating container registry...";
k3d registry create eigen-cloud-dev-reg.localhost \
    --port 5000

echo "\nCreating development K8s cluster...";
k3d cluster create eigen-dev-cluster \
    --registry-use k3d-eigen-cloud-dev-reg.localhost:5000 \
    --port "80:80@loadbalancer" \
    --port "8080:8080@loadbalancer"

echo "\nConfiguring your kubeconfig...";
mkdir ~/.kube
touch ~/.kube/config
k3d kubeconfig get eigen-dev-cluster > ~/.kube/config
chmod 600 ~/.kube/config

echo "\nBuilding Controller container and pushing image to K3d registry...";
docker build -t eigen-controller ..
docker tag eigen-controller k3d-eigen-cloud-dev-reg.localhost:5000/eigen-controller
docker push k3d-eigen-cloud-dev-reg.localhost:5000/eigen-controller

echo "\nInstalling NGINX Ingress Controller...";
./install_nginx_ingress_controller.sh

echo "\n\nContainer registry used by dev cluster: 3d-eigen-cloud-dev-reg.localhost:5000";

echo "You're good to go! :)";